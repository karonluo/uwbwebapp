package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func CreateSiteUser(siteUser entities.SiteUser) error {
	return dao.CreateSiteUser(siteUser)
}
