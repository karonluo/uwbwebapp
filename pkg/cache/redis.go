package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
		PoolSize: 100,          // 连接池
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

// 将用户信息放置到 Redis 缓存中。
// 当设置为0分钟时，按照24小时设置超时时间
func SetUserToRedis(user entities.SysUser, timeOutDurationMinute int) {
	rctx := context.Background()
	result, _ := json.Marshal(&user)
	if timeOutDurationMinute == 0 {
		RedisDatabase.Set(rctx, "sysuser_"+user.LoginName, string(result), time.Hour*24)
	} else {

		tmpTimeout := time.Duration(timeOutDurationMinute) * time.Minute
		RedisDatabase.Set(rctx, "sysuser_"+user.LoginName, string(result), tmpTimeout)
	}

}

func EnumAllUWBBaseStationsBySiteIdFromRedis(siteId string) []entities.UWBBaseStation {
	var result []entities.UWBBaseStation
	rctx := context.Background()
	keyList := RedisDatabase.Keys(rctx, fmt.Sprintf("UWBBaseStation_%s_*", siteId)).Val()
	for _, key := range keyList {
		var station entities.UWBBaseStation
		if err := json.Unmarshal([]byte(RedisDatabase.Get(rctx, key).Val()), &station); err == nil {
			result = append(result, station)

		}
	}
	return result
}

func SetAllUWBBaseStationsToRedis(stations []entities.UWBBaseStation) {
	for _, station := range stations {
		SetUWBBaseStationToRedis(station)
	}
}

func SetUWBBaseStationToRedis(station entities.UWBBaseStation) {
	rctx := context.Background()
	strStation, _ := json.Marshal(station)
	RedisDatabase.Set(rctx, fmt.Sprintf("UWBBaseStation_%s_%s", station.SiteID, station.Code), string(strStation), 24*time.Hour)
	defer rctx.Done()
}

// 将字典表放入到 Redis 缓存中
func SetAllDictsToRedis(dicts []entities.SysDict) {

	rctx := context.Background()
	for _, dict := range dicts {
		bytesDict, _ := json.Marshal(dict)
		RedisDatabase.Set(rctx, "UWB_dict_"+dict.Code, string(bytesDict), 24*time.Hour)
	}

}

// 将字典放入到 Redis 缓存中
func SetDictToRedis(dict *entities.SysDict) {
	rctx := context.Background()
	bytesDict, _ := json.Marshal(dict)
	RedisDatabase.Set(rctx, "UWB_dict_"+dict.Code, string(bytesDict), 24*time.Hour)
}
func GetDictFromRedis(dictCode string) (entities.SysDict, error) {
	rctx := context.Background()
	var dict entities.SysDict
	result, er := RedisDatabase.Get(rctx, "UWB_dict_"+dictCode).Result()
	if er == nil {
		er = json.Unmarshal([]byte(result), &dict)
	}
	return dict, er
}

// 将系统功能页面信息放置到 Redis 缓存中。
func SetFuncPageToRedis(funcPage *entities.SysFuncPage, timeOutDurationMinute int) {
	rctx := context.Background()
	result, _ := json.Marshal(&funcPage)
	if timeOutDurationMinute == 0 {
		RedisDatabase.Set(rctx, "funcpage_"+funcPage.URLAddress+"|"+funcPage.URLMethod, string(result), time.Hour*24)
	} else {
		tmpTimeout := time.Duration(timeOutDurationMinute) * time.Minute
		RedisDatabase.Set(rctx, "funcpage_"+funcPage.URLAddress+"|"+funcPage.URLMethod, string(result), tmpTimeout)
	}
}
func TestRedis() {
	fmt.Println("测试内存数据库")
	rctx := context.Background()
	result, _ := RedisDatabase.Get(rctx, "key_test").Result()
	fmt.Println(result)
	fmt.Println("测试内存数据库完成")
}

func WriteSystemLogToRedis(log *entities.SystemLog) error {
	bytes, err := json.Marshal(log)
	if err == nil {
		rctx := context.Background()
		cmd := RedisDatabase.LPush(rctx, "SystemLog", string(bytes))
		if cmd.Err() != nil {
			err = cmd.Err()
		}
	}
	return err
}

// 设置 UWB 终端标签信息到 Redis
func SetUWBTerminalTagToRedis(uwbTerminal *entities.UWBTerminalMQTTInformation) error {
	rctx := context.Background()
	uwbTerminal.InCacheDateTime = time.Now()
	bytesUWBTermail, err := json.Marshal(uwbTerminal)
	cmd := RedisDatabase.Set(rctx, fmt.Sprintf("UWBTag_%s", uwbTerminal.Properties.SwimmerId), string(bytesUWBTermail), 24*time.Hour)
	if cmd.Err() != nil {
		err = cmd.Err()
	}
	return err
}

func EnumAllUWBTerminalTagFromRedis() ([]entities.UWBTerminalMQTTInformation, error) {
	var uwbTerminals []entities.UWBTerminalMQTTInformation
	rctx := context.Background()
	strUWBTerminals, err := RedisDatabase.Keys(rctx, "UWBTag_*").Result()
	if err == nil {
		for _, strUWBTerminal := range strUWBTerminals {
			var uwbTerminal entities.UWBTerminalMQTTInformation
			result, _ := RedisDatabase.Get(rctx, strUWBTerminal).Result()
			err = json.Unmarshal([]byte(result), &uwbTerminal)
			if err == nil {
				uwbTerminals = append(uwbTerminals, uwbTerminal)
			}
		}
	}
	return uwbTerminals, err
}

// 打开相关 UWB 终端标签缓存接收
func OpenUWBTerminalTagFromRedis(swimmerId string) {
	rctx := context.Background()
	RedisDatabase.Set(rctx, fmt.Sprintf("UWBTag_%s", swimmerId), "open", 24*time.Hour)
}

// 关闭相关 UWB 终端标签缓存接收
func CloseUWBTerminalTagFromRedis(swimmerId string) {
	rctx := context.Background()
	RedisDatabase.Set(rctx, fmt.Sprintf("UWBTag_%s", swimmerId), "close", 24*time.Hour)
}

func GetUWBTerminalTagFromRedis(swimmerId string) (entities.UWBTerminalMQTTInformation, error) {
	var uwbTerminal entities.UWBTerminalMQTTInformation
	var err error
	rctx := context.Background()
	result, _ := RedisDatabase.Get(rctx, fmt.Sprintf("UWBTag_%s", swimmerId)).Result() // 从缓存中获取相关游泳者的UWB标签信息
	if result == "close" {
		// 当缓存中的内容是 close 时，报告 UWB Terminal 缓存服务被关闭（目前主要用在游泳者签出逻辑）
		err = errors.New("the UWB terminal was closed")
	} else if result == "open" {
		// 当缓存中的内容是 open 时， 报告 UWB Terminal 缓存服务被打开，允许将经过处理的 UWB Terminal 信息放入缓存
		err = errors.New("the UWB terminal was opened")
	} else if result == "clear" {
		// 当缓存中的内容是 clear 时， 报告 UWB Terminal 缓存服务的数据被清空。（目前主要用于清空 Web Socket 客户端 UWB 标签信息)
		err = errors.New("the UWB terminal was cleared on web socket client")
	} else {
		err = json.Unmarshal([]byte(result), &uwbTerminal)
	}
	defer rctx.Done()
	return uwbTerminal, err
}

func SetAlertInformationToRedis(alert *entities.AlertInformation) {
	// 设置危险告警到 cache 队列中
	rctx := context.Background()
	RedisDatabase.Set(rctx, fmt.Sprintf("AlertInfor_%s", alert.SwimmerId), alert, 24*time.Hour)
}
func ClearAlertInformationFromRedis(swimmerId string) {
	// 清楚告警信息，并设置五分钟内都忽略告警。
	rctx := context.Background()
	RedisDatabase.Set(rctx, fmt.Sprintf("ClearAlertInfor_%s", swimmerId), "clear", 5*time.Minute)
	RedisDatabase.Del(rctx, fmt.Sprintf("AlertInfor_%s", swimmerId))
}
func EnumAllAlertInformationFromRedis() ([]entities.AlertInformation, error) {
	// 获取所有危险告警
	var alerts []entities.AlertInformation
	var strErrs []string
	rctx := context.Background()
	res, err := RedisDatabase.Keys(rctx, "AlertInfor_*").Result()
	if err == nil {
		for _, re := range res {
			re, err = RedisDatabase.Get(rctx, re).Result()
			if err == nil {

			} else {
				strErrs = append(strErrs, err.Error())
			}
		}
		if len(strErrs) > 0 {
			err = fmt.Errorf(strings.Join(strErrs, ";"))
		}
	}

	return alerts, err

}
func GetTopSiteFence(siteId string) (entities.SiteFence, error) {
	var fence entities.SiteFence
	var err error
	rctx := context.Background()
	key := fmt.Sprintf("SiteTopFence_%s", siteId)
	defer rctx.Done()
	if tmp, err := RedisDatabase.Exists(rctx, key).Result(); nil == err && tmp > 0 {
		sFence := RedisDatabase.Get(rctx, key).Val()
		err = json.Unmarshal([]byte(sFence), &fence)
		return fence, err
	}
	err = fmt.Errorf("not found in cache")
	return fence, err

}

func SetTopSiteFence(fence *entities.SiteFence) {
	rctx := context.Background()
	key := fmt.Sprintf("SiteTopFence_%s", fence.SiteID)
	var bFence []byte
	bFence, _ = json.Marshal(fence)
	RedisDatabase.Set(rctx, key, bFence, 24*time.Hour)
	defer rctx.Done()
}

// TODO: 需要在 BIZ 中编写 UWBTERMINAL 的业务代码。
/*
-------------------------------------------------------------------
1. 写入 Redis
当从 MQTT 服务中获取到最新 UWB Terminal 立刻写入到 Redis，
并将当前距离与之前的 Redis 中的UWB终端信息相加，
注意，只有之前有该 UWB Terminal 信息时才进行加法运算否则只进行记录，该记录每天凌晨清空并记录到数据库进行物理存储。


2. 停止运动告警(普通)
定时获取指定 UWB Terminal 信息（永远保持最后一次），
比对最后一次从MQTT服务中获取的信息的时间差异和距离差异，
如果距离差异小于 N（可配置）米，时间差距 > M 秒并且最后一次出现时在泳池内部，
则发送告警信息到 Redis 并提交到用于警告的 WEB SOCKET 中。
注意这里可能需要用到线程异步。（停止运动告警）

3. 溺水告警(严重)
定时获取所有 UWB Terminal 信息，循环判断每个 UWB Terminal 最后一次的信息与当前时间差 > M 秒并且在泳池内部，则立刻告警。(溺水告警)
-------------------------------------------------------------------

*/
