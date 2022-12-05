// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    sysDict, err := UnmarshalSysDict(bytes)
//    bytes, err = sysDict.Marshal()

package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSysDict(data []byte) (SysDict, error) {
	var r SysDict
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SysDict) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// 系统字典
type SysDict struct {
	Code           string    `json:"Code"`                             // 编码
	CreateDatetime time.Time `json:"CreateDatetime"`                   // 数据创建日期时间
	Creator        string    `json:"Creator"`                          // 数据创建人
	DataType       string    `json:"DataType"`                         // 数据类型
	Key            string    `json:"Key"`                              // 建
	Modifier       string    `json:"Modifier"`                         // 数据修改人
	ModifyDatetime time.Time `json:"ModifyDatetime"`                   // 数据修改日期时间
	Value          string    `json:"Value"`                            // 值
	OrderKey       int       `gorm:"column:order_key" json:"OrderKey"` // 排序键
}
