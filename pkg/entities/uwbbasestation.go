// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBBaseStation, err := UnmarshalUWBBaseStation(bytes)
//    bytes, err = uWBBaseStation.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalUWBBaseStation(data []byte) (UWBBaseStation, error) {
	var r UWBBaseStation
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBBaseStation) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UWBBaseStation
type UWBBaseStation struct {
	Code            string    `gorm:"primaryKey" json:"Code"`       // 基站编码
	CreateDatetime  time.Time `json:"CreateDatetime"`               // 数据创建日期时间，不用填写，后端自动填写。
	Creator         string    `json:"Creator"`                      // 数据创建人，若不填写则后端填写成admin
	Description     string    `json:"Description"`                  // 描述
	Gps             string    `json:"GPS"`                          // 基站GPS地址
	Modifier        string    `json:"Modifier"`                     // 数据修改人，若不填写则后端填写成admin
	ModifyDatetime  time.Time `json:"ModifyDatetime"`               // 数修改日期时间，不用填写，后端自动填写。
	SiteID          string    `gorm:"column:site_id" json:"SiteId"` // 场地唯一编号
	SiteDisplayName string    `json:"SiteDisplayName"`              // 场地显示名
	Position        string    `json:"Position"`                     // 位置信息 x,y 坐标系 100,200 代表 绘制框中的x和y坐标
}

func (UWBBaseStation) TableName() string {
	return "uwb_base_stations"
}
