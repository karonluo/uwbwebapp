package biz

import (
	"strings"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"
)

func ClearSiteUsers(siteID string) error {
	err := dao.DeleteSiteUsersBySiteId(siteID)
	if err == nil {
		err = dao.SetSiteUserNames(siteID, ``)
	}
	return err
}
func EnumSiteUsersBySiteId(siteId string) ([]entities.SiteUser, error) {
	return dao.EnumSiteUsersBySiteId(siteId)
}
func EnumSiteUsersByUserId(userId string) ([]entities.SiteUser, error) {
	return dao.EnumSiteUsersByUserId(userId)
}
func SetSiteUsers(siteId string, siteUsers []entities.SiteUser) error {
	var err error
	var userNames string
	// 清空 该场地的用户信息。
	err = ClearSiteUsers(siteId)
	if len(siteUsers) > 0 {

		// 业务上不允许在场地信息管理中同时更新不同的场地的场地用户集合，因此只允许第一个出现的场地。

		if err == nil {

			for _, siteUser := range siteUsers {
				// 业务上不允许在场地信息管理中同时更新不同的场地的场地用户集合，因此只允许第一个出现的场地。
				siteUser.SiteID = siteId
				siteUser.CreateDatetime = time.Now()
				siteUser.ModifyDatetime = time.Now()
				if siteUser.Modifier == "" {
					siteUser.Modifier = "admin"
				}
				if siteUser.Creator == "" {
					siteUser.Creator = "admin"
				}
				err = dao.CreateSiteUser(&siteUser)
				if err != nil {
					// 一旦出现错误，立刻停止剩下的设置。
					break
				}
				userNames = userNames + siteUser.SysUserDisplayname + "|" + siteUser.JobTitle + ","
			}
		}
		if err == nil {
			userNames = strings.TrimRight(userNames, ",") //去掉最后一个逗号分隔符
			var site entities.Site

			// 注意此处需要根据字段的长度进行字符串截取，否则可能会出现错误。
			var fieldSize int
			fieldSize, err = tools.GetDatabaseTableFieldSize(site, "Users")
			if err == nil {
				if len(userNames) > fieldSize {
					nameRune := []rune(userNames)
					userNames = string(nameRune[0 : fieldSize/2])
					userNames = strings.TrimRight(userNames, ",")
				}
				err = dao.SetSiteUserNames(siteId, userNames) // 设置场地用户信息显示名，方便单表查询。

			}
			//err = dao.SetSiteUserNames(siteId, userNames)
		}
	}
	return err
}
