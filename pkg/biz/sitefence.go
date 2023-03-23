package biz

import (
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func CreateSiteFence(fence *entities.SiteFence) error {
	if fence.Creator == "" {
		fence.Creator = "admin"
	}
	fence.CreateDatetime = time.Now()
	fence.Modifier = fence.Creator
	if fence.Code == "" {
		fence.Code = uuid.NewString()
	}
	fence.ModifyDatetime = fence.CreateDatetime
	return dao.CreateSiteFence(fence)
}

func UpdateSiteFence(fence *entities.SiteFence) error {
	oriFence, err := dao.GetSiteFence(fence.SiteID, fence.Code)
	if err == nil {
		fence.ModifyDatetime = time.Now()
		fence.Creator = oriFence.Creator
		fence.CreateDatetime = oriFence.CreateDatetime
		if fence.Modifier == "" {
			fence.Modifier = "admin"
		}
	}
	return dao.UpdateSiteFence(fence)
}

func EnumSiteFences(siteId string) ([]entities.SiteFence, error) {
	return dao.EnumSiteFences(siteId)
}

func EnumSiteFenceCodes(siteId string) ([]string, error) {
	return dao.EnumSiteFenceCodes(siteId)
}

func GetSiteFence(siteId string, code string) (entities.SiteFence, error) {
	return dao.GetSiteFence(siteId, code)
}
