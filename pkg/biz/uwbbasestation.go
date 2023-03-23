package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
)

func GetUWBBaseStationByCode(id string) (entities.UWBBaseStation, error) {
	return dao.GetUWBBaseStationByCode(id)
}

func CreateUWBBaseStation(station *entities.UWBBaseStation) error {
	station.CreateDatetime = time.Now()
	if station.Creator == "" {
		station.Creator = "admin"
	}
	station.ModifyDatetime = station.CreateDatetime
	station.Modifier = station.Creator
	err := dao.CreateUWBBaseStation(station)
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
	}
	return err
}
