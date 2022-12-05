package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

func UnmarshalSportsCompany(data *[]byte) (SportsCompany, error) {
	var r SportsCompany

	err := json.Unmarshal(*data, &r)
	if err != nil {
		fmt.Println("转换错误:\r\n", err.Error())
	}

	return r, err
}

func (r *SportsCompany) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// 体育运动公司
type SportsCompany struct {
	Id             string `gorm:"primaryKey"`
	Address        string
	Creator        string
	Description    string
	Modifier       string
	ModifyDatetime time.Time
	CreateDatetime time.Time
	Name           string
	TelephoneList  string
}

func (SportsCompany) TableName() string {
	return "sports_companies"
}
