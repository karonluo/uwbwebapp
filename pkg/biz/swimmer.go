package biz

import (
	"fmt"
	"math"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/tools"

	"github.com/google/uuid"
)

func CreateSwimmer(swimmer entities.Swimmer) (string, error) {

	if swimmer.ID == "" {
		swimmer.ID = uuid.New().String()
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
	return swimmer.ID, err
}

// 查询游泳者信息
// 游泳者集合, 页数，数据项总数，错误信息
func QuerySwimmers(queryCondition entities.QueryCondition, companyId string) ([]entities.Swimmer, int64, int64, error) {
	var swimmers []entities.Swimmer
	dataRecordCount, err := dao.GetSwimmersCount(queryCondition, companyId)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		swimmers, err = dao.QuerySwimmers(queryCondition, companyId)
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
	var oriSwimmer entities.Swimmer
	var err error
	oriSwimmer, err = dao.GetSwimmersById(swimmer.ID)
	if err == nil {
		// 防止以下字段被修改
		swimmer.CreateDatetime = oriSwimmer.CreateDatetime
		swimmer.Creator = oriSwimmer.Creator
		swimmer.ModifyDatetime = time.Now()
		if swimmer.Modifier == "" {
			swimmer.Modifier = "admin"
		}
		if swimmer.DisplayName != oriSwimmer.DisplayName {
			err = dao.UpdateSwimmerDisplayNameRelTables(swimmer.ID, swimmer.DisplayName)
		}
		if err == nil {
			err = dao.UpdateSwimmer(swimmer)
			if err == nil {
				// UWBDevicePlatformUpdateTerminal(swimmer.DisplayName, swimmer.Gender, cs.UWBTagCode, conf.WebConfiguration.UWBDevicePlatformConf.DefaultTerminalModelId)
				// 获取该游泳者所有的标签
				css, _ := SwimmerEnumCompanies(swimmer.ID)
				for _, cs := range css {
					UWBDevicePlatformUpdateTerminal(swimmer.DisplayName, swimmer.Gender, swimmer.ID, cs.UWBTagCode, conf.WebConfiguration.UWBDevicePlatformConf.DefaultTerminalModelId)
				}

			} else {
				err = fmt.Errorf("更新失败。")
			}
		}
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
	err = dao.CancelAllUWBTagsBySwimmerId(swimmerId) // 首先取消该游泳者所有标签。
	if err == nil {
		err = dao.ClearAllCompaniesFromSwimmer(swimmerId) // 清楚游泳者和公司的关系
		if err == nil {
			for _, cs := range css {
				err = dao.CreateCompanySwimmer(&cs) // 重新建立公司关系
				cs.CreateDatetime = time.Now()
				cs.ModifyDatetime = time.Now()
				if cs.Creator == "" {
					cs.Creator = "admin"
				}
				cs.Modifier = cs.Creator
				if err == nil {
					// 设置标签绑定状态
					err = dao.SetUWBTagSwimmerBoundStatus(cs.UWBTagCode, true) // 重新绑定 UWB 标签
					if err == nil {
						// 设置UWB管理平台中的终端标签
						UWBDevicePlatformUpdateTerminal(cs.SwimmerDisplayName, cs.SwimmerGender, cs.SwimmerID, cs.UWBTagCode, conf.WebConfiguration.UWBDevicePlatformConf.DefaultTerminalModelId)
					}
					if err == nil {
						newJoined = append(newJoined, cs)
					}
				} else {
					errs = append(errs, err)
					tools.ProcessError("biz.SwimmerJoinInSportsCompany", `err = dao.CreateCompanySwimmer(&cs)`, err)
				}
			}
		} else {
			tools.ProcessError("biz.SwimmerJoinInSportsCompany", `err = dao.ClearAllCompaniesFromSwimmer(swimmerId)`, err)
		}
	} else {
		tools.ProcessError("biz.SwimmerJoinInSportsCompany", `err = dao.CancelAllUWBTagsBySwimmerId(swimmerId)`, err)
	}
	return newJoined, errs
}

func SwimmerEnumCompanies(swimmerId string) ([]entities.CompanySwimmer, error) {
	return EnumSportsCompanySwimmersBySwimmerId(swimmerId)
}

func SetSwimmerVIPLevel(swimmerId string, companyId string, vipLevelDictCode string, vipLevel string, modifier string) error {
	var err error
	var swimmerCompany entities.CompanySwimmer
	var recordCount int64
	swimmerCompany, recordCount, err = dao.GetCompanySwimmerByCompanyIDAndSwimmerID(companyId, swimmerId)
	if recordCount == 1 && err == nil {
		swimmerCompany.VIPLevelDictCode = vipLevelDictCode
		swimmerCompany.VIPLevel = vipLevel
		swimmerCompany.ModifyDatetime = time.Now()
		if modifier != "" {
			swimmerCompany.Modifier = modifier
		} else {
			swimmerCompany.Modifier = "admin"
		}
		err = dao.UpdateCompanySwimmer(&swimmerCompany)
	}
	return err
}
