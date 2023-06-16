package biz

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"
)

func GetUWBBaseStationByCode(id string) (entities.UWBBaseStation, error) {
	return dao.GetUWBBaseStationByCode(id)
}

func EnumUWBBaseStationBySiteId(siteId string) ([]entities.UWBBaseStation, error) {
	return dao.EnumUWBBaseStationBySiteId(siteId)
}

func InitAllUWBBaseStationsToRedis() {
	fmt.Print("将数据库中的所有基站信息放入内存数据库")
	var queryCodition entities.QueryCondition
	queryCodition.LikeValue = ""
	queryCodition.PageIndex = 1
	queryCodition.PageSize = 999999
	stations, _, _, _ := QueryUWBBaseStations(queryCodition)
	cache.SetAllUWBBaseStationsToRedis(stations)
	fmt.Println("......完成")

}
func GetTopSiteFence(siteId string) (entities.SiteFence, error) {
	fence, err := cache.GetTopSiteFence(siteId)
	if err != nil {
		fence, err = dao.GetTopSiteFence(siteId)
		if err == nil {
			cache.SetTopSiteFence(&fence)
		}
	}
	return fence, err
}

// 生成 UWB 基站电子围栏多边形范围
func GenerateUWBBaseStationPointsBySiteId(siteId string) []entities.Point {
	// 优先获取电子围栏信息，当没有电子围栏信息时使用 UWB 基站点位信息组合
	fence, err := GetTopSiteFence(siteId)
	var points []entities.Point
	if err == nil {
		json.Unmarshal([]byte(fence.Coordinate), &points)
	} else {
		stations := cache.EnumAllUWBBaseStationsBySiteIdFromRedis(siteId)
		for _, station := range stations {
			var point entities.Point
			tmpPoint := strings.Split(station.Position, ",")
			if v, e := strconv.ParseFloat(tmpPoint[0], 64); e == nil {
				point.X = v
			}
			if v, e := strconv.ParseFloat(tmpPoint[1], 64); e == nil {
				point.Y = v
			}
			points = append(points, point)
		}
	}
	return points
}

func CreateUWBBaseStation(station *entities.UWBBaseStation) error {
	station.CreateDatetime = time.Now()
	if station.Creator == "" {
		station.Creator = "admin"
	}
	station.ModifyDatetime = station.CreateDatetime
	station.Modifier = station.Creator
	err := dao.CreateUWBBaseStation(station)
	if err == nil {
		// 更新缓存
		cache.SetUWBBaseStationToRedis(*station)
	}
	return err
}

func QueryUWBBaseStations(queryCondition entities.QueryCondition) ([]entities.UWBBaseStation, int64, int64, error) {
	var station []entities.UWBBaseStation
	dataRecordCount, err := dao.GetUWBBaseStationCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		station, err = dao.QueryUWBBaseStations(queryCondition)
	}

	return station, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func DeleteUWBBaseStations(codes []string) error {
	return dao.DeleteUWBBaseStations(codes)
}

func UpdateUWBBaseStation(station *entities.UWBBaseStation) error {
	var tmpStation entities.UWBBaseStation
	var err error
	tmpStation, err = dao.GetUWBBaseStationByCode(station.Code)
	if err == nil {
		// 防止以下字段被修改
		station.CreateDatetime = tmpStation.CreateDatetime
		station.Creator = tmpStation.Creator
		station.ModifyDatetime = time.Now()
		if station.Modifier == "" {
			station.Modifier = "admin"
		}
		err = dao.UpdateUWBBaseStation(station)
		if err == nil {
			// 更新缓存
			cache.SetUWBBaseStationToRedis(*station)
		}
	}
	return err
}

func CheckPointInSite(uwbInfor entities.UWBTerminalMQTTInformation, siteId string) bool {
	result := false
	//! 验证的时候打开 return true
	//return true
	if points := GenerateUWBBaseStationPointsBySiteId(siteId); len(points) >= 3 {
		var point entities.Point
		point.X = uwbInfor.Z //! 注意 UWB 的 Z是X，Y是Y，X未启用
		point.Y = uwbInfor.Y
		debug := fmt.Sprintf("测试: %s 人员位置: X:%f, Y:%f", uwbInfor.Properties.SwimmerDisplayName, point.X, point.Y)
		result, _ = tools.UseGeosInPolygon(point, points)
		if result {
			fmt.Println(debug + " 在基站范围内")
		} else {
			fmt.Println(debug + " 不在基站范围内")
		}
	}
	return result
}
