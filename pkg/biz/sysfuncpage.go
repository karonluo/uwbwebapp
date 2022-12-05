package biz

import (
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func EnumSysFuncPages() ([]entities.SysFuncPage, int64) {

	sysFuncPage, recordCount := dao.EnumSysFuncPages()
	return sysFuncPage, recordCount
}
