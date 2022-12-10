// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    companySite, err := UnmarshalCompanySite(bytes)
//    bytes, err = companySite.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalCompanySite(data []byte) (CompanySite, error) {
	var r CompanySite
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CompanySite) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// CompanySite
type CompanySite struct {
	SportsCompanyName string    `gorm:"column:sports_company_name" json:"SportsCompanyName"` // 体育公司名称
	CreateDatetime    time.Time `json:"CreateDatetime"`                                      // 数据创建日期，前端JSON赋值null
	Creator           string    `json:"Creator"`                                             // 数据创建者
	Modifier          string    `json:"Modifier"`                                            // 数据修改者
	ModifyDatetime    time.Time `json:"ModifyDatetime"`                                      // 数据修改日期，前端JSON赋值null
	SiteDisplayName   string    `json:"SiteDisplayName"`                                     // 场地显示名
	SiteID            string    `json:"SiteId"`                                              // 场地唯一编号
	SportsCompanyID   string    `json:"SportsCompanyId"`                                     // 体育公司唯一编号
}

func (CompanySite) TableName() string {
	return "company_sites"
}
