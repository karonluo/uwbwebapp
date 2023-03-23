package dao

import "uwbwebapp/pkg/entities"

// TODO: CRUD Site Fence

func CreateSiteFence(fence *entities.SiteFence) error {
	return Database.Create(fence).Error
}

func UpdateSiteFence(fence *entities.SiteFence) error {
	result := Database.Table("site_fence").Where("site_id=? AND code=?", fence.SiteID, fence.Code).UpdateColumns(&fence)
	return result.Error
}

func EnumSiteFences(siteId string) ([]entities.SiteFence, error) {
	var fences []entities.SiteFence
	err := Database.Table("site_fence").Select("site_id, code, modifier, creator, modify_datetime, create_datetime").Where("site_id=?", siteId).Order("modify_datetime desc").Find(&fences).Error
	return fences, err
}
func EnumSiteFenceCodes(siteId string) ([]string, error) {
	var codes []string
	err := Database.Table("site_fence").Select("code").Where("site_id=?", siteId).Order("modify_datetime desc").Find(&codes).Error
	return codes, err
}

func GetSiteFence(siteId string, code string) (entities.SiteFence, error) {
	var fence entities.SiteFence
	err := Database.Table("site_fence").Where("site_id=? AND code=?", siteId, code).First(&fence).Error
	return fence, err
}
