package biz

import (
	"math"
	"time"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

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

// 返回
// 本次新加入的公司
// 错误信息集合
func SwimmerJoinInSportsCompanies(css []entities.CompanySwimmer, swimmerId string) ([]entities.CompanySwimmer, []error) {
	var err error
	var errs []error
	// var originJoined []entities.CompanySwimmer
	var newJoined []entities.CompanySwimmer
	err = dao.ClearAllCompaniesFromSwimmer(swimmerId)
	if err == nil {
		for _, cs := range css {
			// TODO: 需要优化
			/*
				_, dataRecordCount, err = dao.GetCompanySwimmerByCompanyIDAndSwimmerID(cs.SportsCompanyID, cs.SwimmerID)
				if dataRecordCount == 1 {
					originJoined = append(originJoined, cs)
				} else {
					err = dao.CreateCompanySwimmer(&cs)
					if err == nil {
						newJoined = append(newJoined, cs)
					} else {
						errs = append(errs, err)
					}
				}
			*/
			// 目前的方式，需要前端进行优先判断，只加入曾经未加入的公司。
			err = dao.CreateCompanySwimmer(&cs)
			cs.CreateDatetime = time.Now()
			cs.ModifyDatetime = time.Now()
			if cs.Creator == "" {
				cs.Creator = "admin"
			}
			cs.Modifier = cs.Creator
			if err == nil {
				newJoined = append(newJoined, cs)
			} else {
				errs = append(errs, err)
				tools.ProcessError("biz.SwimmerJoinInSportsCompany", `err = dao.CreateCompanySwimmer(&cs)`, err)

			}
		}
	} else {
		tools.ProcessError("biz.SwimmerJoinInSportsCompany", `err = dao.ClearAllCompaniesFromSwimmer(swimmerId)`, err)
	}
	return newJoined, errs
}

func SwimmerEnumCompanies(swimmerId string) ([]entities.CompanySwimmer, error) {
	return EnumSportsCompanySwimmersBySwimmerId(swimmerId)
}

// 设置游泳者进入场地日期时间
func SetSwimmerEnterToSite(swimmerId string, siteId string, enterTime time.Time) {

}

// 设置游泳者退出场地日期时间
func SetSwimmerExitestFromSite(swimmerId string, siteId string, enterTime time.Time) {

}
