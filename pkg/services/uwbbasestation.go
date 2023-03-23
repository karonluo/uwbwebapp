package services

import (
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateUWBBaseStation(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	var station entities.UWBBaseStation
	var err error
	body, _ := ctx.GetBody()
	station, err = entities.UnmarshalUWBBaseStation(body)
	if err != nil {
		tools.ProcessError("services.WSCreateUWBBaseStation", "station, err = entities.UnmarshalUWBBaseStation(body)", err)
	} else {
		err = biz.CreateUWBBaseStation(&station)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = station.Code
		}

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetUWBBaseStationByCode(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	station, err := dao.GetUWBBaseStationByCode(ctx.FormValue("code"))
	if err != nil {
		tools.ProcessError("services.WSGetUWBBaseStationByCode", `station, err := dao.GetUWBBaseStationByCode(ctx.FormValue("code"))`, err)
		message.StatusCode = 500
		message.Message = err.Error()
	} else {
		message.Message = station
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 查询 UWB 基站信息
func WSQueryUWBBaseStations(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var station []entities.UWBBaseStation
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	if err != nil {
		tools.ProcessError(`services.WSQueryUWBBaseStations`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		station, pageCount, recordCount, query_err = biz.QueryUWBBaseStations(queryCondition)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQueryUWBBaseStations`, `sites, pageCount, recordCount, query_err = biz.QueryUWBBaseStations(queryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["stations"] = station
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 批量删除场地
func WSDeleteUWBBaseStations(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("codes") != "" {
		stationCodes := strings.Split(ctx.FormValue("codes"), ",")
		err := biz.DeleteUWBBaseStations(stationCodes)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteUWBBaseStations`,
				`err := biz.DeleteUWBBaseStations(stationCodes)`,
				err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 修改场地信息
func WSUpdateUWBBaseStation(ctx iris.Context) {
	var station entities.UWBBaseStation
	var err error
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	station, err = entities.UnmarshalUWBBaseStation(body)

	if err == nil {
		err = biz.UpdateUWBBaseStation(&station)
		if err != nil {
			tools.ProcessError("services.WSUpdateUWBBaseStation", "err = biz.UpdateUWBBaseStation(&station)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateUWBBaseStation", `station, err = entities.UnmarshalUWBBaseStation(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
