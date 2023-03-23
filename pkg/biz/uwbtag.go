package biz

import (
	"errors"
	"math"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func CreateUWBTag(tag *entities.UWBTag) error {
	var err error

	if tag.Code == "" {
		errStr := "新增 UWB 标签必须填写 Code 编码。"
		err = errors.New(errStr)
	} else {

		if tag.Creator == "" {
			tag.Creator = "admin"
		}
		if tag.Modifier == "" {
			tag.Modifier = "admin"
		}
		tag.IsBound = false
		tag.CreateDatetime = time.Now()
		tag.ModifyDatetime = time.Now()
		err = dao.CreateUWBTag(tag)
		if err == nil {
			// 调用 UWB 设备管理平台接口 创建新的标签终端
			_, err = UWBDevicePlatformCreateTerminal("noname", "nogender", "noid", tag.Code, conf.WebConfiguration.UWBDevicePlatformConf.DefaultTerminalModelId)

		}
	}
	return err
}

func GetUWBTag(code string) (entities.UWBTag, error) {
	return dao.GetUWBTag(code)
}

func QueryUWBTags(queryCondition entities.QueryCondition, companyIds []string, isBound interface{}) ([]entities.UWBTag, int64, int64, error) {
	var tags []entities.UWBTag
	dataRecordCount, err := dao.GetUWBTagCount(queryCondition, companyIds, isBound)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		tags, err = dao.QueryUWBTags(queryCondition, companyIds, isBound)
	}
	return tags, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func DeleteUWBTags(codes []string) error {
	return dao.DeleteUWBTags(codes)
}
