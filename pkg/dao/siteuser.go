package dao

import "uwbwebapp/pkg/entities"

func CreateSiteUser(siteUser entities.SiteUser) error {
	return Database.Create(&siteUser).Error
}
