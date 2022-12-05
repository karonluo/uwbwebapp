package dao

import (
	"uwbwebapp/pkg/entities"

	"gorm.io/gorm"
)

func GetSiteById(id string) (entities.Site, error) {
	var site entities.Site
	result := Database.First(&site, "id=?", id)
	// InitDatabase()
	return site, result.Error
}

// func GetSiteByCode(code string) entities.Site {
// 	var site entities.Site
// 	Database.First(&site, "code=?", code)
// 	return site
// }

func CreateSite(site *entities.Site) error {

	return Database.Create(site).Error
}

func DeleteSites(siteIds []string) error {
	var site entities.Site
	result := Database.Delete(&site, siteIds)
	return result.Error
}

func ClearSiteOwners(site_id string) (bool, string) {
	var result bool = true
	var msg string = "OK"
	// tx := Database.Exec("DELETE FROM site_owners WHERE site_id=?", site_id)
	// if tx.Error != nil {
	// 	fmt.Println("Clear is error")
	// 	Database.Rollback()
	// 	result = false
	// } else {
	// 	fmt.Println("Clear is ok")
	// 	Database.Commit()
	// }
	db, _ := Database.DB()
	_, err := db.Exec("DELETE FROM site_users WHERE site_id=$1", site_id)

	if err != nil {
		result = false
		msg = err.Error()
	}
	return result, msg
}

func SetSiteOwner(siteId string, sysUserId string, jobTitle string) (bool, error) {
	var result bool = true
	var sqlcmd string = "INSERT INTO site_users (site_id, sys_user_id, job_title) VALUES ($1, $2, $3)"
	db, _ := Database.DB()
	_, execerr := db.Exec(sqlcmd, siteId, sysUserId, jobTitle)
	if execerr != nil {
		result = false
	}
	return result, execerr
}

func GetSitesCount(queryCodition entities.QueryCondition) (int64, error) {
	var count int64
	var site entities.Site
	var result *gorm.DB
	if queryCodition.LikeValue != "" {
		result = Database.Model(&site).Where(`Address LIKE ? OR Contact LIKE ? OR Users LIKE ? OR Display_Name LIKE ?`,
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%",
			"%"+queryCodition.LikeValue+"%").Count(&count)
	} else {
		result = Database.Model(&site).Count(&count)
	}
	return count, result.Error
}

func QuerySites(queryCodition entities.QueryCondition) ([]entities.Site, error) {
	var site entities.Site
	var sites []entities.Site
	var result *gorm.DB
	var selectFileds = `Id,Address,Contact,Users,Display_Name,create_datetime,modify_datetime,Creator,Modifier`
	if queryCodition.LikeValue != "" {
		result = Database.Model(&site).Select(selectFileds).
			Where(`Address LIKE ? OR Contact LIKE ? OR Users LIKE ? OR Display_Name LIKE ?`,
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%",
				"%"+queryCodition.LikeValue+"%").
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&sites)
	} else {
		result = Database.Model(&site).Select(selectFileds).
			Order("modify_datetime DESC").
			Limit(int(queryCodition.PageSize)).
			Offset(int(queryCodition.PageSize * (queryCodition.PageIndex - 1))).
			Find(&sites)
	}
	return sites, result.Error
}

func UpdateSite(site *entities.Site) error {
	result := Database.Table("sites").Where("id=?", site.Id).UpdateColumns(site)
	return result.Error
}
