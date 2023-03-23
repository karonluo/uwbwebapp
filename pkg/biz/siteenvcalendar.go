package biz

import (
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func CreateSiteEnvCalendar(calendar *entities.SiteEnvCalendar) error {

	calendar.CreateDatetime = time.Now()
	calendar.ModifyDatetime = time.Now()
	if calendar.Creator == "" {
		calendar.Creator = "admin"
	}
	calendar.Modifier = calendar.Creator
	return dao.CreateSiteEnvCalendar(calendar)
}

// 获取场地日期段的环境日历
func EnumSiteEnvCalendars(bDate string, eDate string, siteId string) ([]entities.SiteEnvCalendar, error) {
	return dao.EnumSiteEnvCalendars(bDate, eDate, siteId)
}

// 更新场地环境日历
func UpdateSiteEnvCalendar(calendar *entities.SiteEnvCalendar) error {

	var result error
	tmpCalendar, err := dao.GetSiteEnvCalendar(calendar.SiteID, calendar.Date)
	if err == nil {
		if calendar.Modifier == "" {
			calendar.Modifier = "admin"
		}
		calendar.Creator = tmpCalendar.Creator
		calendar.CreateDatetime = tmpCalendar.CreateDatetime
		calendar.ModifyDatetime = time.Now()
		result = dao.UpdateSiteEnvCalendar(calendar)
	} else {
		result = err
	}
	return result
}
