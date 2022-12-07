// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    companySwimmer, err := UnmarshalCompanySwimmer(bytes)
//    bytes, err = companySwimmer.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalCompanySwimmer(data []byte) (CompanySwimmer, error) {
	var r CompanySwimmer
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CompanySwimmer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// CompanySwimmer
type CompanySwimmer struct {
	CreateDatetime     time.Time `json:"CreateDatetime"`
	Creator            string    `json:"Creator"`
	MemberCode         string    `json:"MemberCode"`         // 会员编号
	MembershipDatetime time.Time `json:"MembershipDatetime"` // 入会日期
	Modifier           string    `json:"Modifier"`
	ModifyDatetime     time.Time `json:"ModifyDatetime"`
	SportsCompanyID    string    `json:"SportsCompanyId"`                                    // 体育公司唯一编码
	SwimmerID          string    `json:"SwimmerId"`                                          // 游泳者(会员)唯一编码
	UWBTagCode         string    `gorm:"column:uwb_tag_code" json:"UWBTagCode"`              // UWB标签编码
	VIPLevel           string    `gorm:"column:vip_level" json:"VIPLevel"`                   // 会员等级字典值
	VIPLevelDictCode   string    `gorm:"column:vip_level_dict_code" json:"VIPLevelDictCode"` // 会员等级字典编码
	WithdrawaDatetime  time.Time `json:"WithdrawaDatetime"`                                  // 退会日期
	SwimmerDisplayName string    `json:"SwimmerDisplayName"`                                 // 游泳者显示姓名
	SportsCompanyName  string    `json:"SportsCompanyName"`                                  // 体育运动公司名称
}

func (CompanySwimmer) TableName() string {
	return "company_swimmers"
}
