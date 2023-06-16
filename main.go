package main

// golang datetime format 2006-01-02 15:04:05
import (
	"fmt"
	"os"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/asynctimer"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/services"
	"uwbwebapp/pkg/tools"
	"uwbwebapp/web"

	"github.com/ahmetb/go-linq/v3"

	"github.com/kataras/iris/v12"
)

func InitWebSite() {
	file, _ := os.Create("./pid")
	pid := fmt.Sprintf("%d", os.Getpid())
	_, err := file.Write([]byte(pid))
	if err != nil {
		fmt.Println(err.Error())
	}
	file.Close()
	fmt.Println("初始化 Web Server")
	conf.LoadWebConfig("conf/WebConfig.json")
	if dao.InitDatabase() {
		if cache.InitRedisDatabase() {
			biz.InitSystemLogger()
			biz.InitSysUserToRedis()
			biz.InitSysFuncPageToRedis()
			biz.InitDictToRedis()
			biz.InitAllUWBBaseStationsToRedis()

			// //!TEST
			// cal, _ := biz.GetSwimmerCalendarBySwimmerIdAndNoExitSite("a956c6b1-c89d-45c0-bfd5-7531dce97f63")
			// fmt.Println(cal.SiteID)
			web.ConnectionUWBTagMQTTServer("go_mqtt_client_subscriber") // 链接UWB MQTT 服务，并建立 Web Socket Server.
		}
	}
}

func main() {
	fmt.Println("==========================================================================")
	InitWebSite()
	fmt.Println("==========================================================================")
	fmt.Printf("UWB Swimmer Safety Management Platform Powered by Karonsoft. Version: %s\r\n", conf.WebConfiguration.Version)
	app := iris.New()
	app.AllowMethods("GET,POST,PUT,DELETE,OPTIONS")
	// 指定静态目录和视图模板目录
	app.HandleDir("/", "./web")
	app.RegisterView(iris.HTML("./web", ".html"))
	app.Use(iris.Compression)
	app.Use(web.Before)
	web.RegisterServices(app)

	app.Post("/upload", services.WSUploadFile)               // 测试上传文件
	app.Post("/uploadsitemap", services.WSUploadSiteMap)     // 上传场地平台图
	app.Get("/nodog/auth", web.NoDogSplashAuthorizedService) // NodogSplashAuthorizedService

	// app.Use(func(ctx iris.Context) {
	// 	if strings.Contains(ctx.GetHeader("Accept-Encoding"), "gzip") {
	// 		ctx.Header("Content-Encoding", "gzip")
	// 		gz := gzip.NewWriter(ctx.ResponseWriter())
	// 		defer gz.Close()
	// 	} else {
	// 		ctx.WriteString("Client browser is not gzip supported")
	// 	}
	// })

	//Begin Web Socket
	web.GlobalUWBWebSocket = web.NewSocket()
	web.MQTTWebSocketHandle(app, "/msg")

	web.GlobalAlertWebSocket = web.NewSocket()
	web.AlertWebSocketHandle(app, "/alert")

	go web.BoardcastToWebSocketClient()
	go asynctimer.AlertDangerInformation()
	go web.BoardcastToAlertWebSocketClient()

	//End Web Socket

	ros := app.GetRoutes()
	var paths []string
	for _, r := range ros {
		paths = append(paths, r.Path)
	}
	linq.From(paths).DistinctByT(
		func(r string) string {
			return r
		}).ToSlice(&paths)
	conf.WebConfiguration.UrlPathList = paths // 所有注册的服务接口地址。
	for _, path := range paths {
		app.Options(path, web.Cors)
	}
	config := iris.WithConfiguration(iris.Configuration{
		DisableStartupLog:   false,
		EnableOptimizations: true,
		Charset:             "UTF-8",
		PostMaxMemory:       conf.WebConfiguration.PostDataMaxMBSize,
	})
	// go biz.SpeakMessageAlert()
	//5.67 -13.487847209008407
	//points := biz.GenerateUWBBaseStationPointsBySiteId("f089a938-90fc-4ca7-9071-588e1c4da0cb")
	// points := []entities.Point{
	// 	{0, 0},
	// 	{24.893, -12.773},
	// 	{0.034, -14.547},
	// 	{11.993, -14.609},
	// 	{11.953, 0.166},
	// 	{24.695, 0.197},
	// }
	//point := entities.Point{X: 5.67, Y: -13.487847209008407}
	// for _, p := range points {
	// 	fmt.Println(p.X, p.Y)
	// }

	//test()
	var points []entities.Point
	points = append(points, entities.Point{X: 0, Y: 0})
	points = append(points, entities.Point{X: 0, Y: 100})
	points = append(points, entities.Point{X: 100, Y: 100})
	points = append(points, entities.Point{X: 100, Y: 0})
	points = append(points, entities.Point{X: 0, Y: 0})
	point := entities.Point{X: 90, Y: 90}
	if res, _ := tools.UseGeosInPolygon(point, points); res {
		fmt.Println("在里面")
	} else {
		fmt.Println("不在里面")
	}

	//! FOR TEST
	app.Run(iris.Addr(conf.WebConfiguration.Port), config)

}

func test() {

	// 多边形6个点的坐标
	points := biz.GenerateUWBBaseStationPointsBySiteId("f089a938-90fc-4ca7-9071-588e1c4da0cb")

	// 测试点
	point := entities.Point{X: 5.67, Y: -13.487847209008407}

	result, _ := tools.UseGeosInPolygon(point, points)
	if result {
		println("该点在多边形内") // 输出结果
	} else {
		println("该点不在多边形内")
	}
}
