package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"gorm.io/gorm"
)

func CreateSysUser(sysuser entities.SysUser) string {
	Database.Create(&sysuser)
	return sysuser.Id

}

func CheckHaveReSysUser(sysuser entities.SysUser) bool {
	var result []entities.SysUser
	var res bool = false
	Database.Raw(`SELECT login_name FROM sys_users WHERE cellphone=? OR email=? OR login_name=? OR id_card_number=? LIMIT 1;`,
		sysuser.Cellphone, sysuser.Email, sysuser.LoginName, sysuser.IDCardNumber).Scan(&result)
	if len(result) > 0 {
		res = true
	}
	return res
}

func EnumSysUserFromDB() ([]entities.SysUser, int64) {

	var sysusers []entities.SysUser
	result := Database.Find(&sysusers)
	return sysusers, result.RowsAffected
}

func GetUserFromRedisByLoginName(login_name string) (entities.SysUser, bool) {
	var sysuser entities.SysUser
	var result bool = true
	ctx := context.Background()
	res := RedisDatabase.Get(ctx, "sysuser_"+login_name)
	if res != nil {
		err := json.Unmarshal([]byte(res.Val()), &sysuser)
		if err != nil {
			result = false
			tools.ProcessError(`dao.GetUserFromRedisByLoginName`, `json.Unmarshal([]byte(res.Val()), &sysuser)`, err, `pkg/dao/sysuser.go`)

		}
	} else {
		result = false
	}
	return sysuser, result
}

func GetUserFromDBByLoginName(login_name string) (entities.SysUser, bool, error) {
	var result entities.SysUser
	var bres bool = false
	res := Database.Where(("login_name=?"), login_name).Find(&result)
	if res.RowsAffected == 1 {
		bres = true
	}
	return result, bres, res.Error
}

func DeleteSysUser1(id string) bool {
	var result bool = true
	var user entities.SysUser
	user.Id = id
	res := Database.Exec("DELETE FROM sys_users WHERE id=?", id)
	if res.Error != nil {
		result = false
		fmt.Printf("删除系统用户: %s 失败, 原因是: %s \r\n", id, res.Error.Error())
	}
	if res.RowsAffected == 0 {
		result = false
		fmt.Printf("删除系统用户: %s 失败, 原因是: %s \r\n", id, "未找到该系统用户")
	}
	return result
}

func DeleteSysUsers(userIds []string) error {
	var user entities.SysUser
	result := Database.Delete(&user, userIds)
	return result.Error
}

func DeleteSysUser(compnay entities.SysUser) (bool, error) {
	result := Database.Delete(&compnay)
	var res bool = true
	if result.RowsAffected == 0 {
		res = false
	}
	return res, result.Error
}

func GetSysUserCount(queryCodition entities.QueryCondition, companyID string) (int64, error) {
	var count int64
	var user entities.SysUser
	var result *gorm.DB
	if queryCodition.LikeValue != "" {
		if companyID != "" {
			result = Database.Model(&user).Where(
				`(login_name LIKE ? OR
			display_name LIKE ? OR
			cellphone LIKE ? OR
			email LIKE ? OR
			wechat LIKE ? OR
			id_card_number LIKE ? OR
			sports_company_name LIKE ? OR
			qq LIKE ?) AND sports_company_id=?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				companyID).Count(&count)
		} else {
			result = Database.Model(&user).Where(
				`login_name LIKE ? OR
			display_name LIKE ? OR
			cellphone LIKE ? OR
			email LIKE ? OR
			wechat LIKE ? OR
			id_card_number LIKE ? OR
			sports_company_name LIKE ? OR
			qq LIKE ?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%").Count(&count)
		}
	} else {
		if companyID != "" {
			result = Database.Model(&user).Where("sports_company_id", companyID).Count(&count)
		} else {
			result = Database.Model(&user).Count(&count)
		}
	}
	return count, result.Error
}

func QuerySysUsers(queryCodition entities.QueryCondition, companyID string) ([]entities.SysUser, error) {
	var user entities.SysUser
	var users []entities.SysUser
	var result *gorm.DB
	selectFields := `id, sports_company_name, id_card_number, login_name, display_name, cellphone,email, wechat, qq, Modifier, Creator, modify_datetime, create_datetime`
	if queryCodition.LikeValue != "" {
		if companyID != "" {
			result = Database.Model(&user).Select(selectFields).Where(
				`(login_name LIKE ? OR display_name LIKE ? OR cellphone LIKE ? OR email LIKE ? OR wechat LIKE ? OR
			id_card_number LIKE ? OR
			sports_company_name LIKE ? OR
			qq LIKE ?) AND sports_company_id = ?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%", companyID).
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&users)
		} else {
			result = Database.Model(&user).Select(selectFields).Where(
				`login_name LIKE ? OR display_name LIKE ? OR cellphone LIKE ? OR email LIKE ? OR wechat LIKE ? OR
			id_card_number LIKE ? OR
			sports_company_name LIKE ? OR
			qq LIKE ?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%").
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&users)
		}

	} else {
		if companyID == "" {
			result = Database.Model(&user).
				Select(selectFields).
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&users)
		} else {
			result = Database.Model(&user).
				Select(selectFields).Where("sports_company_id=?", companyID).
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&users)
		}
	}
	return users, result.Error
}

func GetSysUserFromDBById(sysUserId string) (entities.SysUser, error) {
	var user entities.SysUser
	result := Database.Table("sys_users").Where("id=?", sysUserId).First(&user)
	return user, result.Error

}
func UpdateSysUserPassword(sysUserId string, newPassword string) error {
	result := Database.Table("sys_users").Where("id=?", sysUserId).UpdateColumn("passwd_md5", tools.SHA1(newPassword))
	return result.Error
}
func UpdateSysUser(user entities.SysUser) error {
	result := Database.Table("sys_users").Where("id=?", user.Id).UpdateColumns(&user)
	return result.Error
}

func UpdateSysUserCompanyName(company *entities.SportsCompany) error {
	result := Database.Table("sys_users").
		Where("sports_company_id=?", company.Id).
		UpdateColumn("sports_company_name", company.Name)

	return result.Error
}

// 通过体育公司唯一编号集合获取其所有下属系统用户
func EnumSysUsersFromSportsCompanyIds(siteIds []string) ([]entities.SysUser, error) {
	var users []entities.SysUser
	selectFields := `id, sports_company_name, id_card_number, login_name, display_name, cellphone,email, wechat, qq, Modifier, Creator, modify_datetime, create_datetime`
	result := Database.Table("sys_users").Select(selectFields).Order("modify_datetime DESC, display_name ASC").
		Where("sports_company_id in ?", siteIds).Find(&users)
	return users, result.Error
}
