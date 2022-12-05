package services

import (
	"uwbwebapp/pkg/biz"

	"github.com/kataras/iris/v12"
)

func WSEnumSysFuncPages(ctx iris.Context) {
	var message WebServiceMessage
	pages, _ := biz.EnumSysFuncPages()
	message.Message = pages
	message.StatusCode = 200
	ctx.JSON(message)

}
