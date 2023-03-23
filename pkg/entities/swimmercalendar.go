// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    swimmerCalendar, err := UnmarshalSwimmerCalendar(bytes)
//    bytes, err = swimmerCalendar.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSwimmerCalendar(data []byte) (SwimmerCalendar, error) {
	var r SwimmerCalendar
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SwimmerCalendar) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SwimmerCalendar
type SwimmerCalendar struct {
	CreateDatetime        time.Time `json:"CreateDatetime"`
	Creator               string    `json:"Creator"`
	EnterDatetime         time.Time `gorm:"column:enter_datetime" json:"EnterDatetime"` // 进入场地的日期时间，由前台人员登记
	ExitDateTime          time.Time `gorm:"column:exit_datetime" json:"ExitDateTime"`   // 退出场地的日期时间，由前台人员登记
	ID                    string    `gorm:"primaryKey column:id" json:"Id"`
	Modifier              string    `json:"Modifier"`
	ModifyDatetime        time.Time `json:"ModifyDatetime"`
	SiteID                string    `json:"SiteId"`
	SwimmerID             string    `json:"SwimmerId"`
	SwimmerDisplayName    string    `gorm:"column:swimmer_display_name" json:"SwimmerDisplayName"`
	TotalMileage          float32   `json:"TotalMileage"`                   // 总里程，当前台人员登记出场时即时计算本次游泳总里程
	TrainingBeginDatetime time.Time `json:"TrainingBeginDatetime"`          // 训练开始日期时间，由教练指定
	TrainingEndDatetime   time.Time `json:"TrainingEndDatetime"`            // 训练结束日期时间，由教练指定(注意不能超过开始日期）
	IsCycle               bool      `gorm:"column:is_cycle" json:"IsCycle"` // 如果是训练计划记录，需要指定是否为周期性的。
	// Distance              float64   `gorm:"distance" json:"Distance"`       // 本次进场和退场游泳距离。
}
