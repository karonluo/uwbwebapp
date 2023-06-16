package dao

import (
	"uwbwebapp/pkg/entities"
)

func CreateSwimmerSwimmReport(report *entities.SwimmerSwimmReport) error {
	err := Database.Create(report).Error
	return err
}

// func EnumSwimmerSwimmReports(beginDate time.Time, endDate time.Time, swimmerId string) ([]entities.SwimmerSwimmReport, error) {
// 	var reports []entities.SwimmerSwimmReport
// 	err := Database.Table("swimmer_swimm_reports").Where("date BETWEEN ? AND ? ", beginDate, endDate).Find(&reports).Error
// 	return reports, err
// }
