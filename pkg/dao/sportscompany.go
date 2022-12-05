package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func CreateSportsCompany(compnay *entities.SportsCompany) error {
	result := Database.Create(compnay)
	return result.Error
}
func DeleteSportsCompanies(companyIds []string) error {
	var company entities.SportsCompany
	result := Database.Delete(&company, companyIds)
	return result.Error
}

func DeleteSportsCompany(compnay entities.SportsCompany) (bool, error) {
	result := Database.Delete(&compnay)
	var res bool = true
	if result.RowsAffected == 0 {
		res = false
	}
	return res, result.Error
}

func UpdateSportsCompany(company *entities.SportsCompany) error {
	res := Database.Save(company)
	return res.Error
}

func GetSportsCompanyById(id string) (entities.SportsCompany, error) {

	var company entities.SportsCompany
	result := Database.First(&company, "id=?", id)
	return company, result.Error
}

func GetSportsCompanyCount(companyQueryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var company entities.SportsCompany
	var result *gorm.DB
	if companyQueryCodition.LikeValue != "" {
		result = Database.Model(&company).Where(`name LIKE ? OR 
	address LIKE ? OR 
	telephone_list LIKE ? OR
	description LIKE ?`,
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%").Count(&count)
	} else {
		result = Database.Model(&company).Count(&count)
	}
	return count, result.Error
}

func QuerySportsCompanies(companyQueryCodition entities.QueryCondition) ([]entities.SportsCompany, error) {
	var company entities.SportsCompany
	var companies []entities.SportsCompany
	var result *gorm.DB
	if companyQueryCodition.LikeValue != "" {
		result = Database.Model(&company).Select("id, name, telephone_list, address, Modifier, Creator, modify_datetime, create_datetime").Where(`name LIKE ? OR 
	address LIKE ? OR 
	telephone_list LIKE ? OR
	description LIKE ?`,
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%",
			"%"+companyQueryCodition.LikeValue+"%").
			Order("modify_datetime DESC").
			Limit(int(companyQueryCodition.PageSize)).
			Offset(int(companyQueryCodition.PageSize * (companyQueryCodition.PageIndex - 1))).
			Find(&companies)
	} else {
		result = Database.Model(&company).
			Order("modify_datetime DESC").
			Limit(int(companyQueryCodition.PageSize)).
			Offset(int(companyQueryCodition.PageSize * (companyQueryCodition.PageIndex - 1))).
			Find(&companies)
	}
	return companies, result.Error
}
