package services

import (
	"encoding/json"
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateUser(ctx iris.Context) {
	var message WebServiceMessage
	var sysuser entities.SysUser
	message.Message = "OK"
	message.StatusCode = 200
	bsysuser, _ := ctx.GetBody()
	err := json.Unmarshal(bsysuser, &sysuser)
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		userId, err := biz.CreateSysUser(sysuser)
		if err == nil {
			message.Message = userId

		} else {
			message.Message = err.Error()
			message.StatusCode = 500
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 通过登录名获取系统用户
func WSGetSysUserFromDBByLoginName(ctx iris.Context) {
	var message WebServiceMessage
	message.Message = "OK"
	message.StatusCode = 200
	sysuser := biz.GetSysUserFromDBByLoginName(ctx.FormValue("login_name"))
	message.Message = sysuser

	ctx.JSON(message)
}

func WSGetSysUserByLoginName(ctx iris.Context) {
	var message WebServiceMessage
	message.Message = true
	message.StatusCode = 200
	sysuser, res := biz.GetSysUserByLoginName(ctx.FormValue("login_name"))
	if res {
		message.Message = sysuser
	} else {
		message.Message = false
		message.StatusCode = 404
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 列举所有系统用户
func WSEnumSysUsers(ctx iris.Context) {
	sysusers, _ := biz.EnumSysUserFromDB()
	ctx.JSON(sysusers)
}

// 登录系统
func WSLoginSystem(ctx iris.Context) {
	var message WebServiceMessage
	message.Message = "OK"
	message.StatusCode = 200
	bytesBody, _ := ctx.GetBody()
	var user map[string]string
	json.Unmarshal(bytesBody, &user)
	login_name := user["login_name"]
	password := user["password"]
	if login_name != "" && password != "" {
		msg, result := biz.LoginSystem(login_name, password)
		if result {
			message.Message = msg
			message.StatusCode = 200
		} else {
			message.Message = msg
			message.StatusCode = 500
		}
	}
	ctx.JSON(message)

}

// 获取登录信息
func WSGetLoginInformation(ctx iris.Context) {
	var message WebServiceMessage
	var result interface{}
	message.Message = "OK"
	message.StatusCode = 200
	fieldName := ctx.FormValue("field")
	res, err := biz.GetLoginInformation(ctx.GetHeader("Authorization"), fieldName)
	if err != nil {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError(`services.WSGetLoginInformation`, `res, err := biz.GetLoginInformation(ctx.GetHeader("Authorization"), ctx.FormValue("field"))`, err, `pkg/services/sysuser.go`)
	} else {
		if fieldName == "sysuser" {
			var user entities.SysUser
			json.Unmarshal([]byte(res), &user)
			result = user
		} else if fieldName == "acl" {
			var acl interface{}
			json.Unmarshal([]byte(res), &acl)
			result = acl
		} else {
			result = res
		}
		message.Message = result
		message.StatusCode = 200
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 验证登录状态
func WSCheckLogin(ctx iris.Context) {
	var message WebServiceMessage
	message.Message = "OK"
	message.StatusCode = 200
	_, err := biz.CheckLogin(ctx.GetHeader("Authorization"))
	if err != nil {
		tools.ProcessError(`services.WSCheckLogin`, `_, err := biz.CheckLogin(ctx.GetHeader("Authorization"))`, err, `pkg/services/sysuser.go`)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.JSON(message)
}

// 删除系统用户
func WSDeleteSysUser(ctx iris.Context) {
	var message WebServiceMessage
	var result bool = true
	var err error
	message.StatusCode = 200
	if ctx.FormValue("id") != "" {
		result, err = biz.DeleteSysUser(ctx.FormValue("id"))
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 删除系统用户
func WSDeleteSysUsers(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("ids") != "" {
		userIds := strings.Split(ctx.FormValue("ids"), ",")
		err := biz.DeleteSysUsers(userIds)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteSysUsers`,
				`err := biz.DeleteSysUsers(userIds)`,
				err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 修改系统用户
func WSUpdateSysUser(ctx iris.Context) {
	var user entities.SysUser
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	err := json.Unmarshal(body, &user)
	if err == nil {
		err = biz.UpdateSysUser(user)
		if err != nil {
			tools.ProcessError("services.WSUpdateSysUser", "err = biz.UpdateSysUser(user)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateSysUser", `err:=json.Unmarshal(body, &user)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 查询系统用户列表
func WSQuerySysUsers(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var users []entities.SysUser
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))
	companyID := ctx.FormValue("company_id")
	if err != nil {
		tools.ProcessError(`services.WSQuerySysUsers`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		users, pageCount, recordCount, query_err = biz.QuerySysUsers(queryCondition, companyID)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQuerySysUsers`, `users, pageCount, recordCount, query_err = biz.QuerySysUsers(queryCondition, companyID)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["sys_users"] = users
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 列举公司所属所有系统用户
func WSEnumSysUsersFromSportsCompanyIds(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	users, err := biz.EnumSysUsersFromSportsCompanyIds(strings.Split(ctx.FormValue("company_ids"), ","))
	if err != nil {
		tools.ProcessError(`services.WSEnumSysUsersFromSiteIds`,
			`users, err := biz.EnumSysUsersFromSportsCompanyIds(strings.Split(ctx.FormValue("company_ids"), ","))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500

	} else {
		message.Message = users
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
