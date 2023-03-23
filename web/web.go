package web

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"uwbwebapp/conf"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/services"

	"uwbwebapp/pkg/tools"

	"github.com/ahmetb/go-linq/v3"
	"github.com/kataras/iris/v12"
)

type WebMessage struct {
	StatusCode int
	Message    interface{}
}

func Index(ctx iris.Context) {
	wm := WebMessage{Message: "Hello, Karonsoft!", StatusCode: 200}
	ctx.StatusCode(wm.StatusCode)
	err := ctx.View("/templates/qr.html")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TopCheck(ctx iris.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Credentials", "true")
	ctx.Header("Access-Control-Allow-Methods", "*")
	ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
	ctx.Next()
}

// 可以被跳过认证的页面。
func SkipAuthorizedAddress(address string) bool {
	skipAddresses := []string{
		"/wifidog/ping",
		"/wifidog/auth",
		"/wifidog/login",
		"/wifidog/portal",
		"/wifidog/msg",
	}
	return linq.From(skipAddresses).WhereT(func(s string) bool {
		return (s == address)
	}).Count() > 0

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

	requestPath := ctx.RequestPath(true)
	fmt.Println(requestPath)
	if SkipAuthorizedAddress(requestPath) {
		ctx.Next()
	} else {
		authorization := ctx.GetHeader("Authorization")
		if authorization == "" {
			if requestPath != "/login" &&
				requestPath != "/index" &&
				!strings.Contains(requestPath, "/msg") &&
				!strings.Contains(requestPath, "/alert") {
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
			} else {
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
							if conf.WebConfiguration.IsThroughBackendForWriteOperateLog {
								// TODO: 后端填写操作日志，建议前端进行。后端进行时性能非常低下，需要查询和转换各类信息。
								strUser, errUser := biz.GetLoginInformation(authorization, "sysuser")
								if errUser == nil {
									var user entities.SysUser
									json.Unmarshal([]byte(strUser), &user)
									var log entities.SystemLog
									funcPage := biz.GetSysFuncPageFromRedis(ctx.Path(), ctx.Method())
									log.LogType = "operate"
									log.Message = funcPage.DisplayName
									log.FunctionName = funcPage.URLAddress
									log.ModuleName = "services"
									log.Source = "browser"
									log.UserDisplayName = user.DisplayName
									log.UserName = user.LoginName
									log.Datetime = time.Now()
									// fmt.Println(log.String())
									tools.WriteSystemLog(&log)

								} else {
									fmt.Println(err.Error())
								}
							}
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
	}
}
func Cors(ctx iris.Context) {
	ctx.Text("")
}

// 绑定 UWB 设备相关WEB服务
func RegisterUWBMgtServices(app *iris.Application) {
	app.Post("/uwbtag", services.WSCreateUWBTag)      // 新增 UWB 标签信息
	app.Get("/uwbtag", services.WSGetUWBTag)          // 获取 UWB 标签信息
	app.Get("/uwbtag/query", services.WSQueryUWBTags) // 查询 UWB 标签信息
	app.Delete("/uwbtag", services.WSDeleteUWBTags)   // 删除 UWB 标签信息

	app.Post("/uwbbasestation", services.WSCreateUWBBaseStation)      // 创建UWB基站信息
	app.Get("/uwbbasestation/query", services.WSQueryUWBBaseStations) // 查询基站信息
	app.Delete("/uwbbasestation", services.WSDeleteUWBBaseStations)   // 删除基站信息
	app.Put("/uwbbasestation", services.WSUpdateUWBBaseStation)       // 更新基站信息
	app.Get("/uwbbasestation", services.WSGetUWBBaseStationByCode)    // 获取基站信息
}

// 绑定用户相关WEB服务
func RegisterSysUserServices(app *iris.Application) {

	app.Post("/sysuser", services.WSCreateUser)           //新增系统用户
	app.Delete("/sysuser", services.WSDeleteSysUser)      //删除系统用户
	app.Put("/sysuser", services.WSUpdateSysUser)         // 修改系统用户
	app.Get("/sysuser", services.WSGetSysUserByLoginName) //获取系统用户信息(通过登录名)

	app.Post("/login", services.WSLoginSystem)                                             //登录系统
	app.Get("/sysuser/listall", services.WSEnumSysUsers)                                   //列举系统用户
	app.Get("/sysuser/listall/fromcompanies", services.WSEnumSysUsersFromSportsCompanyIds) // 通过体育公司唯一编号集合获取所有下属系统用户

	app.Get("/sysuser/query", services.WSQuerySysUsers) // 查询系统用户列表
	app.Delete("/sysusers", services.WSDeleteSysUsers)  // 批量删除用户
}

// 绑定体育公司相关WEB服务
func RegisterSportsCompanyServices(app *iris.Application) {

	app.Post("/company", services.WSCreateSportsCompany)   //创建体育运动公司
	app.Delete("/company", services.WSDeleteSportsCompany) // 删除体育运动公司

	app.Delete("/companies", services.WSDeleteSportsCompanies) // 批量删除体育运动公司

	app.Put("/company", services.WSUpdateSportsCompany)  // 更新修改体育运动公司
	app.Get("/company", services.WSGetSportsCompanyById) // 获取体育公司相信信息

	app.Get("/company/query", services.WSQueryCompanies)                        // 查询公司列表
	app.Get("/company/mgtgroup", services.WSEnumSportsCompaniesByGroupId)       // 根据公司管理组查询其下所有公司
	app.Put("/comgtgroup/joincompanies", services.WSJoinInSportCompanyMgtGroup) // 将公司加入到公司管理组
	app.Put("/company/relsites", services.WSRelSportsCompanyAndSites)           // 为体育公司管关联场地
	app.Get("/company/rightsuser", services.WSEnumSportsCompaniesByRightUser)   // 通过系统用户编号枚举他能够管理的公司

}

// 绑定游泳者相关WEB服务
func RegisterSwimmerServices(app *iris.Application) {
	app.Post("/swimmer", services.WSCreateSwimmer)                                     // 创建游泳者信息
	app.Get("/swimmer", services.WSGetSwimmersById)                                    // 获取游泳者信息
	app.Delete("/swimmer", services.WSDeleteSwimmers)                                  // 批量删除游泳者信息
	app.Put("/swimmer", services.WSUpdateSwimmer)                                      // 修改游泳者（会员）信息
	app.Get("/swimmer/query", services.WSQuerySwimmers)                                // 查询游泳者信息
	app.Put("/swimmer/setcompanies", services.WSSwimmerJoinInSportsCompanies)          // 将游泳者加入公司
	app.Put("/swimmer/setviplevel", services.WSSetSwimmerVIPLevel)                     // 设置游泳者会员等级
	app.Put("/swimmer/companyswimmer", services.WSUpdateCompanySwimmer)                // 修改公司会员信息
	app.Get("/swimmer/companyswimmer", services.WSGetCompanySwimmer)                   // 获取公司会员信息
	app.Get("/swimmer/enumcompanies", services.WSEnumSportsCompanySwimmersBySwimmerId) // 通过游泳者编号获取其所属公司集合（仅公司名和编号）

}

// 绑定场地相关WEB服务
func RegisterSiteServices(app *iris.Application) {
	app.Post("/site", services.WSCreateSite)                                  // 创建查询场地信息
	app.Get("/site", services.WSGetSiteById)                                  // 获取查询场地信息
	app.Delete("/site", services.WSDeleteSites)                               // 批量删除场地信息
	app.Put("/site", services.WSUpdateSite)                                   // 修改场地信息
	app.Get("/site/query", services.WSQuerySites)                             // 查询场地信息
	app.Put("/site/setusers", services.WSSetSiteUsers)                        // 设置场地用户集合
	app.Get("/site/enumusers", services.WSEnumSiteUsersBySiteId)              // 通过场地编号获取其所有用户集合
	app.Put("/site/setcompanies", services.WSSiteJoinInSportsCompanines)      // 将场地加入公司
	app.Get("/site/enumcompanies", services.WSEnumSportsCompanySitesBySiteId) // 通过场地编号获取其所属公司集合
	app.Get("/site/enumsites", services.WSEnumSiteUsersByUserId)              // 通过用户编号获取其所属用户场地关联信息集合
	app.Post("/site/envcalendar", services.WSCreateSiteEnvCalendar)           // 创建场地环境日历
	app.Get("/site/enumenvcalendars", services.WSEnumSiteEnvCalendars)        // 枚举日期段的场地环境信息
	app.Put("/site/envcalendar", services.WSUpdateSiteEnvCalendar)            // 更新场地环境日历
	app.Get("/site/swimmerreport", services.WSSiteSwimmerReport)              // 场地游泳者统计报表
	app.Post("/site/fence", services.WSCreateSiteFence)                       // 创建场地泳池电子围栏
	app.Put("/site/fence", services.WSUpdateSiteFence)                        // 更新场地泳池电子围栏
	app.Get("/site/fence", services.WSGetSiteFence)                           // 获取场地泳池电子围栏信息
	app.Get("/site/fence/enumcodes", services.WSEnumSiteFenceCodes)           // 枚举场地所有泳池电子围栏信息编码
	app.Get("/site/fence/enum", services.WSEnumSiteFences)                    // 枚举场地所有泳池电子围栏信息
}

// 注册公司管理组相关WEB服务
func RegisterSportsCompanyMgtGroup(app *iris.Application) {
	app.Post("/comgtgroup", services.WSCreateSportsCompanyGroup)           // 创建公司管理组
	app.Get("/comgtgroup/query", services.WSQuerySportsCompanyGroups)      // 查询公司管理组
	app.Delete("/comgtgroup", services.WSDeleteSportCompanyMgtGroups)      // 删除体育公司管理组
	app.Put("/comgtgroup", services.WSUpdateSportsCompanyMgtGroup)         // 更新体育公司管理组信息
	app.Put("/comgtgroup/users", services.WSSetSystemUsersToGroup)         // 设置系统用户到公司管理组
	app.Delete("/comgtgroup/users", services.WSRemoveSystemUsersFromGroup) // 将系统用户从公司管理组移除

}

// 绑定系统管理相关服务
func RegisterSystemManagement(app *iris.Application) {
	app.Post("/role", services.WSCreateSysRole)                         // 创建系统角色
	app.Get("/role", services.WSGetRoleByRoleId)                        // 获取系统角色信息
	app.Put("/role/joinfuncpage", services.WSSysRoleJoinInSysFuncPages) // 设置角色和功能接口绑定
	app.Get("/role/funcpages", services.WSEnumAllFuncPagesByRoleId)     // 根据角色查询其下所有功能页面
	app.Get("/role/query", services.WSQuerySysRoles)                    // 查询系统角色
	app.Put("/role", services.WSUpdateSysRole)                          // 更新角色信息

	app.Get("/dict", services.WSGetDictValues) // 获取字典值
	app.Put("/dict", services.WSSetSystemDict) // 设置字典
	app.Get("/dict/children", services.WSGetChildrenSystemDictsByParent)

	app.Post("/sysfuncpage", services.WSCreateSysFuncPage)            // 创建系统功能页面
	app.Delete("/sysfuncpage", services.WSDeleteSysFuncPage)          // 删除系统功能界面
	app.Get("/sysfuncpage/list", services.WSEnumSysFuncPages)         //列举所有系统功能页面
	app.Post("/systemlog/operationlog", services.WSWriteOperationLog) // 创建操作日志
	app.Get("/athorizationinfo", services.WSGetLoginInformation)      //获取登录信息

}

// 日程相关接口
func RegisterCalendarService(app *iris.Application) {
	// 根据日期范围和场地编号获取其所有游泳者的日程（包括：计划和入场情况）
	app.Get("/calendar/datescope", services.WSEnumSwimmerCalendarByDateScope)

	// 游泳者入场登记
	app.Put("/calendar/swimmersite/enter", services.WSSwimmerEnterToSite)

	// 游泳者出场签出
	app.Put("/calendar/swimmersite/exit", services.WSSwimmerExitFromSite)

	// 游泳教练配置游泳者训练周期
	app.Put("/calendar/swimmersite/plancycle", services.WSSwimmerPlanCycle)

	// 游泳教练取消训练计划
	app.Put("/calendar/swimmersite/cancleplaycycle", services.WSSwimmerCalendarPlanCancel)

}

// WIFIDOG 支持服务
func RegisterWifiDogService(app *iris.Application) {
	app.Get("/wifidog/ping", DevicePing)
	app.Get("/wifidog/auth", AuthScriptPathFragment)
	app.Get("/wifidog/login", LoginScriptPathFragment)
	app.Get("/wifidog/portal", PortalScriptPathFragment)
	app.Get("/wifidog/msg", MsgScriptPathFragment)
}

func RegisterServices(app *iris.Application) {
	app.Get("/index", Index)
	RegisterSysUserServices(app)
	RegisterSportsCompanyServices(app)
	RegisterSwimmerServices(app)
	RegisterSiteServices(app)
	RegisterSportsCompanyMgtGroup(app)
	RegisterUWBMgtServices(app)
	RegisterSystemManagement(app)
	RegisterCalendarService(app)
	RegisterWifiDogService(app)
	// app.Get("/siteowners/list", services.WSEnumSiteOwners) //列举场地负责人
	// app.Put("/siteowners/setowners", services.WSSetSiteOwners) //设置场地负责人

}
