// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    siteUser, err := UnmarshalSiteUser(bytes)
//    bytes, err = siteUser.Marshal()

package entities

import "encoding/json"

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
	CreateDatetime *string `json:"CreateDatetime"` // 数据创建日期
	Creator        string  `json:"Creator"`        // 数据创建人
	JobTitle       string  `json:"JobTitle"`       // 系统字典-岗位
	Modifier       string  `json:"Modifier"`       // 数据修改人
	ModifyDatetime string  `json:"ModifyDatetime"` // 数据修改日期
	SiteID         string  `json:"SiteId"`         // 场地唯一编号
	SysUserID      string  `json:"SysUserId"`      // 系统用户唯一编号
}
