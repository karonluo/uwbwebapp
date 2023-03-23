package entities

import (
	"encoding/json"
	"time"
)

func UnmarshalSwimmer(data []byte) (Swimmer, error) {
	var r Swimmer
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Swimmer) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// 游泳者(会员)
type Swimmer struct {
	ID             string `gorm:"primaryKey"`
	DisplayName    string
	Gender         string
	IDCardNumber   string `gorm:"column:id_card_number" json:"IDCardNumber"`
	Age            int
	Address        string
	Cellphone      string
	Wechat         string
	Description    string
	CreateDatetime time.Time
	ModifyDatetime time.Time
	Creator        string
	Modifier       string
}
