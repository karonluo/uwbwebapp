package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

// 枚举指定场地所属所有的公司。
func EnumSportsCompanySitesBySiteId(siteId string) ([]entities.CompanySite, error) {
	return dao.EnumSportsCompanySitesBySiteId(siteId)
}
