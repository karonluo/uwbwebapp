package services

import (
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateSite(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	var site entities.Site
	var err error
	body, _ := ctx.GetBody()
	site, err = site.UnmarshalSite(body)
	if err != nil {
		tools.ProcessError("services.WSCreateSite", "site, err = site.UnmarshalSite(body)", err)
	} else {
		err = biz.CreateSite(&site)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = site.Id
		}

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetSiteById(ctx iris.Context) {
	message := WebServiceMessage{StatusCode: 200, Message: true}
	site, err := dao.GetSiteById(ctx.FormValue("id"))
	if err != nil {
		tools.ProcessError("services.WSGetSiteById", `site, err := dao.GetSiteById(ctx.FormValue("id"))`, err)
		message.StatusCode = 500
		message.Message = err.Error()
	} else {
		message.Message = site
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 查询系统用户列表
func WSQuerySites(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var sites []entities.Site
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	if err != nil {
		tools.ProcessError(`services.WSQuerySites`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		sites, pageCount, recordCount, query_err = biz.QuerySites(queryCondition)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQuerySites`, `sites, pageCount, recordCount, query_err = biz.QuerySwimmers(queryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["sites"] = sites
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 批量删除场地
func WSDeleteSites(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("ids") != "" {
		userIds := strings.Split(ctx.FormValue("ids"), ",")
		err := biz.DeleteSites(userIds)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteSysUsers`,
				`err := biz.DeleteSites(userIds)`,
				err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 修改场地信息
func WSUpdateSite(ctx iris.Context) {
	var site entities.Site
	var err error
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	site, err = site.UnmarshalSite(body)
	if err == nil {
		err = biz.UpdateSite(&site)
		if err != nil {
			tools.ProcessError("services.WSUpdateSite", "err = biz.UpdateSite(&site)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateSite", `site, err = site.UnmarshalSite(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
