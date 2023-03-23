package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func CreateSysRole(role entities.SysRole) error {
	result := Database.Create(&role)
	return result.Error
}

// 将系统角色和系统功能页面进行绑定
func SysRoleJoinInSysFuncPages(sysRoleId string, sysFuncPageIds []string) error {
	var rolePages []map[string]interface{}
	for _, funcPageId := range sysFuncPageIds {
		rolePage := make(map[string]interface{})
		rolePage["role_id"] = sysRoleId
		rolePage["func_page_id"] = funcPageId
		rolePages = append(rolePages, rolePage)
	}
	return Database.Table("rel_sysrole_sysfuncpages").Create(&rolePages).Error
}

func ClearSysRoleSysFuncPages(sysRoleId string) error {
	return Database.Exec("DELETE FROM rel_sysrole_sysfuncpages WHERE role_id = ?", sysRoleId).Error
}

// 根据角色编号获取其功能页面
func EnumAllFuncPagesByRoleId(sysRoleId string) ([]entities.SysFuncPage, error) {
	var pages []entities.SysFuncPage
	result := Database.Table("rel_sysrole_sysfuncpages").Select("sys_func_pages.*").
		Joins("left join sys_func_pages on sys_func_pages.id = rel_sysrole_sysfuncpages.func_page_id").
		Where("rel_sysrole_sysfuncpages.role_id=?", sysRoleId).Find(&pages)
	return pages, result.Error
}

func GetSysRoleById(sysRoleId string) (entities.SysRole, error) {
	var role entities.SysRole
	role.Id = sysRoleId
	err := Database.First(&role).Error
	return role, err
}

func GetSysRoleCount(companyQueryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var role entities.SysRole
	var result *gorm.DB
	if companyQueryCodition.LikeValue != "" {
		result = Database.Model(&role).Where(`display_name LIKE ? OR 
		description LIKE ?`,
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%").Count(&count)
	} else {
		result = Database.Model(&role).Count(&count)
	}
	return count, result.Error
}

func QuerySysRoles(queryCodition entities.QueryCondition) ([]entities.SysRole, error) {
	var role entities.SysRole
	var roles []entities.SysRole
	var result *gorm.DB
	selectFields := `id, display_name, description, Modifier, Creator, modify_datetime, create_datetime`
	if queryCodition.LikeValue != "" {

		result = Database.Model(&role).Select(selectFields).Where(`display_name LIKE ? or description LIKE ?`,
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%").
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&roles)

	} else {
		result = Database.Model(&role).
			Select(selectFields).
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&roles)
	}
	return roles, result.Error
}

func UpdateSysRole(role *entities.SysRole) error {
	return Database.Table("sys_roles").Where("id=?", role.Id).UpdateColumns(&role).Error
}
