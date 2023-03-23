// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    siteEnvCalendar, err := UnmarshalSiteEnvCalendar(bytes)
//    bytes, err = siteEnvCalendar.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSiteEnvCalendar(data []byte) (SiteEnvCalendar, error) {
	var r SiteEnvCalendar
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SiteEnvCalendar) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SiteEnvCalendar
type SiteEnvCalendar struct {
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	Date           string    `json:"Date"`                               // 日期
	WaterTemp      float32   `gorm:"column:water_temp" json:"WaterTemp"` // 水温
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
	Ph             float32   `json:"Ph"`                           // 酸碱度
	SiteID         string    `gorm:"column:site_id" json:"SiteId"` // 场地唯一编号
	Temp           float32   `json:"Temp"`                         // 温度
	Clarity        float32   `json:"Clarity"`                      //清澈度
}

func (SiteEnvCalendar) TableName() string {
	return "site_env_calendars"
}
