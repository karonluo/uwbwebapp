package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

var WebConfiguration *WebConfig

type DBConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfiguration struct {
	Host     string
	Port     string
	Password string
	DBId     int
}

type WebConfig struct {
	Port                string
	SessionExpireMinute int
	Version             string
	PostDataMaxMBSize   int64
	DBConf              DBConfiguration
	RedisConf           RedisConfiguration
	MQTTServerConf      MQTTServerConfiguration
	EmailSmtpServerConf EmailSmtpServerConfiguration
	UrlPathList         []string
}

type MQTTServerConfiguration struct {
	Port           int
	Password       string
	WebSockertPort int
	Broker         string
	User           string
}

type EmailSmtpServerConfiguration struct {
	Host                              string
	Password                          string
	Port                              string
	UserName                          string
	Identity                          string
	ResetSysUserPasswordEmailTemplate string
	ResetSysUserPasswordEmailTimeout  int
}

func LoadWebConfig(confpath string) WebConfig {
	fmt.Print("载入配置文件")
	var conf WebConfig
	jsonData, _ := os.ReadFile(confpath)
	//jsonData, _ = tools.Utf8ToGbk(jsonData)
	fmt.Print("..")

	err := json.Unmarshal([]byte(string(jsonData)), &WebConfiguration)
	if err != nil {
		fmt.Println(err.Error())
	}
	// data := make(map[string]interface{})
	// json.Unmarshal(jsonData, &data)
	fmt.Print("..")

	conf = *WebConfiguration
	fmt.Print("..")
	fmt.Println("成功")
	return conf
}
