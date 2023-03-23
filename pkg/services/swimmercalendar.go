package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

// 根据日期范围和场地编号获取其所有游泳者的日程（包括：计划和入场情况）
func WSEnumSwimmerCalendarByDateScope(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}

	strBDate := ctx.FormValue("bdate")
	strEDate := ctx.FormValue("edate")
	siteId := ctx.FormValue("site_id")

	bDate, err := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_TIME, strBDate)
	if err == nil {

		var eDate time.Time
		eDate, err = time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_TIME, strEDate)
		if err == nil {
			betweenBDate := tools.SplitDateTimeToAllDay(bDate)
			betweenEDate := tools.SplitDateTimeToAllDay(eDate)
			var betweenDate entities.BetweenDatetime
			betweenDate.BeginDatetime = betweenBDate.BeginDatetime
			betweenDate.EndDatetime = betweenEDate.EndDatetime
			calendars, _, errq := biz.EnumSwimmerCalendarByDateScope(siteId, betweenDate)
			if errq == nil {
				message.Message = calendars

			} else {
				err = errq
				tools.ProcessError(`services.WSEnumSwimmerCalendarByDateScope`,
					`calendars, _, errq := biz.EnumSwimmerCalendarByDateScope(siteId, betweenDate)`, errq)
			}

		} else {
			tools.ProcessError(`services.WSEnumSwimmerCalendarByDateScope`,
				`eDate, err = time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_TIME, strEDate)`, err)
		}
	} else {
		tools.ProcessError(`services.WSEnumSwimmerCalendarByDateScope`,
			`bDate, err := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_TIME, strBDate)`, err)
	}
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()

	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSwimmerEnterToSite(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	body, _ := ctx.GetBody()
	calendar, err := entities.UnmarshalSwimmerCalendar(body)
	var calendarIds []string
	var errs []string
	if err == nil {
		calendarIds, errs = biz.SwimmerEnterToSite(&calendar)
		if len(errs) == 0 {
			message.Message = calendarIds
		} else {
			msg := make(map[string]interface{}, 0)
			msg["errs"] = errs       // 未成功登记的错误信息。
			msg["ids"] = calendarIds // 被登记的ID。
			message.Message = msg
			for _, etmp := range errs {
				tools.ProcessError(`services.WSSwimmerEnterToSite`,
					`calendarIds, errs = biz.SwimmerEnterToSite(&calendar)`, fmt.Errorf(etmp))
			}
		}
	} else {
		tools.ProcessError(`services.WSSwimmerEnterToSite`,
			`calendar, err := entities.UnmarshalSwimmerCalendar(body)`, err)
	}
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()

	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSwimmerExitFromSite(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var ids, errs []string
	body, _ := ctx.GetBody()
	calendar, err := entities.UnmarshalSwimmerCalendar(body)
	if err == nil {
		ids, errs = biz.SwimmerExitFromSite(&calendar)
		if len(errs) == 0 {
			message.Message = ids
		} else {
			for _, e := range errs {
				tools.ProcessError(`services.WSSwimmerExitFromSite`,
					`ids, errs = biz.SwimmerExitFromSite(&calendar)`, fmt.Errorf(e))
			}
			var msg = make(map[string][]string, 0)
			msg["ids"] = ids
			msg["errs"] = errs
			message.Message = msg
		}
	} else {
		tools.ProcessError(`services.WSSwimmerExitFromSite`,
			`calendar, err := entities.UnmarshalSwimmerCalendar(body)`, err)
		message.StatusCode = 500
		message.Message = err.Error()
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSwimmerPlanCycle(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	body, _ := ctx.GetBody()
	var swimmers []entities.Swimmer
	err := json.Unmarshal(body, &swimmers)
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
		tools.ProcessError("services.WSSwimmerPlanCycle", `err := json.Unmarshal(body, &swimmers)`, err)
	} else {
		cycle := ctx.FormValue("cycle")
		if cycle == "0" {
			ids, errs := biz.SwimmerCalendarPlanToSite(swimmers, ctx.FormValue("site_id"), ctx.FormValue("bdate"), ctx.FormValue("btime"), ctx.FormValue("etime"))
			fmt.Println("一次性训练计划")
			msg := make(map[string][]string, 0)
			msg["ids"] = ids
			msg["errs"] = errs
			message.Message = msg
		} else {
			fmt.Println("分解周期性训练计划")
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSwimmerCalendarPlanCancel(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	ids := strings.Split(ctx.FormValue("ids"), ",")
	ids, errs := biz.SwimmerCalendarPlanCancel(ids)
	msg := make(map[string][]string, 0)
	msg["ids"] = ids
	msg["errs"] = errs
	message.Message = msg
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
