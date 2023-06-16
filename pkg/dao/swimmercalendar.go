package dao

import (
	"time"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"
)

func CreateSwimmerCalendar(calendar *entities.SwimmerCalendar) error {
	return Database.Create(calendar).Error
}

func GetSwimmerCalendarById(id string) (entities.SwimmerCalendar, error) {
	var calendar entities.SwimmerCalendar
	result := Database.Table("swimmer_calendars").Where("id=?", id).First(&calendar)
	return calendar, result.Error
}

// 枚举出所选日期所有游泳者信息，包括有入场记录的和有训练计划的游泳者
// 参数解释：
// traningBBDatetime 规划训练开始日期的查询开始日期
// traningEBDatetime 规划训练开始日期的查询结束日期
// enterBDatetime 进入场地日期的查询开始日期
// enterEDatetime 进入场地日期的查询结束日期
func EnumSwimmerCalendarByTraningOrSiteDateTime(
	traningBBDatetime time.Time, traningEBDatetime time.Time,
	enterBDatetime time.Time, enterEDatetime time.Time) ([]entities.SwimmerCalendar, error) {
	var results []entities.SwimmerCalendar
	var err error
	return results, err
}

// 根据日期范围(分解到入场日期和计划日期）和场地编号枚举其所有游泳者的日程
func EnumSwimmerCalendarByDateScope(siteId string, date entities.BetweenDatetime) ([]entities.SwimmerCalendar, int64, error) {
	var results []entities.SwimmerCalendar

	res := Database.Raw(`SELECT * FROM swimmer_calendars WHERE 
	((enter_datetime between ? and ?) or (training_begin_datetime between ? and ?))

	 AND site_id=?`,
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		siteId).Order("modify_datetime DESC").Find(&results)

	return results, res.RowsAffected, res.Error
}

// 根据日期范围(分解到入场日期和计划日期）和场地编号、游泳者编号枚举其所有日程
func EnumSwimmerCalendarBySwimmerDateScope(siteId string, swimmerId string, date entities.BetweenDatetime) ([]entities.SwimmerCalendar, int64, error) {
	var results []entities.SwimmerCalendar

	res := Database.Raw(`SELECT * FROM swimmer_calendars WHERE 
		( (enter_datetime between ? and ?) or (training_begin_datetime between ? and ?) )
	
		 AND site_id=? and swimmer_id=?`,
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		siteId, swimmerId).Order("modify_datetime DESC").Find(&results)

	return results, res.RowsAffected, res.Error
}

func UpdateSwimmerCalendar(calendar *entities.SwimmerCalendar) error {
	return Database.Table("swimmer_calendars").Where(`id=?`, calendar.ID).UpdateColumns(&calendar).Error

}

func CanclePlanCalendar(id string) error {
	var data map[string]time.Time = make(map[string]time.Time)
	data["training_begin_datetime"], _ = time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "0001-01-01 00:00:00")
	data["training_end_datetime"], _ = time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "0001-01-01 00:00:00")
	return Database.Table("swimmer_calendars").Where("id=?", id).UpdateColumns(data).Error

}

func DeleteSwimmerCalendar(calendar *entities.SwimmerCalendar) error {
	return Database.Delete(calendar).Error
}

// 枚举所有在场的游泳者
func EnumAllSwimmerCalendarsOnSite() ([]entities.SwimmerCalendar, error) {
	var calendars []entities.SwimmerCalendar
	err := Database.Table("swimmer_calendars").Where("exit_datetime = '0001-01-01 00:00:00' AND enter_datetime <> '0001-01-01 00:00:00'").Find(&calendars).Error
	return calendars, err

}

// 枚举所有指定场地在场的游泳者
func EnumAllSwimmerCalendarsOnSiteBySiteId(siteId string) ([]entities.SwimmerCalendar, error) {
	var calendars []entities.SwimmerCalendar
	err := Database.Table("swimmer_calendars").Where("site_id=? AND exit_datetime = '0001-01-01 00:00:00' AND enter_datetime <> '0001-01-01 00:00:00'", siteId).Find(&calendars).Error
	return calendars, err

}

func EnumSwimmerCalendarReportForSwimmer(swimmerId string, date entities.BetweenDatetime) ([]entities.SwimmerCalendar, error) {
	var results []entities.SwimmerCalendar

	res := Database.Raw(`SELECT * FROM swimmer_calendars WHERE 
		(enter_datetime between ? and ?) AND swimmer_id=?`,
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.BeginDatetime.Format(`2006-01-02 15:04:05.000000000`),
		date.EndDatetime.Format(`2006-01-02 15:04:05.000000000`),
		swimmerId).Order("modify_datetime DESC").Find(&results)

	return results, res.Error
}

// 通过游泳者唯一编号和未出场状态获取日程信息
func GetSwimmerCalendarBySwimmerIdAndNoExitSite(swimmerId string) (entities.SwimmerCalendar, error) {
	var result entities.SwimmerCalendar
	r := Database.Model(result).Where("swimmer_id = ? AND enter_datetime <> '0001-01-01 00:00:00' AND exit_datetime = '0001-01-01 00:00:00'", swimmerId).First(&result)
	return result, r.Error

}
