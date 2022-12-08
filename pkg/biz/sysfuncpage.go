package biz

import (
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func EnumSysFuncPages() ([]entities.SysFuncPage, int64) {

	sysFuncPage, recordCount := dao.EnumSysFuncPages()
	return sysFuncPage, recordCount
}

func CreateSysFuncPage(page *entities.SysFuncPage) (string, error) {
	page.ID = uuid.New().String()
	page.CreateDatetime = time.Now()
	page.ModifyDatetime = page.CreateDatetime
	if page.Creator == "" {
		page.Creator = "admin"
	}
	page.Modifier = page.Creator
	err := dao.CreateSysFuncPage(page)
	return page.ID, err
}
