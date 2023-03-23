package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/ahmetb/go-linq/v3"
)

func SiteSwimmerReport(siteId string) (entities.SiteSwimmerReport, error) {
	var report entities.SiteSwimmerReport
	sws, err := dao.EnumSiteSimmerForReport(siteId)
	if err == nil {
		report.Age0014 = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Age <= 14)
		}).Count()

		report.Age1529 = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Age > 14 && sw.Age <= 29)
		}).Count()

		report.Age3049 = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Age > 29 && sw.Age <= 49)
		}).Count()

		report.Age5059 = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Age > 49 && sw.Age <= 59)
		}).Count()

		report.Age60up = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Age > 59)
		}).Count()

		report.GenderFemale = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Gender == "女")
		}).Count()

		report.GenderMale = linq.From(sws).WhereT(func(sw entities.Swimmer) bool {
			return (sw.Gender == "男")
		}).Count()
	}
	return report, err
}
