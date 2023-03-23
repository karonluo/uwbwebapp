package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"

	"github.com/kataras/iris/v12"
)

// 场地游泳者统计报表
func WSSiteSwimmerReport(ctx iris.Context) {
	var report entities.SiteSwimmerReport
	message := WebServiceMessage{StatusCode: 200, Message: true}

	var err error
	siteId := ctx.FormValue("site_id")
	report, err = biz.SiteSwimmerReport(siteId)
	if err == nil {
		message.Message = report
	} else {
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)

}
