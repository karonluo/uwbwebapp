// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    siteFence, err := UnmarshalSiteFence(bytes)
//    bytes, err = siteFence.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSiteFence(data []byte) (SiteFence, error) {
	var r SiteFence
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SiteFence) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SiteFence
type SiteFence struct {
	Code           string    `json:"Code"` // 泳池电子围栏唯一编号，目前后端随机填写，同时支持前端输入。
	Coordinate     string    `json:"Coordinate"`
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
	SiteID         string    `json:"SiteId"` // 场地唯一编号
}

func (SiteFence) TableName() string {
	return "site_fence"
}
