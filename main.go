package main

// golang datetime format 2006-01-02 15:04:05
import (
	"fmt"
	"os"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/dao"
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
		if dao.InitRedisDatabase() {
			dao.InitSysUserToRedis()
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
	app.Use(web.Before)
	web.RegisterServices(app)
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
	app.Run(iris.Addr(conf.WebConfiguration.Port), iris.WithCharset("UTF-8"))

}
