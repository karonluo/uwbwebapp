package entities

import "time"

type SwimmerSwimmReport struct {
	CreateDatetime time.Time `json:"CreateDatetime"` // 数据创建日期时间
	Creator        string    `json:"Creator"`        // 数据创建人
	SwimmerID      string    `json:"SwimmerId"`      // 游泳者唯一编号
	Distence       float32   `json:"Distence"`       // 游泳距离
	Date           time.Time `json:"Date"`           // 报告日期
	ModifyDatetime time.Time `json:"ModifyDatetime"` // 数据修改日期时间
	Modifier       string    `json:"Modifier"`       // 数据修改人
}

func (SwimmerSwimmReport) TableName() string {
	return "swimmer_swimm_reports"
}
