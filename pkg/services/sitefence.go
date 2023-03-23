package services

import (
	"encoding/json"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateSiteFence(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	var fence entities.SiteFence
	err := json.Unmarshal(body, &fence)
	if err == nil {
		err = biz.CreateSiteFence(&fence)
		if err == nil {
			message.Message = fence

		} else {
			tools.ProcessError("services.WSCreateSiteFence", `err = biz.CreateSiteFence(&fence)`, err)
		}

	} else {
		tools.ProcessError("services.WSCreateSiteFence", `err:=json.Unmarshal(body, &fence)`, err)
	}
	if err != nil {
		message.Message = err
		message.StatusCode = 500
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSUpdateSiteFence(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	var fence entities.SiteFence
	err := json.Unmarshal(body, &fence)
	if err == nil {
		err = biz.UpdateSiteFence(&fence)
		if err == nil {
			message.Message = fence

		} else {
			tools.ProcessError("services.WSUpdateSiteFence", `err = biz.UpdateSiteFence(&fence)`, err)
		}

	} else {
		tools.ProcessError("services.WSUpdateSiteFence", `err:=json.Unmarshal(body, &fence)`, err)
	}
	if err != nil {
		message.Message = err
		message.StatusCode = 500
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSEnumSiteFences(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	siteId := ctx.FormValue("site_id")
	fences, err := biz.EnumSiteFences(siteId)
	if err == nil {
		message.Message = fences

	} else {
		message.Message = err
		message.StatusCode = 500
		tools.ProcessError("services.WSEnumSiteFences", `fences,err:=biz.EnumSiteFences(siteId)`, err)

	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
func WSEnumSiteFenceCodes(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	siteId := ctx.FormValue("site_id")
	codes, err := biz.EnumSiteFenceCodes(siteId)
	if err == nil {
		message.Message = codes

	} else {
		message.Message = err
		message.StatusCode = 500
		tools.ProcessError("services.WSEnumSiteFences", `fences,err:=biz.EnumSiteFences(siteId)`, err)

	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
func WSGetSiteFence(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	siteId := ctx.FormValue("site_id")
	code := ctx.FormValue("code")
	fence, err := biz.GetSiteFence(siteId, code)
	if err == nil {
		message.Message = fence
	} else {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.WSGetSiteFence", `fence, err := biz.GetSiteFence(siteId, code)`, err)
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
