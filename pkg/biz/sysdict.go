package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

// 通过字典键获取字典值
func GetSystemDictValues(key string) ([]entities.SysDict, error) {
	return dao.GetSystemDictValues(key)
}

// 通过字典编码获取字典值
func GetSystemDictValue(code string) (entities.SysDict, error) {
	return dao.GetSystemDictValue(code)
}

// 通过字典上级键获取其下级键字典信息
func GetChildrenSystemDictsByParentKey(parentKey string) ([]entities.SysDict, error) {
	return dao.GetChildrenSystemDictsByParentKey(parentKey)
}

// 通过字典上级编码获取其下级字典信息
func GetChildrenSystemDictsByParentCode(parentCode string) ([]entities.SysDict, error) {
	return dao.GetChildrenSystemDictsByParentCode(parentCode)
}
