package biz

import (
	"encoding/json"
	"fmt"
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/ahmetb/go-linq/v3"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

// 通过游泳者唯一编号和未出场状态获取日程信息
func GetSiteIdBySwimmerIdAndNoExitSite(swimmerId string) (string, error) {
	var err error
	siteId := ""
	rctx := context.Background()
	key := fmt.Sprintf("SwimmerOnSite_%s", swimmerId)
	if cache.RedisDatabase.Exists(rctx, key).Val() > 0 {
		siteId = cache.RedisDatabase.Get(rctx, key).Val()
		if siteId == "" {
			err = fmt.Errorf("未找到该游泳者所在泳池")
		}
	} else {
		var cal entities.SwimmerCalendar
		cal, err = dao.GetSwimmerCalendarBySwimmerIdAndNoExitSite(swimmerId)
		if err == nil {
			cache.RedisDatabase.Set(rctx, key, cal.SiteID, 24*time.Hour)
		}
	}
	return siteId, err
}

func EnumSwimmerCalendarByDateScope(siteId string, date entities.BetweenDatetime) ([]entities.SwimmerCalendar, int64, error) {
	return dao.EnumSwimmerCalendarByDateScope(siteId, date)
}

// 获取所有在场游泳者
func EnumAllSwimmerCalendarsOnSite() ([]entities.SwimmerCalendar, error) {
	return dao.EnumAllSwimmerCalendarsOnSite()
}

// 获取所有指定场所的在场游泳者
func EnumAllSwimmerCalendarsOnSiteBySiteId(siteId string) ([]entities.SwimmerCalendar, error) {
	return dao.EnumAllSwimmerCalendarsOnSiteBySiteId(siteId)
}

// 判断游泳者是否在场
func SwimmerInSite(siteId string, swimmerId string, date time.Time) (bool, error) {
	betweenDate := tools.SplitDateTimeToAllDay(date)
	calendars, count, errq := dao.EnumSwimmerCalendarBySwimmerDateScope(siteId, swimmerId, betweenDate)
	var inSite bool = false
	if errq == nil && count > 0 {
		c := linq.From(calendars).WhereT(func(t *entities.SwimmerCalendar) bool {
			return t.EnterDatetime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01" &&
				t.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01"
		}).Count()
		inSite = (c > 0)
	}
	return inSite, errq
}

// 前台人员登记游泳者入场(前台无需告知入场日期时间，系统会自动记录入场日期时间)
func SwimmerEnterToSite(calendar *entities.SwimmerCalendar) ([]string, []string) {
	var err error
	var ids []string
	var errs []string
	var timeNow time.Time
	if tools.CheckNoDate(calendar.EnterDatetime) {
		timeNow = time.Now()
	} else {
		timeNow = calendar.EnterDatetime
	}
	betweenDate := tools.SplitDateTimeToAllDay(timeNow) // 分成全天日期时间段
	// 获取该日期该游泳者在该场地的数据。
	calendars, count, errq := dao.EnumSwimmerCalendarBySwimmerDateScope(calendar.SiteID, calendar.SwimmerID, betweenDate)
	if errq != nil {
		errs = append(errs, errq.Error())
	} else {
		if count == 0 {
			// 证明指定用户不再计划内或者还未入场。
			fmt.Println("--------------------------------------------------")
			calendar.CreateDatetime = time.Now()
			calendar.ModifyDatetime = calendar.CreateDatetime
			calendar.EnterDatetime = timeNow
			if calendar.Creator == "" {
				calendar.Creator = "admin"
			}
			calendar.Modifier = calendar.Creator
			calendar.ID = uuid.NewString()
			err = dao.CreateSwimmerCalendar(calendar)
			rctx := context.Background()
			cache.RedisDatabase.Del(rctx, fmt.Sprintf("UWBTag_%s", calendar.SwimmerID)) // 若有该标签数据则删除。
			defer rctx.Done()
			if err != nil {
				errs = append(errs, err.Error())
			} else {
				ids = append(ids, calendar.ID)
			}
		} else {

			// 查询这些数据中有已入场未出场的数据
			if linq.From(calendars).WhereT(
				func(c entities.SwimmerCalendar) bool {
					r := c.EnterDatetime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01" &&
						c.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01"
					return r
				}).Count() > 0 {
				errs = append(errs, "尚有已入场但未登记出场的数据，请先做游泳者退场登记。")

			} else {
				var tmpCalendars []entities.SwimmerCalendar
				// 首先判断每条数据是否都是有入场和退场的，如果都有则需要新建一条入场登记信息。
				if linq.From(calendars).WhereT(func(c entities.SwimmerCalendar) bool {
					r := c.EnterDatetime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01" &&
						c.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01"
					return r
				}).Count() == len(calendars) {
					calendar.CreateDatetime = time.Now()
					calendar.ModifyDatetime = calendar.CreateDatetime
					calendar.EnterDatetime = timeNow
					if calendar.Creator == "" {
						calendar.Creator = "admin"
					}
					calendar.Modifier = calendar.Creator
					calendar.ID = uuid.NewString()
					err = dao.CreateSwimmerCalendar(calendar)
					rctx := context.Background()
					cache.RedisDatabase.Del(rctx, fmt.Sprintf("UWBTag_%s", calendar.SwimmerID)) // 若有该标签数据则删除。
					defer rctx.Done()
					if err != nil {
						errs = append(errs, err.Error())
					} else {
						ids = append(ids, calendar.ID)
					}
				} else {

					// 找到所有未入场也未退场的数据，并将其统一都设置入场日期时间。
					linq.From(calendars).WhereT(func(c entities.SwimmerCalendar) bool {
						r := c.EnterDatetime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01" &&
							c.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01"
						return r
					}).ToSlice(&tmpCalendars)
					for _, tmpCalendar := range tmpCalendars {
						tmpCalendar.EnterDatetime = timeNow
						tmpCalendar.ModifyDatetime = time.Now()
						if calendar.Modifier == "" {
							calendar.Modifier = "admin"
						}
						tmpCalendar.Modifier = calendar.Modifier
						err = dao.UpdateSwimmerCalendar(&tmpCalendar)
						if err != nil {
							errs = append(errs, err.Error())
						} else {
							ids = append(ids, tmpCalendar.ID)
						}
					}
					rctx := context.Background()
					cache.RedisDatabase.Del(rctx, fmt.Sprintf("UWBTag_%s", calendar.SwimmerID)) // 若有该标签数据则删除。
					defer rctx.Done()
				}
			}
		}
	}
	return ids, errs
}

// 前台人员登记游泳者出场，由于有特殊性，因此允许前台人员设定签出日期时间，但会根据输入的签出的日期，
func SwimmerExitFromSite(calendar *entities.SwimmerCalendar) ([]string, []string) {
	// var err error
	var ids, errs []string
	var timeNow time.Time
	if tools.CheckNoDate(calendar.ExitDateTime) {
		timeNow = time.Now()
	} else {
		timeNow = calendar.ExitDateTime
	}
	betweenDate := tools.SplitDateTimeToAllDay(timeNow) // 分成全天日期时间段
	// 查询当日该游泳者的所有日程
	calendars, count, errq := dao.EnumSwimmerCalendarBySwimmerDateScope(calendar.SiteID, calendar.SwimmerID, betweenDate)
	if errq != nil {
		// 当有错误则将错误加入到错误列表中
		errs = append(errs, errq.Error())
	} else {
		if count == 0 {
			// 当没有找到任何日历时，报告该游泳者未入场不允许签出。
			errs = append(errs, "该游泳者目前未入场，无法签出。")
		} else if count > 0 {
			// 判断所有数据都有签退日期时间, 如果都有签退时间证明不能进行签退。
			if linq.From(calendars).WhereT(func(c entities.SwimmerCalendar) bool {
				return c.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01"
			}).Count() == len(calendars) {
				errs = append(errs, "该游泳者目前未入场，无法签出。")
			} else {
				var err error
				var tmpCalendars []entities.SwimmerCalendar
				linq.From(calendars).WhereT(func(c entities.SwimmerCalendar) bool {
					return (c.EnterDatetime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) != "0001-01-01" &&
						c.ExitDateTime.Format(tools.GOOGLE_DATETIME_FORMAT_NO_TIME) == "0001-01-01")
				}).ToSlice(&tmpCalendars)
				// 找到相关 UWB Terminal 标签，获取最后的距离
				var infor entities.UWBTerminalMQTTInformation
				infor, _ = cache.GetUWBTerminalTagFromRedis(calendar.SwimmerID)
				cache.CloseUWBTerminalTagFromRedis(calendar.SwimmerID) // 关闭禁止更新
				// 更新所有有入场并且无签退日期时间的数据，统一签退时间。
				for _, tmpCalendar := range tmpCalendars {
					tmpCalendar.ExitDateTime = timeNow
					tmpCalendar.ModifyDatetime = time.Now()
					if calendar.Modifier == "" {
						calendar.Modifier = "admin"
					}
					tmpCalendar.Modifier = calendar.Modifier
					tmpCalendar.TotalMileage = tmpCalendar.TotalMileage + infor.SumDistance // 将标签中的距离信息录入到数据库中
					err = dao.UpdateSwimmerCalendar(&tmpCalendar)
					if err != nil {
						errs = append(errs, err.Error())
					} else {
						// 这里清空关于入场游泳者场地编号相关缓存。
						rctx := context.Background()
						key := fmt.Sprintf("SwimmerOnSite_%s", tmpCalendar.SwimmerID)
						cache.RedisDatabase.Del(rctx, key) // 若有该标签数据则删除。
						defer rctx.Done()
						ids = append(ids, tmpCalendar.ID)
					}
				}

			}

		} else {
			errs = append(errs, "数据异常，无法定位游泳者入场日程")
		}
	}
	return ids, errs
}

// TODO: 需要增加取消周期性的训练计划
// 取消单次训练计划
func SwimmerCalendarPlanCancel(calendarIds []string) ([]string, []string) {
	var ids, errs []string
	for _, id := range calendarIds {
		calendar, err := dao.GetSwimmerCalendarById(id)
		if err != nil {
			errs = append(errs, err.Error())
		} else {
			if tools.CheckNoDate(calendar.EnterDatetime) {
				err = dao.DeleteSwimmerCalendar(&calendar)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					ids = append(ids, id)
				}
			} else {
				calendar.TrainingBeginDatetime, _ = time.Parse(tools.GOOGLE_DATETIME_FORMAT, "0001-01-01 00:00:00.000000000")
				calendar.TrainingEndDatetime, _ = time.Parse(tools.GOOGLE_DATETIME_FORMAT, "0001-01-01 00:00:00.000000000")
				calendar.ModifyDatetime = time.Now()
				calendar.Modifier = "admin"
				err = dao.UpdateSwimmerCalendar(&calendar)
				if err != nil {
					errs = append(errs, err.Error())
				} else {
					ids = append(ids, id)
				}
			}
		}

	}
	return ids, errs
}

// TODO: 需要增加制定周期性的训练计划
// 单次训练计划
func SwimmerCalendarPlanToSite(swimmers []entities.Swimmer, siteId string, beginDate string, beginTime string, endTime string) ([]string, []string) {
	var ids, errs []string
	for _, swimmer := range swimmers {
		beginDatetime, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, fmt.Sprintf("%s %s", beginDate, beginTime))
		endDatetime, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, fmt.Sprintf("%s %s", beginDate, endTime))
		var td entities.BetweenDatetime
		td.BeginDatetime = beginDatetime
		td.EndDatetime = endDatetime
		calendars, count, _ := dao.EnumSwimmerCalendarBySwimmerDateScope(siteId, swimmer.ID, td) // 查询计划训练时间中有没有日程信息（包括入场和出场日期）。
		if count == 0 {
			// 如果没有相关日程则直接新增计划
			var calendar entities.SwimmerCalendar
			calendar.CreateDatetime = time.Now()
			if swimmer.Modifier == "" {
				calendar.Modifier = "admin"
				calendar.Creator = "admin"
			} else {
				calendar.Modifier = "admin"
				calendar.Creator = "admin"
			}
			calendar.ID = uuid.NewString()
			calendar.CreateDatetime = time.Now()
			ids = append(ids, calendar.ID)
			calendar.IsCycle = false
			calendar.ModifyDatetime = calendar.CreateDatetime
			calendar.SiteID = siteId
			calendar.SwimmerID = swimmer.ID
			calendar.SwimmerDisplayName = swimmer.DisplayName
			calendar.TotalMileage = 0
			calendar.TrainingBeginDatetime = td.BeginDatetime
			calendar.TrainingEndDatetime = td.EndDatetime
			err := dao.CreateSwimmerCalendar(&calendar)
			if err != nil {
				errs = append(errs, err.Error())
			}

		} else {
			// 判断是否在这个事件范围有训练计划的，如果有则直接返回错误，错误信息是：请先取消原有计划，如果没有找到则可以在所有记录中加入训练计划时间开始和结束时间。
			if linq.From(calendars).WhereT(func(c entities.SwimmerCalendar) bool {
				result, _ := tools.CheckTimesHasOverlap(c.TrainingBeginDatetime, c.TrainingEndDatetime, td.BeginDatetime, td.EndDatetime)
				return result
			}).Count() > 0 {
				errs = append(errs, fmt.Sprintf("[游泳者编号: %s, 游泳者姓名:%s] 在计划训练的时间范围内，已有安排训练无法创建计划请先取消原有计划。", swimmer.ID, swimmer.DisplayName))

			} else {
				for _, tmpCalendar := range calendars {
					tmpCalendar.TrainingBeginDatetime = beginDatetime
					tmpCalendar.TrainingEndDatetime = endDatetime
					if swimmer.Modifier == "" {
						swimmer.Modifier = "admin"
					}
					tmpCalendar.Modifier = swimmer.Modifier
					tmpCalendar.ModifyDatetime = time.Now()
					err := dao.UpdateSwimmerCalendar(&tmpCalendar)
					if err != nil {
						errs = append(errs, err.Error())
					} else {
						ids = append(ids, tmpCalendar.ID)
					}
				}
			}
		}
	}

	return ids, errs
}

func EnumSwimmerCalendarReportForSwimmer(swimmerId string, date entities.BetweenDatetime) ([]entities.SwimmerCalendar, error) {
	return dao.EnumSwimmerCalendarReportForSwimmer(swimmerId, date)
}

func GetSiteFenceForSwimmer(swimmerId string) ([]entities.Point, error) {
	var err error
	var fence entities.SiteFence
	var points []entities.Point
	// 数据库方式 dao.Database.Raw("")
	// cache 方式
	sql := `SELECT
	* 
FROM
	site_fence 
WHERE
site_id IN ( 
	SELECT site_id FROM swimmer_calendars WHERE 
		enter_datetime <> '0001-01-01 00:00:00' AND 
		exit_datetime = '0001-01-01 00:00:00' AND 
		swimmer_id = ?
		ORDER BY modify_datetime DESC LIMIT 1 ) ORDER BY modify_datetime DESC LIMIT 1`
	err = dao.Database.Raw(sql, swimmerId).Find(&fence).Error
	if err == nil {
		err = json.Unmarshal([]byte(fence.Coordinate), &points)
	}
	return points, err
}
