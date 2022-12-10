package dao

import "uwbwebapp/pkg/entities"

func CreateCompanySite(companySite *entities.CompanySite) error {
	return Database.Create(companySite).Error
}

// 枚举指定游泳者所属所有的公司。
func EnumSportsCompanySitesBySiteId(siteId string) ([]entities.CompanySite, error) {
	var compines []entities.CompanySite
	result := Database.Table("company_sites").Where("site_id=?", siteId).Find(&compines)
	return compines, result.Error
}
