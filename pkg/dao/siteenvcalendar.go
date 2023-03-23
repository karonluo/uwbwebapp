package dao

import "uwbwebapp/pkg/entities"

func CreateSiteEnvCalendar(calendar *entities.SiteEnvCalendar) error {
	return Database.Table("site_env_calendars").Create(calendar).Error
}

// 获取场地日期段的环境日历
func EnumSiteEnvCalendars(bDate string, eDate string, siteId string) ([]entities.SiteEnvCalendar, error) {
	var calendars []entities.SiteEnvCalendar
	err := Database.Table("site_env_calendars").Where("date BETWEEN ? AND ? AND site_id = ?", bDate, eDate, siteId).Find(&calendars).Error
	return calendars, err

}

// 更新场地环境日历
func UpdateSiteEnvCalendar(calendar *entities.SiteEnvCalendar) error {
	return Database.Table("site_env_calendars").Where("date = ? AND site_id =?", calendar.Date, calendar.SiteID).UpdateColumns(calendar).Error
}

// 获取场地环境日历
func GetSiteEnvCalendar(siteId string, date string) (entities.SiteEnvCalendar, error) {
	var calendar entities.SiteEnvCalendar
	err := Database.Table("site_env_calendars").Where("date = ? AND site_id =?", date, siteId).First(&calendar).Error
	return calendar, err

}
