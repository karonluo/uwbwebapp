package services

import (
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"

	"github.com/kataras/iris/v12"
)

func WSEnumSysFuncPages(ctx iris.Context) {
	var message WebServiceMessage
	pages, _ := biz.EnumSysFuncPages()
	message.Message = pages
	message.StatusCode = 200
	ctx.JSON(message)

}

func WSCreateSysFuncPage(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	var page entities.SysFuncPage
	var id string
	body, err := ctx.GetBody()
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		page, err = entities.UnmarshalSysFuncPage(body)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			id, err = biz.CreateSysFuncPage(&page)
			if err == nil {
				message.Message = id
			}
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
