package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSysRole(data []byte) (SysRole, error) {
	var r SysRole
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SysRole) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// 系统角色
type SysRole struct {
	Id             string    `json:"Id"`
	CreateDatetime time.Time `json:"CreateDatetime"`
	Creator        string    `json:"Creator"`
	DisplayName    string    `json:"DisplayName"`
	Modifier       string    `json:"Modifier"`
	ModifyDatetime time.Time `json:"ModifyDatetime"`
}
