package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/entities"
)

func UWBDevicePlatformLoginInformationFromCache() (entities.UWBDevicePlatformLoginInformation, error) {

	// fmt.Println("调用从 Cache 中获取 UWB 设备平台登录信息")
	var uwbLoginInfo entities.UWBDevicePlatformLoginInformation
	rctx := context.Background()
	result, err := cache.RedisDatabase.Get(rctx, "UWBPlatLoginToken").Result()
	if err == nil {
		err = json.Unmarshal([]byte(result), &uwbLoginInfo)
	}
	return uwbLoginInfo, err
}

// UWB 设备平台登录
// fouceRefresh 强制刷新登录信息
func UWBDevicePlatformLogin(fouceRefresh bool) (entities.UWBDevicePlatformLoginInformation, error) {
	// fmt.Println("调用 UWB 设备平台登录功能")
	var uwbLoginInfo entities.UWBDevicePlatformLoginInformation
	var err error
	if !fouceRefresh {
		uwbLoginInfo, err = UWBDevicePlatformLoginInformationFromCache()
		uwbLoginInfo.Data.AccessToken = fmt.Sprintf("Bearer %s", uwbLoginInfo.Data.AccessToken)
	}
	if err != nil || fouceRefresh {
		// fmt.Println("调用 UWB 设备平台登录接口")
		uwbDevicePlatformConf := conf.WebConfiguration.UWBDevicePlatformConf
		loginInterface := fmt.Sprintf("%s%s", uwbDevicePlatformConf.Address, uwbDevicePlatformConf.LoginInterface)
		client := &http.Client{}
		data := make(map[string]interface{})
		data["name"] = uwbDevicePlatformConf.User
		data["password"] = uwbDevicePlatformConf.Password
		bytesData, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", loginInterface, bytes.NewReader(bytesData))
		req.Header.Add("Content-Type", "application/json")
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &uwbLoginInfo)
		uwbLoginInfo.Data.AccessToken = fmt.Sprintf("Bearer %s", uwbLoginInfo.Data.AccessToken)
		fmt.Printf("登录成功, Token: %s", uwbLoginInfo.Data.AccessToken)
		// 将登录信息写入 cache(redis)
		rctx := context.Background()
		cache.RedisDatabase.Set(rctx, "UWBPlatLoginToken", string(body), 5*time.Minute) // 写入认证信息到 cache, 方便下次使用，超时时间5分钟。
	}

	return uwbLoginInfo, err

}

func UWBDevicePlatformDeleteTerminal() error {
	var err error
	return err
}

// 修改 UWB 平台终端信息
func UWBDevicePlatformUpdateTerminal(swimmerDisplayName string, swimmerGender string, swimmerId string, serialNumber string, modelId int64) (string, error) {
	type Terminal struct {
		ApplicationId int64                  `json:"applicationId"`
		ModelId       int64                  `json:"modelId"`
		Name          string                 `json:"name"`
		Properties    map[string]interface{} `json:"properties"`
		SerialNumber  string                 `json:"serialNumber"`
	}
	var result string
	var ter Terminal
	ter.Properties = make(map[string]interface{})
	ter.Properties["swimmerDisplayName"] = swimmerDisplayName
	ter.Properties["swimmerGender"] = swimmerGender
	ter.Properties["swimmerId"] = swimmerId
	ter.ApplicationId = conf.WebConfiguration.UWBDevicePlatformConf.ApplicationId
	ter.ModelId = modelId
	ter.Name = swimmerDisplayName
	ter.SerialNumber = serialNumber
	var bytesData []byte
	var err error
	var loginInfo entities.UWBDevicePlatformLoginInformation
	loginInfo, err = UWBDevicePlatformLogin(false)
	if err == nil {
		uwbDevicePlatformConf := conf.WebConfiguration.UWBDevicePlatformConf
		updateInterface := fmt.Sprintf("%s%s%s", uwbDevicePlatformConf.Address, uwbDevicePlatformConf.UpdateTerminalInterface, serialNumber)
		client := &http.Client{}
		updateInterface = fmt.Sprintf("%s?serialNumber=%s", updateInterface, serialNumber)
		bytesData, err = json.Marshal(ter)
		fmt.Println(string(bytesData))
		req, _ := http.NewRequest("PUT", updateInterface, bytes.NewReader(bytesData))
		fmt.Println(updateInterface)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", loginInfo.Data.AccessToken)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)
	}

	return result, err
}

// 创建 UWB 平台终端信息
func UWBDevicePlatformCreateTerminal(swimmerDisplayName string, swimmerGender string, swimmerId string, serialNumber string, modelId int64) (string, error) {
	type Terminal struct {
		ApplicationId int64                  `json:"applicationId"`
		ModelId       int64                  `json:"modelId"`
		Name          string                 `json:"name"`
		Properties    map[string]interface{} `json:"properties"`
		SerialNumber  string                 `json:"serialNumber"`
	}
	var result string
	var ter Terminal
	ter.Properties = make(map[string]interface{})
	ter.Properties["swimmerDisplayName"] = swimmerDisplayName
	ter.Properties["swimmerGender"] = swimmerGender
	ter.Properties["swimmerId"] = swimmerId
	ter.ApplicationId = conf.WebConfiguration.UWBDevicePlatformConf.ApplicationId
	ter.ModelId = modelId
	ter.Name = swimmerDisplayName
	ter.SerialNumber = serialNumber
	var bytesData []byte
	var err error
	var loginInfo entities.UWBDevicePlatformLoginInformation
	loginInfo, err = UWBDevicePlatformLogin(false)
	if err == nil {
		uwbDevicePlatformConf := conf.WebConfiguration.UWBDevicePlatformConf
		updateInterface := fmt.Sprintf("%s%s", uwbDevicePlatformConf.Address, uwbDevicePlatformConf.CreateTerminalInterface)
		client := &http.Client{}
		// updateInterface = fmt.Sprintf("%s?serialNumber=%s", updateInterface, serialNumber)
		bytesData, err = json.Marshal(ter)
		fmt.Println(string(bytesData))
		req, _ := http.NewRequest("POST", updateInterface, bytes.NewReader(bytesData))
		fmt.Println(updateInterface)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", loginInfo.Data.AccessToken)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)
	}

	return result, err
}

// 获取 UWB 平台终端信息
func UWBDevicePlatformGetTerminal(serialNumber string) (string, error) {
	var err error
	var result string
	var loginInfo entities.UWBDevicePlatformLoginInformation
	loginInfo, err = UWBDevicePlatformLogin(false)
	if err == nil {
		uwbDevicePlatformConf := conf.WebConfiguration.UWBDevicePlatformConf
		getInterface := fmt.Sprintf("%s%s%s", uwbDevicePlatformConf.Address, uwbDevicePlatformConf.GetTerminalInterface, serialNumber)
		client := &http.Client{}
		data := make(map[string]interface{})
		data["name"] = uwbDevicePlatformConf.User
		data["password"] = uwbDevicePlatformConf.Password
		bytesData, _ := json.Marshal(data)
		getInterface = fmt.Sprintf("%s?serialNumber=%s", getInterface, serialNumber)
		req, _ := http.NewRequest("GET", getInterface, bytes.NewReader(bytesData))

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", loginInfo.Data.AccessToken)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)

	} else {
		fmt.Println(err.Error())
	}
	return result, err
}

// 搜索UWB平台终端信息
func UWBDevicePlatformQueryTerminal(queryCondition entities.QueryCondition) (string, error) {
	var err error
	var result string
	var loginInfo entities.UWBDevicePlatformLoginInformation
	loginInfo, err = UWBDevicePlatformLogin(false)
	if err == nil {
		uwbDevicePlatformConf := conf.WebConfiguration.UWBDevicePlatformConf
		searchInterface := fmt.Sprintf("%s%s", uwbDevicePlatformConf.Address, uwbDevicePlatformConf.SearchTerminalInterface)
		applicationId := fmt.Sprintf("%d", uwbDevicePlatformConf.ApplicationId)
		organizationId := uwbDevicePlatformConf.OrganizationId
		page := queryCondition.PageIndex
		pageSize := queryCondition.PageSize
		searchKey := queryCondition.LikeValue

		client := &http.Client{}
		data := make(map[string]interface{})
		data["name"] = uwbDevicePlatformConf.User
		data["password"] = uwbDevicePlatformConf.Password
		bytesData, _ := json.Marshal(data)
		searchInterface = fmt.Sprintf(searchInterface+"?page=%d&pageSize=%d&applicationId=%d&organizationId=%d&searchKey=%s",
			page, pageSize, applicationId, organizationId, searchKey)
		req, _ := http.NewRequest("GET", searchInterface, bytes.NewReader(bytesData))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", loginInfo.Data.AccessToken)
		resp, _ := client.Do(req)
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)
	} else {
		fmt.Println(err.Error())
	}
	return result, err
}
