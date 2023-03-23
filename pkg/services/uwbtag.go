package services

import (
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateUWBTag(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	var tag entities.UWBTag
	var err error
	body, _ := ctx.GetBody()
	tag, err = entities.UnmarshalUWBTag(body)
	if err != nil {
		tools.ProcessError("services.WSCreateUWBTag", "tag, err = entities.UnmarshalUWBTag(body)", err)
	} else {
		err = biz.CreateUWBTag(&tag)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSCreateUWBTag`, `code, err = biz.CreateUWBTag(&tag)`, err)
		} else {
			message.Message = tag.Code
		}

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetUWBTag(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	tag, err := biz.GetUWBTag(ctx.FormValue("code"))
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError(`services.WSGetUWBTag`, `tag, err := biz.GetUWBTag(ctx.FormValue("code"))`, err)
	} else {
		message.Message = tag
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)

}

func WSQueryUWBTags(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var tags []entities.UWBTag
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})

	// 处理查询条件 - 搜索条件
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	// 处理查询条件 - 检索条件 - 公司
	companyIds := ctx.FormValue("company_ids")
	var slice_companyIds []string
	if strings.Contains(companyIds, ",") {
		slice_companyIds = strings.Split(companyIds, ",")
	} else if companyIds != "" {
		slice_companyIds = append(slice_companyIds, companyIds)
	}
	// 处理查询条件 - 检索条件 - 是否绑定给了用户
	isBound := ctx.FormValue("isbound")
	var bIsBound interface{}
	if isBound == "true" {
		bIsBound = true
	} else if isBound == "false" {
		bIsBound = false
	} else {
		bIsBound = nil
	}

	if err != nil {
		tools.ProcessError(`services.WSQueryUWBTags`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		tags, pageCount, recordCount, query_err = biz.QueryUWBTags(queryCondition, slice_companyIds, bIsBound)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQueryUWBTags`,
				`tags, pageCount, recordCount, query_err = biz.QueryUWBTags(queryCondition, slice_companyIds, bIsBound)`,
				query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["uwb_tags"] = tags
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)

}

// 批量删除标签
func WSDeleteUWBTags(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("codes") != "" {
		tagCodes := strings.Split(ctx.FormValue("codes"), ",")
		err := biz.DeleteUWBTags(tagCodes)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteUWBTags`, `err := biz.WSDeleteUWBTags(tagCodes)`, err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
