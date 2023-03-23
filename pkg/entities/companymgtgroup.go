// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sportsCompanyMgtGroup, err := UnmarshalSportsCompanyMgtGroup(bytes)
//    bytes, err = sportsCompanyMgtGroup.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSportsCompanyMgtGroup(data []byte) (SportsCompanyMgtGroup, error) {
	var r SportsCompanyMgtGroup
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SportsCompanyMgtGroup) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SportsCompanyMgtGroup
type SportsCompanyMgtGroup struct {
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	Description    string    `json:"Description"` // 描述
	Id             string    `json:"Id"`          // 体育公司管理组唯一编号
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
	Name           string    `json:"Name"`                           // 名称
	IsBound        bool      `grom:"column:is_bound" json:"IsBound"` // 是否已经绑定公司 冗余字段
}
