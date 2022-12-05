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

// SysFuncPage
type SysFuncPage struct {
	CreateDatetime time.Time `json:"CreateDatetime"` // 数据创建日期时间
	Creator        string    `json:"Creator"`        // 数据创建者
	DisplayName    string    `json:"DisplayName"`    // 页面名称
	ID             *string   `json:"Id"`             // 唯一编号
	Modifier       string    `json:"Modifier"`       // 数据修改者
	ModifyDatetime time.Time `json:"ModifyDatetime"` // 数据修改日期时间
	ParentID       *string   `json:"ParentId"`       // 上级功能或页面编号
	URLAddress     string    `json:"UrlAddress"`     // 页面URL地址
}
