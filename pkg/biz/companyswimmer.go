package biz

import (
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

// 枚举指定游泳者所属所有的公司。
func EnumSportsCompanySwimmersBySwimmerId(swimmerId string) ([]entities.CompanySwimmer, error) {
	return dao.EnumSportsCompanySwimmersBySwimmerId(swimmerId)
}

// 通过游泳者唯一编号和公司唯一编号获取公司会员信息
func GetCompanySwimmerByCompanyIDAndSwimmerID(companyId string, swimmerId string) (entities.CompanySwimmer, int64, error) {
	return dao.GetCompanySwimmerByCompanyIDAndSwimmerID(companyId, swimmerId)
}

func UpdateCompanySwimmer(companySwimmer *entities.CompanySwimmer) error {
	tmp, _, err := dao.GetCompanySwimmerByCompanyIDAndSwimmerID(companySwimmer.SportsCompanyID, companySwimmer.SwimmerID)
	if err == nil {
		if companySwimmer.Modifier == "" {
			companySwimmer.Modifier = "admin"
		}
		companySwimmer.ModifyDatetime = time.Now()
		// 下面两句代码是为了防止被接口调用者改变。
		companySwimmer.CreateDatetime = tmp.CreateDatetime
		companySwimmer.Creator = tmp.Creator
	}
	return dao.UpdateCompanySwimmer(companySwimmer)
}
