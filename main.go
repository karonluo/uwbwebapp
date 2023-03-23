package main

// golang datetime format 2006-01-02 15:04:05
import (
	"fmt"
	"os"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/services"
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
	app.Use(web.Before)
	web.RegisterServices(app)

	app.Post("/upload", services.WSUploadFile)           // 测试上传文件
	app.Post("/uploadsitemap", services.WSUploadSiteMap) // 上传场地平台图

	//Begin Web Socket
	web.GlobalUWBWebSocket = web.NewSocket()
	web.MQTTWebSocketHandle(app, "/msg")

	web.GlobalAlertWebSocket = web.NewSocket()
	web.AlertWebSocketHandle(app, "/alert")

	go web.BoardcastToWebSocketClient()
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
	app.Run(iris.Addr(conf.WebConfiguration.Port), iris.WithCharset("UTF-8"))

}
