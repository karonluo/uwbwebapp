package dao

import "uwbwebapp/pkg/entities"

func CreateSiteUser(siteUser *entities.SiteUser) error {
	return Database.Create(siteUser).Error
}

func DeleteSiteUsersBySiteId(siteId string) error {
	return Database.Exec("DELETE FROM site_users WHERE site_id = ?", siteId).Error
}

func EnumSiteUsersBySiteId(siteId string) ([]entities.SiteUser, error) {
	var siteusers []entities.SiteUser
	err := Database.Table("site_users").Where("site_id=?", siteId).Order("modify_datetime DESC").Find(&siteusers).Error
	return siteusers, err
}

func EnumSiteUsersByUserId(userId string) ([]entities.SiteUser, error) {
	var siteusers []entities.SiteUser
	err := Database.Table("site_users").Where("sys_user_id=?", userId).Order("modify_datetime DESC").Find(&siteusers).Error
	return siteusers, err
}

func UpdateSitUserSiteDisplayName(siteId string, displayName string) error {
	return Database.Table("site_users").Where("site_id=?", siteId).UpdateColumn("site_display_name", displayName).Error
}

func UpdateSiteUserUserDisplayName(userId string, displayName string) error {
	return Database.Table("site_users").Where("sys_user_id=?", userId).UpdateColumn("sys_user_displayname", displayName).Error
}
