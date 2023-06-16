package services

import (
	"strings"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/kataras/iris/v12"
)

func WSCreateSportsCompanyGroup(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	group, err := entities.UnmarshalSportsCompanyMgtGroup(body)
	if err == nil {
		err = biz.CreateSportsCompanyGroup(&group)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
		} else {
			message.Message = group.Id
		}

	} else {
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSJoinInSportCompanyMgtGroup(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	co_ids := strings.Split(ctx.FormValue("co_ids"), ",")
	cmg_id := ctx.FormValue("cmg_id")
	err := biz.ClearSportCompanyMgtGroupCompanies(cmg_id)
	if err == nil {
		if len(co_ids) > 0 && co_ids[0] != "" {
			err = biz.JoinInSportCompanyMgtGroup(cmg_id, co_ids)
			if err != nil {
				tools.ProcessError(`services.WSQuerySportsCompanyGroups`,
					`err = biz.JoinInSportCompanyMgtGroup(cmg_id, co_ids)`,
					err)
			} else {
				err = biz.SetBoundCompanyMgtGroup(cmg_id, true)
				if err != nil {
					tools.ProcessError(`services.WSQuerySportsCompanyGroups`,
						`err = biz.SetBoundCompanyMgtGroup(cmg_id, true)`,
						err)
				}
			}
		} else {
			err = biz.SetBoundCompanyMgtGroup(cmg_id, false)
			if err != nil {
				tools.ProcessError(`services.WSQuerySportsCompanyGroups`,
					`err = biz.SetBoundCompanyMgtGroup(cmg_id, false)`,
					err)
			}
		}
	} else {
		tools.ProcessError(`services.WSQuerySportsCompanyGroups`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
	}

	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSQuerySportsCompanyGroups(ctx iris.Context) {
	message := WebServiceMessage{Message: "OK", StatusCode: 200}
	var queryCondition entities.QueryCondition
	var groups []entities.SportsCompanyMgtGroup
	var err, query_err error
	var pageCount, recordCount int64
	result := make(map[string]interface{})
	// body, _ := ctx.GetBody()
	// err := json.Unmarshal(body, &companyQueryCondition)
	queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))

	if err != nil {
		tools.ProcessError(`services.WSQuerySportsCompanyGroups`,
			`queryCondition, err = tools.GenerateQueryConditionFromWebParameters(ctx.FormValue("page_size"), ctx.FormValue("page_index"), ctx.FormValue("like_value"))`,
			err)
		message.Message = err.Error()
		message.StatusCode = 500
	} else {
		groups, pageCount, recordCount, query_err = biz.QuerySportCompanyMgtGroups(queryCondition)
		if query_err != nil {
			message.Message = query_err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSQuerySportsCompanyGroups`, `sites, pageCount, recordCount, query_err = biz.QuerySportCompanyMgtGroups(queryCondition)`, query_err)
		} else {
			result["page_count"] = pageCount
			result["record_count"] = recordCount
			result["groups"] = groups
			message.Message = result
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

// 批量删除公司管理组
func WSDeleteSportCompanyMgtGroups(ctx iris.Context) {
	message := WebServiceMessage{Message: true, StatusCode: 200}
	if ctx.FormValue("ids") != "" {
		groupIds := strings.Split(ctx.FormValue("ids"), ",")
		err := biz.DeleteSportCompanyMgtGroups(groupIds)
		if err != nil {
			message.Message = err.Error()
			message.StatusCode = 500
			tools.ProcessError(`services.WSDeleteSportCompanyMgtGroups`,
				`err := biz.DeleteSportCompanyMgtGroups(userIds)`,
				err)
		}
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
func WSUpdateSportsCompanyMgtGroup(ctx iris.Context) {
	var group entities.SportsCompanyMgtGroup
	var err error
	message := WebServiceMessage{Message: true, StatusCode: 200}
	body, _ := ctx.GetBody()
	group, err = entities.UnmarshalSportsCompanyMgtGroup(body)
	if err == nil {
		err = biz.UpdateSportsCompanyMgtGroup(&group)
		if err != nil {
			tools.ProcessError("services.WSSportsCompanyMgtGroup", "err = biz.UpdateSportsCompanyMgtGroup(&group)", err)
			message.Message = err.Error()
			message.StatusCode = 500
		}
	} else {
		tools.ProcessError("services.WSSportsCompanyMgtGroup", `group, err = entities.UnmarshalSportsCompanyMgtGroup(body)`, err)
		message.Message = err.Error()
		message.StatusCode = 500
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSSetSystemUsersToGroup(ctx iris.Context) {

	message := WebServiceMessage{Message: true, StatusCode: 200}
	userIds := strings.Split(ctx.FormValue("user_ids"), ",")
	groupId := ctx.FormValue("group_id")
	erok := biz.SetSystemUsersToGroup(userIds, groupId)
	message.Message = erok
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSRemoveSystemUsersFromGroup(ctx iris.Context) {

	message := WebServiceMessage{Message: true, StatusCode: 200}
	userIds := strings.Split(ctx.FormValue("user_ids"), ",")
	groupId := ctx.FormValue("group_id")
	erok := biz.RemoveSystemUsersFromGroup(userIds, groupId)
	message.Message = erok
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}

func WSEnumSystemUsersFromGroup(ctx iris.Context) {

	message := WebServiceMessage{Message: true, StatusCode: 200}

	groupId := ctx.FormValue("group_id")
	users, err := biz.EnumSystemUsersFromGroup(groupId)
	if err != nil {
		message.Message = err
		message.StatusCode = 500
	} else {
		message.Message = users
	}
	ctx.StatusCode(message.StatusCode)
	ctx.JSON(message)
}
