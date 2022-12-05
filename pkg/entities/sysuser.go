package entities

import "time"

type SysUser struct {
	Id                string `gorm:"primaryKey"`
	LoginName         string
	DisplayName       string
	Cellphone         string
	Email             string
	Wechat            string
	QQ                string
	IDCardNumber      string `gorm:"column:id_card_number" json:"IDCardNumber"`
	IsDisableLogin    bool   `gorm:"column:is_disable_login"` // 指定 字段名为 is_disable_login
	CreateDatetime    time.Time
	ModifyDatetime    time.Time
	Creator           string
	Modifier          string
	PasswdMD5         string
	SportsCompanyName string `gorm:"column:sports_company_name" json:"SportsCompanyName"`
	SportsCompanyID   string `gorm:"column:sports_company_id" json:"SportsCompanyID"`
}

// 指定表名 为 sys_users
func (SysUser) TableName() string {
	return "sys_users"
}
