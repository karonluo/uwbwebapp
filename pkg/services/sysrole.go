package services

import (
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateSysRole(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	bytesRole, _ := ctx.GetBody()
	role, err := entities.UnmarshalSysRole(bytesRole)

	if err != nil {
		tools.ProcessError(`services.WSCreateSysRole`, `role, err := entities.UnmarshalSysRole(bytesRole)`, err, `pkg/services/sysrole.go`)
		message.StatusCode = 500
		message.Message = err.Error()

	} else {
		id, err := biz.CreateSysRole(role)
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

func WSSysRoleJoinInSysFuncPages(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	page_ids := strings.Split(ctx.FormValue("page_ids"), ",")
	role_id := ctx.FormValue("role_id")
	err := biz.ClearSysRoleSysFuncPages(role_id)
	if err == nil {
		if len(page_ids) > 0 && page_ids[0] != "" {
			err = biz.SysRoleJoinInSysFuncPages(role_id, page_ids)
			if err != nil {
				message.StatusCode = 500
				message.Message = err.Error()
			}
		}
	} else {
		message.StatusCode = 500
		message.Message = err.Error()
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSEnumAllFuncPagesByRoleId(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}

	result, err := biz.EnumAllFuncPagesByRoleId(ctx.FormValue("role_id"))
	if err != nil {
		message.StatusCode = 500
		message.Message = err.Error()
		tools.ProcessError("services.WSEnumAllFuncPagesByRoleId", `result, err := biz.WSEnumAllFuncPagesByRoleId(ctx.FormValue("role_id"))`, err)
	} else {
		message.Message = result
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSGetRoleByRoleId(ctx iris.Context) {
	var message = WebServiceMessage{Message: true, StatusCode: 200}
	message.Message = true
	message.StatusCode = 200
	sysRole, err := biz.GetSysRoleById(ctx.FormValue("role_id"))
	if err == nil {
		message.Message = sysRole
	} else {
		message.Message = err.Error()
		message.StatusCode = 500
		tools.ProcessError("services.WSGetRoleByRoleId", `sysRole, err := biz.GetSysRoleById(ctx.FormValue("role_id"))`, err)
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 查询列表
func WSQuerySysRoles(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var roles []entities.SysRole
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	if err != nil {
		tools.ProcessError(`services.WSQuerySysRoles`,
			`companyQueryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		roles, pageCount, recordCount, query_err = biz.QuerySysRoles(queryCondition)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQuerySysRoles`, `roles, pageCount, recordCount, query_err = biz.QuerySysRoles(queryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["sys_roles"] = roles
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 修改系统用户
func WSUpdateSysRole(ctx iris.Context) {
	var role entities.SysRole
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	role, err := entities.UnmarshalSysRole(body)
	if err == nil {
		err = biz.UpdateSysRole(&role)
		if err != nil {
			tools.ProcessError("services.WSUpdateSysRole", "err = biz.UpdateSysRole(&role)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSUpdateSysRole", `role, err := entities.UnmarshalSysRole(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
