package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func CreateSwimmer(swimmer entities.Swimmer) error {
	result := Database.Create(&swimmer)
	return result.Error
}

func DeleteSwimmers(ids []string) error {
	var swimmer entities.Swimmer
	result := Database.Delete(&swimmer, ids)
	return result.Error
}
func GetSwimmersById(id string) (entities.Swimmer, error) {
	var swimmer entities.Swimmer
	result := Database.Where("id=?", id).First(&swimmer)
	return swimmer, result.Error
}
func GetSwimmersCount(queryCodition entities.QueryCondition, companyId string) (int64, error) {
	var count int64
	var swimmer entities.Swimmer
	var result *gorm.DB
	if companyId == "" {
		if queryCodition.LikeValue != "" {
			result = Database.Model(&swimmer).
				Where(
					`
			display_name LIKE ? OR 
			gender LIKE ? OR 
			address LIKE ? OR 
			cellphone LIKE ? OR 
			wechat LIKE ? OR 
			id_card_number LIKE ? OR 
			description LIKE ?
		`,
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%",
					"%"+queryCodition.LikeValue+"%").
				Count(&count)
		} else {
			result = Database.Model(&swimmer).Count(&count)
		}
	} else {
		sql := `SELECT count(sw.display_name) FROM swimmers as sw LEFT JOIN company_swimmers as co ON co.swimmer_id = sw."id" where co.sports_company_id = ? `
		if queryCodition.LikeValue != "" {
			where := ` AND (
				sw.display_name LIKE ? OR 
				sw.gender LIKE ? OR 
				sw.address LIKE ? OR 
				sw.cellphone LIKE ? OR 
				sw.wechat LIKE ? OR 
				sw.id_card_number LIKE ? OR 
				sw.description LIKE ?)
				`
			result = Database.Raw(sql+where,
				companyId,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%").Count(&count)
		} else {
			result = Database.Raw(sql, companyId).Count(&count)
		}

	}
	return count, result.Error
}

func QuerySwimmers(queryCodition entities.QueryCondition, companyId string) ([]entities.Swimmer, error) {
	var swimmer entities.Swimmer
	var swimmers []entities.Swimmer
	var result *gorm.DB
	var selectFields = `Id, display_name, id_card_number, gender, age, address, cellphone, wechat, create_datetime, modify_datetime, Creator, Modifier`
	if companyId == "" {

		if queryCodition.LikeValue != "" {
			result = Database.Model(&swimmer).Select(selectFields).
				Where(
					`
				display_name LIKE ? OR 
				gender LIKE ? OR 
				address LIKE ? OR 
				cellphone LIKE ? OR 
				wechat LIKE ? OR 
				id_card_number LIKE ? OR 
				description LIKE ?
			`,
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
				Find(&swimmers)
		} else {
			result = Database.Model(&swimmer).Select(selectFields).
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&swimmers)
		}
	} else {
		sql := `
		SELECT sw.Id, sw.display_name, sw.id_card_number, sw.gender, sw.age, sw.address, sw.cellphone, sw.wechat, 
		sw.create_datetime, sw.modify_datetime, sw.Creator, sw.Modifier 
		FROM swimmers as sw LEFT JOIN company_swimmers as co ON co.swimmer_id = sw."id" where co.sports_company_id = ? `
		if queryCodition.LikeValue != "" {
			where := ` AND (
				sw.display_name LIKE ? OR 
				sw.gender LIKE ? OR 
				sw.address LIKE ? OR 
				sw.cellphone LIKE ? OR 
				sw.wechat LIKE ? OR 
				sw.id_card_number LIKE ? OR 
				sw.description LIKE ?)
				`
			result = Database.Raw(sql+where,
				companyId,
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
				Find(&swimmers)
		} else {
			result = Database.Raw(sql, companyId).
				Order("modify_datetime DESC").
				Limit(int(queryCodition.PageSize)).
				Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
				Find(&swimmers)
		}
	}
	return swimmers, result.Error
}

func UpdateSwimmer(swimmer entities.Swimmer) error {
	result := Database.Table("swimmers").Where("id=?", swimmer.ID).UpdateColumns(&swimmer)
	return result.Error
}

func ClearAllCompaniesFromSwimmer(swimmerId string) error {
	return Database.Exec("DELETE FROM company_swimmers WHERE swimmer_id = ?", swimmerId).Error

}

// 当游泳者姓名变更时需要变更所有涉及的数据表
func UpdateSwimmerDisplayNameRelTables(swimmerId string, swimmerDisplayName string) error {
	err := Database.Table("company_swimmers").Where("swimmer_id=?", swimmerId).UpdateColumn("swimmer_display_name", swimmerDisplayName).Error
	if err == nil {
		err = Database.Table("swimmer_calendar").Where("swimmer_id=?", swimmerId).UpdateColumn("swimmer_display_name", swimmerDisplayName).Error
	}
	return err

}

func EnumSiteSimmerForReport(siteId string) ([]entities.Swimmer, error) {
	var result []entities.Swimmer
	sql := `SELECT swimmers.id, gender, age, display_name FROM swimmers WHERE ID IN (SELECT swimmer_id FROM swimmer_calendars WHERE site_id=? AND 
	exit_datetime = '0001-01-01 00:00:00' AND enter_datetime < NOW() AND enter_datetime <> '0001-01-01 00:00:00')
	`
	err := Database.Raw(sql, siteId).Find(&result).Error
	return result, err
}
