package services

import (
	"uwbwebapp/pkg/biz"

	"github.com/kataras/iris/v12"
)

func WSGetDictValues(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var err error
	var result interface{}
	if ctx.FormValue("code") != "" {
		result, err = biz.GetSystemDictValue(ctx.FormValue("code"))
	} else {
		if ctx.FormValue("key") != "" {
			result, err = biz.GetSystemDictValues(ctx.FormValue("key"))
		}
	}
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 通过字典上级键或编码获取其下级字典信息
func WSGetChildrenSystemDictsByParent(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var err error
	var result interface{}
	if ctx.FormValue("code") != "" {
		result, err = biz.GetChildrenSystemDictsByParentCode(ctx.FormValue("code"))
	} else {
		if ctx.FormValue("key") != "" {
			result, err = biz.GetChildrenSystemDictsByParentKey(ctx.FormValue("key"))
		}
	}
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
