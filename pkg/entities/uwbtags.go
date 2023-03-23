// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    uWBTag, err := UnmarshalUWBTag(bytes)
//    bytes, err = uWBTag.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalUWBTag(data []byte) (UWBTag, error) {
	var r UWBTag
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UWBTag) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UWBTag
type UWBTag struct {
	Code              string    `gorm:"primaryKey" json:"Code"` // 需要手工填写的唯一编码
	CreateDatetime    time.Time `json:"CreateDatetime"`         // 后端干预前端不用填写
	Creator           string    `json:"Creator"`                // 前端填写
	Description       string    `json:"Description"`            // 描述信息128个字符
	Modifier          string    `json:"Modifier"`               // 前端填写
	ModifyDatetime    time.Time `json:"ModifyDatetime"`         // 后端干预前端不用填写
	SportsCompanyID   string    `json:"SportsCompanyId"`        // 体育公司唯一编码
	SportsCompanyName string    `json:"SportsCompanyName"`      // 体育公司名称
	IsBound           bool      `json:"IsBound"`                // 是否已经绑定了游泳者会员
}

func (me *UWBTag) String() string {
	return me.Code
}

func (UWBTag) TableName() string {
	return "uwb_tags"
}
