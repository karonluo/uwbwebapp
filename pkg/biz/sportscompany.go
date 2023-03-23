package biz

import (
	"fmt"
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func CreateSportsCompany(company *entities.SportsCompany) (string, error) {
	if company.Id == "" {
		company.Id = uuid.New().String()
	}

	if company.Creator == "" {
		company.Creator = "admin"
	}
	if company.Modifier == "" {
		company.Modifier = "admin"

	}
	company.CreateDatetime = time.Now()
	company.ModifyDatetime = time.Now()
	err := dao.CreateSportsCompany(company)
	return company.Id, err
}
func DeleteSportsCompanies(companyIds []string) error {
	return dao.DeleteSportsCompanies(companyIds)
}
func DeleteSportsCompany(company entities.SportsCompany) (bool, error) {
	return dao.DeleteSportsCompany(company)
}

func UpdateSportsCompany(company *entities.SportsCompany) error {

	tmpCompany, err := dao.GetSportsCompanyById(company.Id)
	if err == nil {
		tmpCompany.Address = company.Address
		tmpCompany.Description = company.Description
		tmpCompany.Modifier = company.Modifier
		tmpCompany.ModifyDatetime = time.Now()
		if company.Name != tmpCompany.Name {
			// 当修改了公司名称，需要更新用户表中所有相关用户的公司名称。
			dao.UpdateSysUserCompanyName(company.Id, company.Name)
			dao.UpdateUWBTagSportsCompanyName(company.Id, company.Name)
			dao.UpdateCompanySwimmerCompanyName(company.Id, company.Name)
		}
		tmpCompany.Name = company.Name
		tmpCompany.TelephoneList = company.TelephoneList
		err = dao.UpdateSportsCompany(&tmpCompany)
	}
	return err
}

func GetSportsCompanyById(id string) (entities.SportsCompany, error) {
	return dao.GetSportsCompanyById(id)
}

func RelSportsCompanyAndSite(company_id string, site_id string) {
	fmt.Println(company_id, site_id)

}

func QueryCompanies(queryCondition entities.QueryCondition) ([]entities.SportsCompany, int64, int64, error) {
	var companies []entities.SportsCompany
	dataRecordCount, err := dao.GetSportsCompanyCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		companies, err = dao.QuerySportsCompanies(queryCondition)
	}
	return companies, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func EnumSportsCompaniesByGroupId(groupId string) ([]entities.SportsCompany, error) {
	return dao.EnumSportsCompaniesByGroupId(groupId)
}

func EnumSportsCompaniesByRightUser(userId string) ([]entities.SportsCompany, error) {
	return dao.EnumSportsCompaniesByRightUser(userId)
}
