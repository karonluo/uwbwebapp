package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

// 列举场地所有负责人
func WSEnumSiteOwners(ctx iris.Context) {
	if ctx.FormValue("site_id") != "" {
		ctx.JSON(biz.EnumSiteOwners(ctx.FormValue("site_id")))

	}
}

// 设置场地负责人
func WSSetSiteOwners(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	if ctx.FormValue("site_id") != "" && ctx.FormValue("sys_user_ids") != "" && ctx.FormValue("job_title") != "" {
		res, clear_msg := biz.ClearSiteOwners(ctx.FormValue("site_id"))
		if res {
			err := biz.SetSiteOwners(ctx.FormValue("site_id"), ctx.FormValue("sys_user_ids"), ctx.FormValue("job_title"))
			if err != nil {
				tools.ProcessError("services.WSSetSiteOwners", `biz.SetSiteOwners(ctx.FormValue("site_id"), ctx.FormValue("sys_user_ids"), ctx.FormValue("job_title"))`, err)
				message.Message = err.Error()
				message.StatusCode = 500
			}
		} else {
			message.StatusCode = 500
			message.Message = clear_msg
		}
	} else {
		message.StatusCode = 400
		message.Message = "参数不正确"

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// Web Service Message
type WebServiceMessage struct {
	StatusCode int
	Message    interface{}
}
