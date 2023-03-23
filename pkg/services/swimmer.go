package services

import (
	"encoding/json"
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSDeleteSwimmers(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("ids") != "" {
		swimmerIds := strings.Split(ctx.FormValue("ids"), ",")
		err := biz.DeleteSwimmers(swimmerIds)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteSiwmmers`, `err := biz.DeleteSwimmers(swimmerIds)`, err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSCreateSwimmer(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var swimmer entities.Swimmer

	bytesSwimmer, _ := ctx.GetBody()
	err := json.Unmarshal(bytesSwimmer, &swimmer)
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
		tools.ProcessError(`services.WSCreateSwimmer`, `err := json.Unmarshal(bytesSwimmer, &swimmer)`, err)

	} else {
		id, err := biz.CreateSwimmer(swimmer)
		if err != nil {
			message.StatusCode = 500
			message.Message = err.Error()
			tools.ProcessError(`services.WSCreateSwimmer`, `id, err := biz.CreateSwimmer(swimmer)`, err)
		} else {
			message.Message = id
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 查询系统用户列表
func WSQuerySwimmers(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var swimmers []entities.Swimmer
	var err, query_err error
	var pageCount, recordCount int64
	var companyId string
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))
	companyId = ctx.FormValue("company_id")
	if err != nil {
		tools.ProcessError(`services.WSQuerySwimmers`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		swimmers, pageCount, recordCount, query_err = biz.QuerySwimmers(queryCondition, companyId)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQuerySwimmers`, `swimmers, pageCount, recordCount, query_err = biz.QuerySwimmers(queryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["swimmers"] = swimmers
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetSwimmersById(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	swimmer, err := biz.GetSwimmersById(ctx.FormValue("id"))
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError(`services.WSGetSwimmersById`, `swimmer, err := biz.GetSwimmersById(ctx.FormValue("id"))`, err)
	} else {
		message.Message = swimmer
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)

}

// 修改游泳者（会员）
func WSUpdateSwimmer(ctx iris.Context) {
	//var swimmer entities.Swimmer
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	swimmer, err := entities.UnmarshalSwimmer(body)
	if err == nil {
		err = biz.UpdateSwimmer(swimmer)
		if err != nil {
			tools.ProcessError("services.WSUpdateSwimmer", "err = biz.UpdateSwimmer(swimmer)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateSwimmer", `swimmer, err := entities.UnmarshalSwimmer(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
