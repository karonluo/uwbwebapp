package dao

import "uwbwebapp/pkg/entities"

func CreateSiteUser(siteUser *entities.SiteUser) error {
	return Database.Create(siteUser).Error
}

func DeleteSiteUsersBySiteId(siteId string) error {
	return Database.Exec("DELETE FROM site_users WHERE site_id = ?", siteId).Error
}
