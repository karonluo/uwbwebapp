package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

// 枚举指定游泳者所属所有的公司。
func EnumSportsCompanySwimmersBySwimmerId(swimmerId string) ([]entities.CompanySwimmer, error) {
	return dao.EnumSportsCompanySwimmersBySwimmerId(swimmerId)
}
