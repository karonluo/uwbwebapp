package dao

import (
	"uwbwebapp/pkg/entities"
)

func CreateSysFuncPage(sysFuncPage *entities.SysFuncPage) error {
	result := Database.Create(&sysFuncPage)
	return result.Error
}

func EnumSysFuncPages() ([]entities.SysFuncPage, int64) {

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
