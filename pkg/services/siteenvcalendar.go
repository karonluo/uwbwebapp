package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

// 创建场地环境日程
func WSCreateSiteEnvCalendar(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	var calendar entities.SiteEnvCalendar
	var err error
	body, _ := ctx.GetBody()
	calendar, err = entities.UnmarshalSiteEnvCalendar(body)
	if err != nil {
		tools.ProcessError("services.WSCreateSiteEnvCalendar", "calendar, err = entities.UnmarshalSiteEnvCalendar(body)", err)
	} else {
		err = biz.CreateSiteEnvCalendar(&calendar)
		if err != nil {
			tools.ProcessError("services.WSCreateSiteEnvCalendar", "err = biz.CreateSiteEnvCalendar(&calendar)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = calendar
		}

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 枚举场地环境日期段
func WSEnumSiteEnvCalendars(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	bdate := ctx.FormValue("bdate")
	edate := ctx.FormValue("edate")
	siteId := ctx.FormValue("site_id")
	cals, err := biz.EnumSiteEnvCalendars(bdate, edate, siteId)

	if err != nil {
		tools.ProcessError("services.WSEnumSiteEnvCalendars", "cals, err := biz.EnumSiteEnvCalendars(bdate, edate, siteId)", err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		message.Message = cals
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 更新场地环境日历
func WSUpdateSiteEnvCalendar(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	var calendar entities.SiteEnvCalendar
	var err error
	body, _ := ctx.GetBody()
	calendar, err = entities.UnmarshalSiteEnvCalendar(body)
	if err != nil {
		tools.ProcessError("services.WSUpdateSiteEnvCalendar", "calendar, err = entities.UnmarshalSiteEnvCalendar(body)", err)
	} else {
		err = biz.UpdateSiteEnvCalendar(&calendar)
		if err != nil {
			tools.ProcessError("services.WSUpdateSiteEnvCalendar", "err = biz.UpdateSiteEnvCalendar(&calendar)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = calendar
		}

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
