// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sysFuncPages, err := UnmarshalSysFuncPages(bytes)
//    bytes, err = sysFuncPages.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSysFuncPage(data []byte) (SysFuncPage, error) {
	var r SysFuncPage
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SysFuncPage) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// SysFuncPages
type SysFuncPage struct {
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	DisplayName    string    `json:"DisplayName"` // 菜单或功能显示名
	ID             string    `json:"Id"`          // 菜单功能唯一编号
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
	ParentID       string    `json:"ParentId"`                       // 上级编号
	URLAddress     string    `json:"UrlAddress"`                     // 调用URL地址
	URLMethod      string    `json:"UrlMethod"`                      // 调用URL地址方法
	URLType        string    `gorm:"column:url_type" json:"URLType"` // 地址类型
	OrderKey       int       `json:"OrderKey"`                       // 排序键
}

func (SysFuncPage) TableName() string {
	return "sys_func_pages"
}
