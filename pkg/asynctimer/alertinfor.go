package asynctimer

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/ahmetb/go-linq/v3"
)

func AlertDangerInformation() {
	for {
		time.Sleep(10 * time.Second)
		// 获取在场游泳者，循环查看在场游泳者最后一次收到信息的区域和时间，如果时间 >30 秒且区域在泳池范围则危险告警。
		calendars, err := biz.EnumAllSwimmerCalendarsOnSite()
		if err == nil {
			for _, cal := range calendars {
				uwbInfor, _ := cache.GetUWBTerminalTagFromRedis(cal.SwimmerID)
				tmp := uwbInfor.InCacheDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_NANO)
				inCacheDateTime, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, tmp)
				tmp = time.Now().Format(tools.GOOGLE_DATETIME_FORMAT_NO_NANO)
				now, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, tmp)
				conf, _ := biz.GetSystemDict("4001")
				threshold, _ := strconv.ParseFloat(conf.Value, 64)
				fmt.Println(conf.Value)
				fmt.Println(threshold)
				if now.Sub(inCacheDateTime).Seconds() > threshold {

					fmt.Printf("Danger: %s\r\n", uwbInfor.Properties.SwimmerDisplayName) // 把数据塞入告警缓存
					var info entities.AlertInformation
					info.DevEui = uwbInfor.DevEui
					info.Message = fmt.Sprintf("危险告警，%s 可能出现溺水，已经 %d 秒未在泳池内检测到该人员。", info.SwimmerDisplayName, int(threshold))
					info.Type = "danger"
					info.X = uwbInfor.X
					info.Y = uwbInfor.Y
					cache.SetAlertInformationToRedis(&info)
				}
			}
		}
	}
}
func AlertNormalInformation() {
	type Infor struct {
		SumDistence        float32
		SwimmerId          string
		SwimmerDisplayName string
	}
	var firstForProcess bool = true
	var infors []Infor
	for {
		if firstForProcess {
			tags, er := cache.EnumAllUWBTerminalTagFromRedis()
			if er == nil {
				for _, tag := range tags {
					if tag.Properties.SwimmerId != "" {
						infors = append(infors, Infor{SumDistence: tag.SumDistance, SwimmerId: tag.Properties.SwimmerId, SwimmerDisplayName: tag.Properties.SwimmerDisplayName})
					}
				}
			}
			firstForProcess = false
		} else {
			tags, er := cache.EnumAllUWBTerminalTagFromRedis()
			if er == nil {
				for _, tag := range tags {
					infor := linq.From(infors).WhereT(func(t *Infor) bool {
						return (t.SwimmerId == tag.Properties.SwimmerId)
					}).First().(*Infor)
					if math.Abs(float64(infor.SumDistence-tag.SumDistance)) < 5.0 {
						// 一分钟内移动距离少于5米，则告警。
						fmt.Printf("Alert %s\r\n", infor.SwimmerDisplayName)
					}
				}
			}
		}
		time.Sleep(1 * time.Minute) // 定时一分钟检测一次
	}
}

// 每日 23:59 分统计当日所有游泳者距离
func SumSwimmerDistance() {
	stimeout := "23:59"
	for {
		st := time.Now().Format("15:04")
		time.Sleep(500 * time.Millisecond)
		if stimeout == st {
			// 获取所有已入场但未出场的游泳者
			cals, err := biz.EnumAllSwimmerCalendarsOnSite()
			if err == nil {
				for _, cal := range cals {
					cache.GetUWBTerminalTagFromRedis(cal.SwimmerID)
					cal.ExitDateTime = time.Now()
					biz.SwimmerExitFromSite(&cal) // 签出
				}
			}
		}
	}
}
