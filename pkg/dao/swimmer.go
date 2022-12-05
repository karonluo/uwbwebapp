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
func GetSwimmersCount(queryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var swimmer entities.Swimmer
	var result *gorm.DB
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
			uwb_tag_code LIKE ? OR 
			description LIKE ?
		`,
				"%"+queryCodition.LikeValue+"%",
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
	return count, result.Error
}

func QuerySwimmers(queryCodition entities.QueryCondition) ([]entities.Swimmer, error) {
	var swimmer entities.Swimmer
	var swimmers []entities.Swimmer
	var result *gorm.DB
	var selectFields = `Id, display_name, uwb_tag_code, id_card_number, gender, age, address, cellphone, wechat, create_datetime, modify_datetime, Creator, Modifier`
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
				uwb_tag_code LIKE ? OR 
				description LIKE ?
			`,
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
			Find(&swimmers)
	} else {
		result = Database.Model(&swimmer).Select(selectFields).
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&swimmers)
	}
	return swimmers, result.Error
}

func UpdateSwimmer(swimmer entities.Swimmer) error {
	result := Database.Table("swimmers").Where("id=?", swimmer.Id).UpdateColumns(&swimmer)
	return result.Error
}
