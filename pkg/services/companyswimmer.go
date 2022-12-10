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
