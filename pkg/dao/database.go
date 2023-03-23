package dao

import (
	"fmt"
	"time"
	"uwbwebapp/conf"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database *gorm.DB

func InitDatabase() bool {
	var result bool
	if Database == nil {
		fmt.Print("数据库初始化")
		dbconf := conf.WebConfiguration.DBConf
		dsn := "host=" + dbconf.Host + " user=" + dbconf.User + " password=" + dbconf.Password + " dbname=" + dbconf.DBName + " port=" + dbconf.Port + " sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error)})

		if err == nil {

			sql, _ := db.DB()
			sql.SetConnMaxIdleTime(time.Hour)
			sql.SetConnMaxLifetime(24 * time.Hour)
			sql.SetMaxIdleConns(100)
			sql.SetMaxOpenConns(200)
			Database = db
			fmt.Println("......成功")
			result = true
		} else {
			fmt.Println("......失败")
			result = false

		}

	}
	return result
}

func TestDB() {
	fmt.Println("测试数据库")
	InitDatabase()
	//site := GetSiteById("05d6739d-6ed7-4ff9-bbad-3648b78b19bc")
	sysusers, recordCount := EnumSysUserFromDB()
	fmt.Printf("一共查询到: %d 条数据。\r\n", recordCount)
	for _, v := range sysusers {
		fmt.Println(v)
	}
	//fmt.Println(site)
	fmt.Println("测试数据库完成")
}
