package services

import (
	"encoding/json"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSEnumSportsCompanySitesBySiteId(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	result, err := biz.EnumSportsCompanySitesBySiteId(ctx.FormValue("site_id"))
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.EnumSportsCompanySitesBySiteId", `result, err := biz.EnumSportsCompanySitesBySiteId(ctx.FormValue("site_id"))`, err)
	} else {
		message.Message = result

	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSiteJoinInSportsCompanines(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	var css []entities.CompanySite
	var errs []error
	body, _ := ctx.GetBody()
	err := json.Unmarshal(body, &css)
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.WSSiteJoinInSportsCompanines", `err := json.Unmarshal(body, &css)`, err)
	} else {
		css, errs = biz.SiteJoinInSportsCompanies(css, ctx.FormValue("site_id"))
		if len(errs) > 0 {
			result := make(map[string]interface{})
			result["CompanySites"] = css // 返回成功加入的公司
			result["errors"] = errs      // 返回错误加入的公司
			message.Message = result

		} else {
			message.Message = css
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
