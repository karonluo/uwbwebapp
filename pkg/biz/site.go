package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func GetSiteById(id string) (entities.Site, error) {
	return dao.GetSiteById(id)
}

func CreateSite(site *entities.Site) error {
	site.Id = uuid.New().String()
	site.CreateDatetime = time.Now()
	if site.Creator == "" {
		site.Creator = "admin"
	}
	site.ModifyDatetime = site.CreateDatetime
	site.Modifier = site.Creator
	err := dao.CreateSite(site)
	return err
}

func QuerySites(queryCondition entities.QueryCondition) ([]entities.Site, int64, int64, error) {
	var sites []entities.Site
	dataRecordCount, err := dao.GetSitesCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		sites, err = dao.QuerySites(queryCondition)
	}

	return sites, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func DeleteSites(ids []string) error {
	return dao.DeleteSites(ids)
}

func UpdateSite(site *entities.Site) error {
	var tmpSite entities.Site
	var err error
	tmpSite, err = dao.GetSiteById(site.Id)
	if err == nil {
		// 防止以下字段被修改
		site.CreateDatetime = tmpSite.CreateDatetime
		site.Creator = tmpSite.Creator
		site.ModifyDatetime = time.Now()
		if site.Modifier == "" {
			site.Modifier = "admin"
		}
		err = dao.UpdateSite(site)
	}
	return err
}
