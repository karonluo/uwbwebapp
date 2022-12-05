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
	Id                 string `gorm:"primaryKey"`
	UWBTagCode         string `gorm:"uwb_tag_code" json:"UWBTagCode"`
	DisplayName        string
	Gender             string
	IDCardNumber       string `gorm:"column:id_card_number" json:"IDCardNumber"`
	Age                int
	MemberCode         string
	VIPLevel           string `gorm:"column:vip_level"`
	VIPLevelDictCode   string `gorm:"column:vip_level_dict_code"`
	Address            string
	Cellphone          string
	Wechat             string
	Description        string
	CreateDatetime     time.Time
	ModifyDatetime     time.Time
	Creator            string
	Modifier           string
	MembershipDatetime time.Time
	WithdrawaDatetime  time.Time
}
