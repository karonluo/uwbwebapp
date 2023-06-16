package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func CreateSportsCompanyGroup(group *entities.SportsCompanyMgtGroup) error {
	return Database.Create(group).Error
}

// 枚举组内公司
func EnumSportsCompaniesByGroupId(groupId string) ([]entities.SportsCompany, error) {
	var compines []entities.SportsCompany
	result := Database.Table("rel_company_groups").Select("sports_companies.*").
		Joins("left join sports_companies on sports_companies.id = rel_company_groups.co_id").
		Where("rel_company_groups.co_group_id=?", groupId).Find(&compines)

	// Database.Table("go_service_info").Select("go_service_info.serviceId as service_id, go_service_info.serviceName as service_name, go_system_info.systemId as system_id, go_system_info.systemName as system_name").Joins("left join go_system_info on go_service_info.systemId = go_system_info.systemId where go_service_info.serviceId <> ? and go_system_info.systemId = ?", "xxx", "xxx").Scan(&results)

	return compines, result.Error
}

// 将体育运动公司加入到管理组
func JoinInSportCompanyMgtGroup(groupId string, coIds []string) error {
	var dataList []map[string]interface{}
	for _, coId := range coIds {
		data := make(map[string]interface{})
		data["co_id"] = coId
		data["co_group_id"] = groupId
		dataList = append(dataList, data)
	}

	return Database.Table("rel_company_groups").Create(&dataList).Error
}

// 清空管理组中的体育运动公司信息
func ClearSportCompanyMgtGroupCompanies(groupId string) error {
	return Database.Exec("DELETE FROM rel_company_groups WHERE co_group_id = ?", groupId).Error
}
func SetBoundCompanyMgtGroup(groupId string, isBound bool) error {
	return Database.Table("sports_company_mgt_groups").Where("id=?", groupId).UpdateColumn("is_bound", isBound).Error
}
func GetSportCompanyMgtGroupCount(queryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var group entities.SportsCompanyMgtGroup
	var result *gorm.DB

	if queryCodition.LikeValue != "" {

		result = Database.Model(&group).Where(`name like ? or description like ?`,
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%").Count(&count)

	} else {

		result = Database.Model(&group).Count(&count)
	}
	return count, result.Error
}

func QuerySportCompanyMgtGroups(queryCodition entities.QueryCondition) ([]entities.SportsCompanyMgtGroup, error) {
	var group entities.SportsCompanyMgtGroup
	var groups []entities.SportsCompanyMgtGroup
	var result *gorm.DB
	selectFields := `id, name, description, Modifier, Creator, modify_datetime, create_datetime, is_bound`
	if queryCodition.LikeValue != "" {

		result = Database.Model(&group).Select(selectFields).Where(`name LIKE ? or description LIKE ?`,
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%").
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&groups)

	} else {
		result = Database.Model(&group).
			Select(selectFields).
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&groups)
	}
	return groups, result.Error
}

func DeleteSportCompanyMgtGroups(groupIds []string) error {
	var group entities.SportsCompanyMgtGroup
	result := Database.Delete(&group, groupIds)
	return result.Error
}

func UpdateSportsCompanyMgtGroup(group *entities.SportsCompanyMgtGroup) error {
	return Database.Model(group).UpdateColumns(group).Error
}

func GetSportsCompanyMgtGroupById(groupId string) (entities.SportsCompanyMgtGroup, error) {
	var group entities.SportsCompanyMgtGroup
	result := Database.Model(group).Where("id=?", groupId).First(&group)
	return group, result.Error
}

func SetSystemUserToGroup(userId string, groupId string) error {
	return Database.Exec("INSERT INTO company_group_right_users (sys_user_id, co_group_id) VALUES(?, ?)", userId, groupId).Error
}

func RemoveSystemUserFromGroup(userId string, groupId string) error {
	return Database.Exec("DELETE FROM company_group_right_users WHERE sys_user_id = ? and co_group_id = ?", userId, groupId).Error
}

func EnumSystemUsersFromGroup(groupId string) ([]entities.SysUser, error) {
	var users []entities.SysUser
	sql := `SELECT u.id, u.login_name, u.display_name FROM company_group_right_users AS mgtu LEFT JOIN sys_users AS u ON mgtu.sys_user_id = u.id 
	WHERE mgtu.co_group_id = ?`
	err := Database.Raw(sql, groupId).Find(&users).Error
	return users, err
}
