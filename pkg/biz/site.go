package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

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

// 返回
// 本次新加入的公司
// 错误信息集合
func SiteJoinInSportsCompanies(css []entities.CompanySite, siteId string) ([]entities.CompanySite, []error) {
	var err error
	var errs []error
	// var originJoined []entities.CompanySwimmer
	var newJoined []entities.CompanySite

	err = dao.ClearAllCompaniesFromSite(siteId)
	if err == nil {
		for _, cs := range css {
			// TODO: 需要优化

			/*
				_, dataRecordCount, err = dao.GetCompanySwimmerByCompanyIDAndSwimmerID(cs.SportsCompanyID, cs.SwimmerID)
				if dataRecordCount == 1 {
					originJoined = append(originJoined, cs)
				} else {
					err = dao.CreateCompanySwimmer(&cs)
					if err == nil {
						newJoined = append(newJoined, cs)
					} else {
						errs = append(errs, err)
					}
				}
			*/
			// 目前的方式，需要前端进行优先判断，只加入曾经未加入的公司。
			cs.CreateDatetime = time.Now()
			cs.ModifyDatetime = cs.CreateDatetime
			if cs.Creator == "" {
				cs.Creator = "admin"
			}
			cs.Modifier = cs.Creator
			err = dao.CreateCompanySite(&cs)
			if err == nil {
				newJoined = append(newJoined, cs)
			} else {
				errs = append(errs, err)
				tools.ProcessError("biz.SiteJoinInSportsCompanies", `err = dao.CreateCompanySite(&cs)`, err)

			}
		}
	} else {
		tools.ProcessError("biz.SiteJoinInSportsCompanies", `err = dao.ClearAllCompaniesFromSite(siteId)`, err)
	}
	return newJoined, errs
}
