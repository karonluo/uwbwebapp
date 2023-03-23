// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    site, err := UnmarshalSite(bytes)
//    bytes, err = site.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

// site
type Site struct {
	Address        string    `json:"Address"` // 场地地址
	Contact        string    `json:"Contact"` // 场地联系方式，填写手机或电话以逗号分割
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	Id             string    `gorm:"primaryKey" json:"Id"` // 唯一编号
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
	Users          string    `gorm:"size:128" json:"Users"` // 在绑定系统用户时自动填写，用户不能在界面上编辑文字的方式修改，多个用户显示名以逗号分割。
	DisplayName    string    `json:"DisplayName"`
	GPS            string    `gorm:"column:gps" json:"GPS"` // 场地GPS地址 经度纬度标识法例如：103.975466,30.680291
}

// 指定表名 为 sites
func (Site) TableName() string {
	return "sites"
}

func (r Site) UnmarshalSite(data []byte) (Site, error) {
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Site) SiteMarshal() ([]byte, error) {
	return json.Marshal(r)
}
