package biz

import (
	"fmt"
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func InitSysFuncPageToRedis() {
	funcPages, _ := dao.EnumSysFuncPagesFromDB()
	fmt.Print("将数据库中的所有系统功能页面信息放入内存数据库")
	for _, funcPage := range funcPages {
		cache.SetFuncPageToRedis(&funcPage, 0)
	}
	fmt.Println("......完成")
}

func EnumSysFuncPages() ([]entities.SysFuncPage, int64) {

	sysFuncPage, recordCount := dao.EnumSysFuncPagesFromDB()
	return sysFuncPage, recordCount
}
func ClearSysFuncPages() error {

	return dao.Database.Exec("DELETE FROM sys_func_pages").Error
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
	if err == nil {
		cache.SetFuncPageToRedis(page, 0)
	}
	return page.ID, err
}
func DeleteSysFuncPage(id string) error {
	return dao.DeleteSysFuncPage(id)
}

func GetSysFuncPageFromRedis(url string, method string) entities.SysFuncPage {
	return dao.GetSysFuncPageFromRedis(url, method)
}
