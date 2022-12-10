package services

import (
	"encoding/json"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSEnumSportsCompanySwimmersBySwimmerId(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	result, err := biz.EnumSportsCompanySwimmersBySwimmerId(ctx.FormValue("swimmer_id"))
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.WSEnumSportsCompanySwimmersBySwimmerId", `result, err := biz.EnumSportsCompanySwimmersBySwimmerId(ctx.FormValue("swimmer_id"))`, err)
	} else {
		message.Message = result

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSwimmerJoinInSportsCompanies(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	var css []entities.CompanySwimmer
	var errs []error
	body, _ := ctx.GetBody()
	err := json.Unmarshal(body, &css)
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.WSSwimmerJoinInSportsCompanies", `err := json.Unmarshal(body, &css)`, err)
	} else {
		css, errs = biz.SwimmerJoinInSportsCompanies(css, ctx.FormValue("swimmer_id"))
		if len(errs) > 0 {
			result := make(map[string]interface{})
			result["CompanySwimmers"] = css
			result["errors"] = errs
			message.Message = result

		} else {
			message.Message = css
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSetSwimmerVIPLevel(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}

	err := biz.SetSwimmerVIPLevel(ctx.FormValue("swimmer_id"),
		ctx.FormValue("company_id"),
		ctx.FormValue("vip_level_dict_code"),
		ctx.FormValue("vip_level"),
		ctx.FormValue("modifier"))
	if err != nil {
		tools.ProcessError(`services.WSSetSwimmerVIPLevel`, `err := biz.SetSwimmerVIPLevel(ctx.FormValue("swimmer_id"),
	ctx.FormValue("company_id"),
	ctx.FormValue("vip_level_dict_code"),
	ctx.FormValue("vip_level"),
	ctx.FormValue("modifier"))`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)

}

func WSUpdateCompanySwimmer(ctx iris.Context) {
	var companySwimmer entities.CompanySwimmer
	var err error
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	companySwimmer, err = entities.UnmarshalCompanySwimmer(body)
	if err == nil {
		err = biz.UpdateCompanySwimmer(&companySwimmer)
		if err != nil {
			tools.ProcessError("services.WSUpdateCompanySwimmer", "err = biz.UpdateCompanySwimmer(&companySwimmer)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateCompanySwimmer", `companySwimmer, err = entities.UnmarshalCompanySwimmer(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetCompanySwimmer(ctx iris.Context) {
	var companySwimmer entities.CompanySwimmer
	var err error
	message := WebServiceMessage{Message: true, StatusCode: 200}

	companySwimmer, _, err = biz.GetCompanySwimmerByCompanyIDAndSwimmerID(ctx.FormValue("company_id"), ctx.FormValue("swimmer_id"))
	if err != nil {
		tools.ProcessError("services.WSGetCompanySwimmer", `companySwimmer, _, err = biz.GetCompanySwimmerByCompanyIDAndSwimmerID(ctx.FormValue("company_id"), ctx.FormValue("swimmer_id"))`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		message.Message = companySwimmer
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
