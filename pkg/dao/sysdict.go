package dao

import "uwbwebapp/pkg/entities"

// 通过字典键获取字典值
func GetSystemDictValues(key string) ([]entities.SysDict, error) {
	// var sysdict entities.SysDict
	var results []entities.SysDict
	result := Database.Table("sys_dicts").Where("key=?", key).Order("order_key ASC").Find(&results)
	return results, result.Error

}

// 通过字典编码获取字典信息
func GetSystemDictValue(code string) (entities.SysDict, error) {
	var sysDict entities.SysDict
	result := Database.Table("sys_dicts").Where("code=?", code).First(&sysDict)
	return sysDict, result.Error
}

// 通过字典上级键获取其下级键字典信息
func GetChildrenSystemDictsByParentKey(parentKey string) ([]entities.SysDict, error) {
	var results []entities.SysDict
	result := Database.Table("sys_dicts").Where("parent_key=?", parentKey).Order("order_key ASC").Find(&results)
	return results, result.Error
}

// 通过字典上级编码获取其下级字典信息
func GetChildrenSystemDictsByParentCode(parentCode string) ([]entities.SysDict, error) {
	var results []entities.SysDict
	result := Database.Table("sys_dicts").Where("parent_code=?", parentCode).Order("order_key ASC").Find(&results)
	return results, result.Error
}
