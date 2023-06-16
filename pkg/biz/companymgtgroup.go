package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func CreateSportsCompanyGroup(group *entities.SportsCompanyMgtGroup) error {
	group.Id = uuid.New().String()
	if group.Creator == "" {
		group.Creator = "admin"
	}
	if group.Modifier == "" {
		group.Modifier = "admin"
	}
	group.CreateDatetime = time.Now()
	group.ModifyDatetime = time.Now()
	return dao.CreateSportsCompanyGroup(group)
}

func JoinInSportCompanyMgtGroup(groupId string, coIds []string) error {

	return dao.JoinInSportCompanyMgtGroup(groupId, coIds)

}

func ClearSportCompanyMgtGroupCompanies(groupId string) error {
	return dao.ClearSportCompanyMgtGroupCompanies(groupId)
}

func QuerySportCompanyMgtGroups(queryCondition entities.QueryCondition) ([]entities.SportsCompanyMgtGroup, int64, int64, error) {
	var sites []entities.SportsCompanyMgtGroup
	dataRecordCount, err := dao.GetSportCompanyMgtGroupCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		sites, err = dao.QuerySportCompanyMgtGroups(queryCondition)
	}

	return sites, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func DeleteSportCompanyMgtGroups(groupIds []string) error {
	// for _, id := range groupIds {
	// 	dao.ClearSportCompanyMgtGroupCompanies(id)
	// }
	// 通过前台判断是否有公司关联到公司管理组中，需要先移除这些公司。
	return dao.DeleteSportCompanyMgtGroups(groupIds)
}

func SetBoundCompanyMgtGroup(groupId string, isBound bool) error {
	return dao.SetBoundCompanyMgtGroup(groupId, isBound)
}

func UpdateSportsCompanyMgtGroup(group *entities.SportsCompanyMgtGroup) error {
	tmpGroup, err := dao.GetSportsCompanyMgtGroupById(group.Id)
	if err == nil {
		group.CreateDatetime = tmpGroup.CreateDatetime
		group.Creator = tmpGroup.Creator
		group.ModifyDatetime = time.Now()
		group.IsBound = tmpGroup.IsBound
		if group.Modifier == "" {
			group.Modifier = "admin"
		}
		err = dao.UpdateSportsCompanyMgtGroup(group)
	}
	return err
}

// 设置系统用户进入公司组
func SetSystemUsersToGroup(userIds []string, groupId string) entities.EROK {
	var erok entities.EROK
	for _, userId := range userIds {
		err := dao.SetSystemUserToGroup(userId, groupId)
		if err != nil {
			erok.ErrList = append(erok.ErrList, userId)

		} else {
			erok.SuccessList = append(erok.SuccessList, userId)
		}
	}
	return erok
}

// 将系统用户从公司组移除
func RemoveSystemUsersFromGroup(userIds []string, groupId string) entities.EROK {
	var erok entities.EROK
	for _, userId := range userIds {
		err := dao.RemoveSystemUserFromGroup(userId, groupId)
		if err != nil {
			erok.ErrList = append(erok.ErrList, userId)

		} else {
			erok.SuccessList = append(erok.SuccessList, userId)
		}
	}
	return erok
}

func EnumSystemUsersFromGroup(groupId string) ([]entities.SysUser, error) {

	return dao.EnumSystemUsersFromGroup(groupId)
}
