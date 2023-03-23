package biz

import (
	"fmt"
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

// 通过字典键获取字典
func GetSystemDicts(key string) ([]entities.SysDict, error) {
	return dao.GetSystemDicts(key)
}

// 通过字典编码获取字典
func GetSystemDict(code string) (entities.SysDict, error) {
	dict, err := cache.GetDictFromRedis(code)
	if err != nil {
		dict, err = dao.GetSystemDict(code)
		cache.SetDictToRedis(&dict)
	}
	return dict, err
}
func SetSystemDict(dict *entities.SysDict) error {
	oriDict, err := GetSystemDict(dict.Code)
	if err == nil {
		dict.CreateDatetime = oriDict.CreateDatetime
		dict.Creator = oriDict.Creator
		if dict.Modifier == "" {
			dict.Modifier = "admin"
		}
		dict.ModifyDatetime = time.Now()
		err = dao.SetSystemDict(dict)
		if err == nil {
			cache.SetDictToRedis(dict)
		}
	}
	return err
}

// 通过字典上级键获取其下级键字典信息
func GetChildrenSystemDictsByParentKey(parentKey string) ([]entities.SysDict, error) {
	return dao.GetChildrenSystemDictsByParentKey(parentKey)
}

// 通过字典上级编码获取其下级字典信息
func GetChildrenSystemDictsByParentCode(parentCode string) ([]entities.SysDict, error) {
	return dao.GetChildrenSystemDictsByParentCode(parentCode)
}

// 将所有字典设置到Cache.
func InitDictToRedis() error {

	fmt.Print("将数据库中的所有字典信息放入内存数据库")
	dicts, err := dao.EnumAllDicts()
	if err == nil {
		cache.SetAllDictsToRedis(dicts)
	}
	fmt.Println("......完成")

	return err
}

// 更新字典缓存
func SetDictToRedis(dict *entities.SysDict) {
	cache.SetDictToRedis(dict)
}
