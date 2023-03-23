package services

import (
	"encoding/json"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"

	"github.com/kataras/iris/v12"
)

func WSSetSystemDict(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var err error
	var dict entities.SysDict
	var body []byte
	body, err = ctx.GetBody()
	if err == nil {
		err = json.Unmarshal(body, &dict)
		if err == nil {
			err = biz.SetSystemDict(&dict)
		}
	}
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		message.Message = dict
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
func WSGetDictValues(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var err error
	var result interface{}
	if ctx.FormValue("code") != "" {
		result, err = biz.GetSystemDict(ctx.FormValue("code"))
	} else {
		if ctx.FormValue("key") != "" {
			result, err = biz.GetSystemDicts(ctx.FormValue("key"))
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
