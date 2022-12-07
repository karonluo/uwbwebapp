package dao

import "uwbwebapp/pkg/entities"

func CreateCompanySwimmer(companySwimmer *entities.CompanySwimmer) error {
	return Database.Create(companySwimmer).Error
}

// 通过公司编号及游泳者编号查找，主要用于证明是否该用户是否已经加入到制定体育运动公司了。
func GetCompanySwimmerByCompanyIDAndSwimmerID(sportsCompanyId string, swimmerId string) (*entities.CompanySwimmer, int64, error) {
	var companySwimmer *entities.CompanySwimmer
	result := Database.Table("company_swimmers").Select("sports_company_id, swimmer_id").Where("sports_company_id=? and swimmer_id=?", sportsCompanyId, swimmerId).First(&companySwimmer)
	return companySwimmer, result.RowsAffected, result.Error
}

// 枚举指定游泳者所属所有的公司。
func EnumSportsCompanySwimmersBySwimmerId(swimmerId string) ([]entities.CompanySwimmer, error) {
	var compines []entities.CompanySwimmer
	result := Database.Table("company_swimmers").Where("swimmer_id=?", swimmerId).Find(&compines)
	return compines, result.Error
}
