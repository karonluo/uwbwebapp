package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"

	"github.com/google/uuid"
)

func CreateSwimmer(swimmer entities.Swimmer) (string, error) {

	if swimmer.Id == "" {
		swimmer.Id = uuid.New().String()
	}

	if swimmer.Creator == "" {
		swimmer.Creator = "admin"
	}
	if swimmer.Modifier == "" {
		swimmer.Modifier = "admin"

	}
	swimmer.CreateDatetime = time.Now()
	swimmer.ModifyDatetime = time.Now()
	err := dao.CreateSwimmer(swimmer)
	return swimmer.Id, err
}

func QuerySwimmers(queryCondition entities.QueryCondition) ([]entities.Swimmer, int64, int64, error) {
	var swimmers []entities.Swimmer
	dataRecordCount, err := dao.GetSwimmersCount(queryCondition)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		swimmers, err = dao.QuerySwimmers(queryCondition)
	}

	return swimmers, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func GetSwimmersById(id string) (entities.Swimmer, error) {

	swimmer, err := dao.GetSwimmersById(id)
	return swimmer, err
}

func DeleteSwimmers(ids []string) error {
	return dao.DeleteSwimmers(ids)
}

func UpdateSwimmer(swimmer entities.Swimmer) error {
	var tmpSwimmer entities.Swimmer
	var err error
	tmpSwimmer, err = dao.GetSwimmersById(swimmer.Id)
	if err == nil {
		// 防止以下字段被修改
		swimmer.CreateDatetime = tmpSwimmer.CreateDatetime
		swimmer.Creator = tmpSwimmer.Creator
		swimmer.ModifyDatetime = time.Now()
		if swimmer.Modifier == "" {
			swimmer.Modifier = "admin"
		}
		err = dao.UpdateSwimmer(swimmer)
	}

	return err
}
