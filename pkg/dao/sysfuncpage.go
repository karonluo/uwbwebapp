package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"
)

func CreateSysFuncPage(sysFuncPage *entities.SysFuncPage) error {
	result := Database.Create(&sysFuncPage)
	return result.Error
}

func EnumSysFuncPagesFromDB() ([]entities.SysFuncPage, int64) {

	var sysFuncPages []entities.SysFuncPage
	result := Database.Find(&sysFuncPages)
	return sysFuncPages, result.RowsAffected
}

func EnumSysFuncPagesByRoleId() ([]entities.SysFuncPage, error) {
	var sysFuncPages []entities.SysFuncPage
	result := Database.Find(&sysFuncPages)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return sysFuncPages, result.Error
	}
}

func GetSysFuncPageFromRedis(url string, method string) entities.SysFuncPage {
	var funcPage entities.SysFuncPage
	rctx := context.Background()
	strCmd := cache.RedisDatabase.Get(rctx, fmt.Sprintf("funcpage_%s|%s", url, method))
	json.Unmarshal([]byte(strCmd.Val()), &funcPage)
	return funcPage
}

func DeleteSysFuncPage(id string) error {
	page, err := GetSysFuncPage(id)
	if err == nil {
		err = Database.Delete(&page).Error
		if err == nil {
			rctx := context.Background()
			cache.RedisDatabase.Del(rctx, fmt.Sprintf("funcpage_%s|%s", page.URLAddress, page.URLMethod))
			defer rctx.Done()
		}
	}
	return err
}

func GetSysFuncPage(id string) (entities.SysFuncPage, error) {
	var page entities.SysFuncPage
	err := Database.Model(page).Where("id=?", id).First(&page).Error
	return page, err
}
