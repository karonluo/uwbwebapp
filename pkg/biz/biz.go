package biz

import (
	"context"
	"strings"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/dao"
)

func EnumSiteOwners(site_id string) []map[string]interface{} {
	return dao.EnumSiteOwners(site_id)
}

func ClearSiteOwners(site_id string) (bool, string) {
	return dao.ClearSiteOwners(site_id)
}

// 设置场地负责人
func SetSiteOwners(siteId string, sysUserId string, jobTitle string) error {
	ids := strings.Split(sysUserId, ",")
	var errResult error = nil
	for _, id := range ids {
		_, errResult = dao.SetSiteOwner(siteId, id, jobTitle)
		if errResult != nil {
			break
		}
	}
	return errResult
}

func GetLoginInformation(authorization string, fieldName string) (string, error) {
	rctx := context.Background()
	result := dao.RedisDatabase.HGet(rctx, "token_"+authorization, fieldName)
	err := result.Err()
	msg := result.Val()
	return msg, err
}

func CheckLogin(authorization string) (bool, error) {
	result := true

	// TODO: 用于测试代码，允许 authorization 为 test 的时候通过。
	if authorization == "test" {
		return true, nil
	}
	rctx := context.Background()
	res := dao.RedisDatabase.HExists(rctx, "token_"+authorization, "token")
	err := res.Err()
	if err != nil {
		result = false

	} else {
		result = res.Val()
		dao.RedisDatabase.Expire(rctx, "token_"+authorization, time.Duration(conf.WebConfiguration.SessionExpireMinute)*time.Minute) // 重置超时时间
	}
	return result, err
}
