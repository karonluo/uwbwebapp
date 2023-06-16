package asynctimer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/ahmetb/go-linq/v3"
)

// 返回最后一次接收到的标签信息是否大于告警阈值
func lastedRecvUWBInformationThanThreshold(uwbInfor entities.UWBTerminalMQTTInformation, threshold float64) bool {
	tmp := uwbInfor.InCacheDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_NANO)
	inCacheDateTime, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, tmp)
	tmp = time.Now().Format(tools.GOOGLE_DATETIME_FORMAT_NO_NANO)
	now, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, tmp)
	return now.Sub(inCacheDateTime).Seconds() > threshold
}

func AlertDangerInformation() {
	for {
		time.Sleep(10 * time.Second)
		// 获取在场游泳者，循环查看在场游泳者最后一次收到信息的区域和时间，如果时间 >30 秒且区域在泳池范围则危险告警。
		calendars, err := biz.EnumAllSwimmerCalendarsOnSite()
		if err == nil {
			for _, cal := range calendars {
				uwbInfor, _ := cache.GetUWBTerminalTagFromRedis(cal.SwimmerID) // 从缓存中找到标签信息
				//! 注意 UWB 标签，目前 Z 为 X，Y 为 Y，X 未启用
				timeoutConf, _ := biz.GetSystemDict("4001") // 字典 4001 是危险告警的 timeout
				threshold, _ := strconv.ParseFloat(timeoutConf.Value, 64)
				conf_danger_alert_template, _ := biz.GetSystemDict("4003") // 字典 4003 是危险告警的模板
				// 最后一次收到的UWB信息超过了告警超时阈值并且该UWB消息在场地中，此时可能产生溺水风险。
				if lastedRecvUWBInformationThanThreshold(uwbInfor, threshold) && biz.CheckPointInSite(uwbInfor, cal.SiteID) {
					var info entities.AlertInformation
					info.DevEui = uwbInfor.DevEui
					msg := strings.ReplaceAll(conf_danger_alert_template.Value, "[人员]", info.SwimmerDisplayName)
					msg = strings.ReplaceAll(msg, "[时间]", fmt.Sprintf("%d", int(threshold)))
					info.Message = msg
					info.Type = "danger"
					info.X = uwbInfor.Z
					info.Y = uwbInfor.Y
					cache.SetAlertInformationToRedis(&info) // 将告警信息插入告警缓存，让客户端通过 websocket 获取
					fmt.Println(msg)
					// biz.SpeakMessageAlert()
					//TODO: 调用音响功能
					//biz.SpeakMessage(info.Message)
				}
			}
		}
	}
}

func AlertNormalInformation() {
	type Infor struct {
		SumDistence        float64
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
