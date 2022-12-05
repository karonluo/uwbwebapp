package web

import (
	"encoding/json"
	"fmt"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/services"

	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

type WebMessage struct {
	StatusCode int
	Message    interface{}
}

func Index(ctx iris.Context) {
	wm := WebMessage{Message: "Hello, Karonsoft!", StatusCode: 200}
	ctx.StatusCode(wm.StatusCode)
	ctx.JSON(wm)
}
func TopCheck(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
	ctx.Next()
}

// 拦截器，用于验证用户信息和权限信息
func Before(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	if ctx.Request().Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.StatusCode(204)
		return
	}
	authorization := ctx.GetHeader("Authorization")
	requestPath := ctx.RequestPath(true)
	if authorization == "" {
		if requestPath != "/login" {
			message := WebMessage{Message: "未登录", StatusCode: 401}
			fmt.Println("Authorization Error!")
			ctx.StatusCode(401)
			ctx.JSON(message)
		} else {
			ctx.Next()
		}

	} else {
		// TODO: Check authorization in redis
		res, _ := biz.CheckLogin(authorization)
		// TODO: 为了方便测试，authorization 为 test 时，直接通过，不进行权限检测。
		if authorization == "test" {
			ctx.Next()
			return
		}

		if res {
			acl, err := biz.GetLoginInformation(authorization, "acl") // 获取访问权限
			if err != nil {
				tools.ProcessError("web.Before", `acl, err := biz.GetLoginInformation(authorization, "acl")`, err)
				message := WebMessage{Message: "无访问权限", StatusCode: 403}
				ctx.StatusCode(403)
				ctx.JSON(message)
			} else {

				var t []interface{}
				json.Unmarshal([]byte(acl), &t)
				if tools.HaveElementInRange(t, requestPath) {
					ctx.Next()
				} else {
					message := WebMessage{Message: "无访问权限", StatusCode: 403}
					ctx.StatusCode(403)
					ctx.JSON(message)
				}
			}
		} else {
			message := WebMessage{Message: "未登录", StatusCode: 401}
			ctx.StatusCode(401)
			ctx.JSON(message)
		}
	}
}
func Cors(ctx iris.Context) {
	ctx.Text("")
}

// 绑定用户相关WEB服务
func RegisterSysUserServices(app *iris.Application) {
	app.Options("/sysuser", Cors)
	app.Post("/sysuser", services.WSCreateUser)                 //新增系统用户
	app.Delete("/sysuser", services.WSDeleteSysUser)            //删除系统用户
	app.Put("/sysuser", services.WSUpdateSite)                  // 修改系统用户
	app.Get("/sysuser", services.WSGetSysUserFromDBByLoginName) //获取系统用户信息(通过登录名)

	app.Options("/login", Cors)
	app.Post("/login", services.WSLoginSystem)                                             //登录系统
	app.Get("/sysuser/listall", services.WSEnumSysUsers)                                   //列举系统用户
	app.Get("/sysuser/listall/fromcompanies", services.WSEnumSysUsersFromSportsCompanyIds) // 通过体育公司唯一编号集合获取所有下属系统用户

	app.Options("/sysuser/query", Cors)
	app.Get("/sysuser/query", services.WSQuerySysUsers) // 查询系统用户列表
}

// 绑定体育公司相关WEB服务
func RegisterSportsCompanyServices(app *iris.Application) {

	app.Options("/company", Cors)
	app.Post("/company", services.WSCreateSportsCompany)   //创建体育运动公司
	app.Delete("/company", services.WSDeleteSportsCompany) // 删除体育运动公司

	app.Options("/companies", Cors)
	app.Delete("/companies", services.WSDeleteSportsCompanies) // 批量删除体育运动公司

	app.Put("/company", services.WSUpdateSportsCompany)  // 更新修改体育运动公司
	app.Get("/company", services.WSGetSportsCompanyById) // 获取体育公司相信信息

	app.Options("/company/query", Cors)
	app.Get("/company/query", services.WSQueryCompanies) // 查询公司列表

	app.Options("/company/relsites", Cors)
	app.Put("/company/relsites", services.WSRelSportsCompanyAndSites) // 为体育公司管关联场地
}

// 绑定游泳者相关WEB服务
func RegisterSwimmerServices(app *iris.Application) {
	app.Options("/swimmer", Cors)
	app.Post("/swimmer", services.WSCreateSwimmer)    //创建游泳者信息
	app.Get("/swimmer", services.WSGetSwimmersById)   // 获取游泳者信息
	app.Delete("/swimmer", services.WSDeleteSwimmers) //批量删除游泳者信息
	app.Put("/swimmer", services.WSUpdateSwimmer)     // 修改游泳者（会员）信息

	app.Options("/swimmer/query", Cors)
	app.Get("/swimmer/query", services.WSQuerySwimmers) // 查询游泳者信息

}

// 绑定场地相关WEB服务
func RegisterSiteServices(app *iris.Application) {
	app.Options("/site", Cors)
	app.Post("/site", services.WSCreateSite)    //创建查询场地信息
	app.Get("/site", services.WSGetSiteById)    // 获取查询场地信息
	app.Delete("/site", services.WSDeleteSites) //批量删除场地信息
	app.Put("/site", services.WSUpdateSite)     // 修改场地信息

	app.Options("/site/query", Cors)
	app.Get("/site/query", services.WSQuerySites) // 查询场地信息

}
func RegisterServices(app *iris.Application) {

	app.Get("/", Index)
	RegisterSysUserServices(app)
	RegisterSportsCompanyServices(app)
	RegisterSwimmerServices(app)
	RegisterSiteServices(app)
	// app.Options("/siteowners/list", Cors)
	// app.Get("/siteowners/list", services.WSEnumSiteOwners) //列举场地负责人
	// app.Put("/siteowners/setowners", services.WSSetSiteOwners) //设置场地负责人

	app.Options("/athorizationinfo", Cors)
	app.Get("/athorizationinfo", services.WSGetLoginInformation) //获取登录信息

	app.Post("/role", services.WSCreateSysRole)               // 创建系统角色
	app.Get("/sysfuncpage/list", services.WSEnumSysFuncPages) //列举所有系统功能页面

	app.Options("/dict", Cors)
	app.Get("/dict", services.WSGetDictValues) // 获取字典值

	app.Options("/dict/children", Cors)
	app.Get("/dict/children", services.WSGetChildrenSystemDictsByParent)

}
