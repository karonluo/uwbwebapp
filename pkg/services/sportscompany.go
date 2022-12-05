package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

// 新增体育公司
func WSCreateSportsCompany(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	bytesCompany, _ := ctx.GetBody()
	company, err := entities.UnmarshalSportsCompany(&bytesCompany)
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()

	} else {
		id, err := biz.CreateSportsCompany(&company)
		if err != nil {
			message.StatusCode = 500
			message.Message = err.Error()
		} else {
			message.Message = id
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 删除体育运动公司
func WSDeleteSportsCompanies(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("ids") != "" {
		compnayIds := strings.Split(ctx.FormValue("ids"), ",")
		err := biz.DeleteSportsCompanies(compnayIds)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 删除体育运动公司
func WSDeleteSportsCompany(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("id") != "" {
		company := entities.SportsCompany{Id: ctx.FormValue("id")}
		result, err := biz.DeleteSportsCompany(company)
		if err != nil {
			message.StatusCode = 500
			message.Message = err.Error()
		} else {
			if !result {
				message.StatusCode = 500
				message.Message = "未找到相关体育运动公司"
			} else {
				message.Message = true
			}
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 修改体育运动公司信息
func WSUpdateSportsCompany(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	bytesCompany, _ := ctx.GetBody()
	company, err := entities.UnmarshalSportsCompany(&bytesCompany)
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
	} else {
		err := biz.UpdateSportsCompany(&company)
		if err != nil {
			message.StatusCode = 500
			message.Message = err.Error()
		} else {
			message.Message = true
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetSportsCompanyById(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	result, err := biz.GetSportsCompanyById(ctx.FormValue("id"))
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500

	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 关联场地
func WSRelSportsCompanyAndSites(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	ids, _ := ctx.GetBody()
	rel := make(map[string]string)
	err := json.Unmarshal(ids, &rel)
	if err == nil {
		arrayIds := strings.Split(rel["site_ids"], ",")
		fmt.Println((arrayIds))
		if len(arrayIds) > 0 {
			for _, site_id := range arrayIds {
				biz.RelSportsCompanyAndSite(rel["company_id"], site_id)
			}
		}
	}
	ctx.JSON(message)
}

// 查询列表
func WSQueryCompanies(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var companyQueryCondition entities.QueryCondition
	var companies []entities.SportsCompany
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	companyQueryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	if err != nil {
		tools.ProcessError(`services.WSQueryCompanies`,
			`companyQueryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		companies, pageCount, recordCount, query_err = biz.QueryCompanies(companyQueryCondition)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQueryCompanies`, `companies, pageCount, query_err := biz.QueryCompanies(companyQueryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["sports_companies"] = companies
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
