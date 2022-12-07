// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    siteUser, err := UnmarshalSiteUser(bytes)
//    bytes, err = siteUser.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSiteUser(data []byte) (SiteUser, error) {
	var r SiteUser
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SiteUser) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SiteUser
type SiteUser struct {
	CreateDatetime     time.Time `json:"CreateDatetime"`     // 数据创建日期
	Creator            string    `json:"Creator"`            // 数据创建人
	JobTitle           string    `json:"JobTitle"`           // 系统字典-岗位
	JobTitleDictCode   string    `json:"JobTitleDictCode"`   // 字典编码
	Modifier           string    `json:"Modifier"`           // 数据修改人
	ModifyDatetime     time.Time `json:"ModifyDatetime"`     // 数据修改日期
	SiteDisplayName    string    `json:"SiteDisplayName"`    // 场地显示名
	SiteID             string    `json:"SiteId"`             // 场地唯一编号
	SysUserDisplayname string    `json:"SysUserDisplayname"` // 用户显示名
	SysUserID          string    `json:"SysUserId"`          // 系统用户唯一编号
}
