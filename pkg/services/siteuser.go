package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSEnumSiteUsersBySiteId(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}

	result, err := biz.EnumSiteUsersBySiteId(ctx.FormValue("site_id"))
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
		tools.ProcessError("services.WSEnumSiteUsersBySiteId", `result, err := biz.EnumSiteUsersBySiteId(ctx.FormValue("site_id"))`, err)
	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSEnumSiteUsersByUserId(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}

	result, err := biz.EnumSiteUsersByUserId(ctx.FormValue("sysuser_id"))
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
		tools.ProcessError("services.WSEnumSiteUsersByUserId", `result, err := biz.EnumSiteUsersByUserId(ctx.FormValue("sysuser_id"))`, err)
	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
