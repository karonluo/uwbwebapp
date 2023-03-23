package biz

import (
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/ahmetb/go-linq/v3"
)

func CreateSwimmerSwimmReport(report *entities.SwimmerSwimmReport) error {
	report.Date = time.Now()
	report.CreateDatetime = report.Date
	report.ModifyDatetime = report.Date
	report.Creator = "admin"
	report.Modifier = report.Creator
	return dao.CreateSwimmerSwimmReport(report)

}
func EnumSwimmerSwimmReports(beginDate time.Time, endDate time.Time, swimmerId string) ([]entities.SwimmerSwimmReport, error) {
	beginDate = tools.SplitDateTimeToAllDay(beginDate).BeginDatetime
	endDate = tools.SplitDateTimeToAllDay(endDate).EndDatetime
	return dao.EnumSwimmerSwimmReports(beginDate, endDate, swimmerId)
}

func GetSwimmersSumDistenceOrder(siteId string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	// 从缓存中获取所有该游泳场在场的游泳者信息
	// first from db get swimmers on site
	swimmers, err := dao.EnumAllSwimmerCalendarsOnSiteBySiteId(siteId)
	for _, sw := range swimmers {

		tmp := make(map[string]interface{})
		info, err := cache.GetUWBTerminalTagFromRedis(sw.SwimmerID)
		if err == nil {
			tmp["swimmerId"] = info.Properties.SwimmerId
			tmp["swimmerDisplayName"] = info.Properties.SwimmerDisplayName
			tmp["swimmerGender"] = info.Properties.SwimmerGender
			tmp["sumDistence"] = info.SumDistance
			result = append(result, tmp)
		}
	}
	linq.From(result).OrderByDescendingT(func(s map[string]interface{}) float32 {
		return s["sumDistence"].(float32)
	}).Take(5).ToSlice(&result)

	return result, err
}
