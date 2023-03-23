package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func main() {
	conf.LoadWebConfig("../conf/WebConfig.json")
	dao.InitDatabase()
	if cache.InitRedisDatabase() {
		fmt.Println("日志记录器客户端初始化完成。")
		rctx := context.Background()
		for {

			strCmd := cache.RedisDatabase.BLPop(rctx, time.Hour*24, "SystemLog")
			if strCmd.Err() != nil {
				fmt.Println(strCmd.Err().Error())
			} else {
				// TODO： 将日志写入数据库，建议使用 MangoDB.
				var log entities.SystemLog
				strLog := strCmd.Val()[1]
				err := json.Unmarshal([]byte(strLog), &log)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println(log.String())
					// 暂时使用数据库方式记录。
					dao.WriteSystemLogToDB(&log)
				}
			}

		}
	}

}
