package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func GetUWBBaseStationByCode(code string) (entities.UWBBaseStation, error) {
	var baseStation entities.UWBBaseStation
	result := Database.First(&baseStation, "code=?", code)
	// InitDatabase()
	return baseStation, result.Error
}

func CreateUWBBaseStation(station *entities.UWBBaseStation) error {

	return Database.Create(station).Error
}

func DeleteUWBBaseStations(codes []string) error {
	var station entities.UWBBaseStation
	result := Database.Delete(&station, codes)
	return result.Error
}

func UpdateUWBBaseStation(station *entities.UWBBaseStation) error {
	result := Database.Model(station).Where("code=?", station.Code).UpdateColumns(station)
	return result.Error
}
func UpdateUWBBaseStationSiteDisplayName(siteId string, displayName string) error {
	var station entities.UWBBaseStation
	result := Database.Model(&station).Where("site_id=?", siteId).UpdateColumn("site_display_name", displayName)
	return result.Error
}
func GetUWBBaseStationCount(queryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var station entities.UWBBaseStation
	var result *gorm.DB
	if queryCodition.LikeValue != "" {
		result = Database.Model(&station).Where(`site_display_name LIKE ? OR code LIKE ? OR description LIKE ?`,
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%").Count(&count)
	} else {
		result = Database.Model(&station).Count(&count)
	}
	return count, result.Error
}

func QueryUWBBaseStations(queryCodition entities.QueryCondition) ([]entities.UWBBaseStation, error) {
	var station entities.UWBBaseStation
	var stations []entities.UWBBaseStation
	var result *gorm.DB
	var selectFileds = ` * `
	if queryCodition.LikeValue != "" {
		result = Database.Model(&station).Select(selectFileds).
			Where(`site_display_name LIKE ? OR code LIKE ? OR description LIKE ?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%").
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&stations)
	} else {
		result = Database.Model(&station).Select(selectFileds).
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&stations)
	}
	return stations, result.Error
}
