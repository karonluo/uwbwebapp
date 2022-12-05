package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/entities"

	"github.com/go-redis/redis/v8"
)

var RedisDatabase *redis.Client

func InitRedisDatabase() bool {
	var res bool
	rdbconf := conf.WebConfiguration.RedisConf
	fmt.Print("内存数据库初始化")
	RedisDatabase = redis.NewClient(&redis.Options{
		Addr:     rdbconf.Host + ":" + rdbconf.Port,
		Password: rdbconf.Password,
		DB:       rdbconf.DBId, // use default DB
	})
	ctx := context.Background()
	result := RedisDatabase.Ping(ctx)
	if result.Val() == "PONG" {
		res = true
		fmt.Println("......成功")
	} else {
		res = false
		fmt.Println("......失败")
	}
	return res

}

func EnumSysUserFromRedis() {

}

func InitSysUserToRedis() {
	// sysusers, _ := EnumSysUserFromDB()
	fmt.Print("将数据库中的用户信息放入内存数据库")
	sysusers, _ := EnumSysUserFromDB()

	for _, sysuser := range sysusers {

		// fmt.Println(sysuser)
		// data := tools.ReflectMethod(sysuser)
		SetUserToRedis(sysuser)

	}
	fmt.Println("......完成")
}
func SetUserToRedis(user entities.SysUser) {
	ctx := context.Background()
	result, _ := json.Marshal(&user)
	RedisDatabase.Set(ctx, "sysuser_"+user.LoginName, string(result), time.Hour*24)
}
func TestRedis() {
	fmt.Println("测试内存数据库")
	rctx := context.Background()
	result, _ := RedisDatabase.Get(rctx, "key_test").Result()
	fmt.Println(result)
	fmt.Println("测试内存数据库完成")
}
