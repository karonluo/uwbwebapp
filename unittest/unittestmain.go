package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/asynctimer"
	"uwbwebapp/pkg/biz"
	"uwbwebapp/pkg/cache"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/message"

	"uwbwebapp/pkg/tools"

	"github.com/google/uuid"
)

func TestGetSysUserFromDBById() {
	result, err := dao.GetSysUserFromDBById("ecb22d21-e346-44f7-bfd1-4b58f66ed226")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}
func TestUpdateSysUser() {
	var user entities.SysUser
	user.QQ = "12345"
	user.Cellphone = "13608887886"
	user.Id = "ecb22d21-e346-44f7-bfd1-4b58f66ed226"
	user.DisplayName = "我爱杨琴"
	user.LoginName = "yangqin"
	user.Modifier = "admin"
	user.IsDisableLogin = false
	user.Email = "aiyangqin@gmail.com"
	user.Wechat = "12345"
	err := biz.UpdateSysUser(user)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Update success!")
	}
}
func TestSendEmail() {
	rev := []string{`lk139@126.com`}
	message.SendEmail(rev, `<div>点击链接 <a href='https://www.baidu.com?token=u123456uid' target='_blank'>重置密码</a></div>`, `Karonsoft 注册账号 [admin] 重置密码`)
}
func TestResetPasswordEmail() {
	biz.SendForgetPasswordEmail("lk139@126.com")
}

// func testRefVariable(site *entities.Site) {

//		fmt.Println(site.Address)
//		site.Address = "上海市"
//	}
//
//	func testRelVariable(site entities.Site) {
//		fmt.Println(site.Address)
//		site.Address = "成都市"
//	}
func main() {
	conf.LoadWebConfig("../conf/WebConfig.json")
	dao.InitDatabase()
	cache.InitRedisDatabase()
	// nodate := "2001-01-01 00:00:00"
	// tnodate, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, nodate)
	// fmt.Println(tnodate)
	// fmt.Println(tools.CheckNoDate(tnodate))
	cbt, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "2021-01-01 01:30:35")
	cet, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "2021-01-01 05:30:35")

	tbt, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "2021-01-01 01:25:00")
	tet, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "2021-01-01 01:29:00")

	fmt.Println(tet)
	result, _ := tools.CheckTimesHasOverlap(cbt, cet, tbt, tet)
	fmt.Println(result)
	// data, err := biz.UWBDevicePlatformLogin()
	// if err == nil {
	// 	// fmt.Println(data["data"])

	// 	fmt.Println(data.Data.AccessToken)
	// }
	var query entities.QueryCondition
	query.LikeValue = ""
	query.PageIndex = 1
	query.PageSize = 10

	res, er := biz.UWBDevicePlatformQueryTerminal(query)
	if er == nil {
		// fmt.Println(res)
		var searchResult entities.UWBDevicePlatformTerminalSearchResult
		json.Unmarshal([]byte(res), &searchResult)
		for _, v := range searchResult.Data.Data {
			fmt.Printf("UWB TERMINAL INFORMATION %s\r\nSerial Number: %s, Property-Test: %s\r\n", v.Name, v.SerialNumber, v.Properties.Test)
		}
	} else {
		fmt.Println(er.Error())
	}

	// 获取场地游泳者统计报表
	DoEnumSiteSimmerForReport()

	// // 获取 UWB 终端标签信息
	DoGetTerminal()
	// // 更新 UWB 终端标签信息
	// DoUpdateTerminal()
	// // 获取 UWB 终端标签信息
	// DoGetTerminal()

	// 做时间相关测试
	//DoTimeTest()
	// 测试危险告警
	// DoDangerAlert()

	// DoCacheTest()
	DoQuerySwimmerTest()

	cos, err := dao.EnumSportsCompaniesByRightUser("d2f3de8b-b091-434f-b66f-62c9ec0eab48")
	if err == nil {
		text, _ := json.Marshal(cos)
		fmt.Println(string(text))
	}

	codes, ere := dao.EnumSiteFences("543b8172-5bcc-aeaa-f5a4-a637608fb737")
	if ere == nil {
		for _, code := range codes {
			fmt.Println(code)
		}
	}
}
func DoQuerySwimmerTest() {
	var query entities.QueryCondition
	query.LikeValue = ""
	query.PageIndex = 1
	query.PageSize = 10
	var company_id string = `5a606630-d515-4724-8579-b12af464ddeb`
	res, pageCount, recordCount, err := biz.QuerySwimmers(query, company_id)
	fmt.Printf("总页数: %d\r\n", pageCount)
	fmt.Printf("总记录数: %d\r\n", recordCount)
	if err == nil {
		for _, swimmer := range res {
			fmt.Println(swimmer.DisplayName)
		}
	} else {
		fmt.Println(err.Error())
	}
}
func DoCacheTest() {
	// 测试当 key 为空时返回什么内容。
	rctx := context.Background()
	result, err := cache.RedisDatabase.Get(rctx, "hello").Result()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("result: '%s'", result)
}

func DoDangerAlert() {
	go asynctimer.AlertDangerInformation()
	for {
		time.Sleep(500 * time.Millisecond)
	}
}
func DoGetTerminal() {
	res, er := biz.UWBDevicePlatformGetTerminal("00002AC111000002")
	if er == nil {
		fmt.Println(res)
		var o entities.UWBDevicePlatformTerminalGetResult
		er = json.Unmarshal([]byte(res), &o)
		if er != nil {
			fmt.Println(er.Error())
		}
	} else {
		fmt.Println(er.Error())
	}
}

func DoUpdateTerminal() {

	_, er := biz.UWBDevicePlatformUpdateTerminal("钱芳", "女", "6B03FaD8-c259-bd6E-5d0c-CafAB79DE9c8", "00002AC111000002",
		conf.WebConfiguration.UWBDevicePlatformConf.DefaultTerminalModelId)

	if er != nil {
		fmt.Println(er.Error())
	}
}
func DoTimeTest() {
	t, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, "2023-03-13 17:55:05")
	st1 := time.Now().Format(tools.GOOGLE_DATETIME_FORMAT_NO_NANO)
	now, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_NANO, st1)
	tsub := now.Sub(t)
	fmt.Printf("当前时间与指定时间相差：%f 秒\r\n", tsub.Seconds())

	// fmt.Println(time.Now().Format(tools.GOOGLE_DATETIME_FORMAT_NO_DATE_NO_NANO))
	// t, _ := time.Parse(tools.GOOGLE_DATETIME_FORMAT_NO_DATE_NO_NANO, "23:59:59")
	// fmt.Println(t.Format(tools.GOOGLE_DATETIME_FORMAT))
	// fmt.Println(t.Hour() - time.Now().Hour())
	// go asynctimer.SumSwimmerDistance()
	// for {
	// 	time.Sleep(500 * time.Millisecond)
	// }
}
func DoEnumSiteSimmerForReport() {
	// result, err := dao.EnumSiteSimmerForReport("29972eb5-6419-48d2-94e1-e3bc5bb99188")
	// if err == nil {
	// 	for _, swimmer := range result {
	// 		fmt.Printf("Swimmer: %s, age:%d, gender:%s\r\n", swimmer.DisplayName, swimmer.Age, swimmer.Gender)
	// 	}
	// } else {
	// 	fmt.Println(err.Error())
	// }
	report, _ := biz.SiteSwimmerReport("29972eb5-6419-48d2-94e1-e3bc5bb99188")
	bytes, _ := json.Marshal(report)
	fmt.Println(string(bytes))
}
func DoGenerateSysFuncPagesToDB() {
	type TMP struct {
		Url         string `json:"url"`
		Method      string `json:"method"`
		DisplayName string `json:"display_name"`
	}
	var tmp []TMP
	var bcontent []byte
	bcontent, _ = os.ReadFile("../sysfuncpages.json")
	fmt.Println(string(bcontent))
	json.Unmarshal(bcontent, &tmp)
	var pageId string
	var err error
	for idx, t := range tmp {
		// fmt.Println(t.Url)
		var funcPage entities.SysFuncPage
		funcPage.DisplayName = t.DisplayName
		funcPage.OrderKey = idx
		funcPage.URLAddress = t.Url
		funcPage.ParentID = "top"
		funcPage.URLType = "INTERFACE"
		funcPage.URLMethod = strings.ToUpper(t.Method)
		pageId, err = biz.CreateSysFuncPage(&funcPage)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(pageId)
		}

	}
}
func TestMQTTPulishMessage() {
	var i int
	for {
		i = i + 1
		text := fmt.Sprintf("Message %d", i)
		token := message.MQTTClient.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
		fmt.Println(text)
	}

}
func EnumSiteOwners(siteId string) []map[string]interface{} {
	result := dao.EnumSiteOwners(siteId)
	return result
}

func CreateUser() {
	var sysuser entities.SysUser
	sysuser.Cellphone = "18900000000"
	sysuser.CreateDatetime = time.Now()
	sysuser.Creator = "admin"
	sysuser.DisplayName = "Jasmine"
	sysuser.Email = "Jasmine@126.com"
	sysuser.Id = uuid.New().String()
	sysuser.IsDisableLogin = false
	sysuser.LoginName = "jasmine"
	sysuser.Modifier = "admin"
	sysuser.ModifyDatetime = time.Now()
	sysuser.PasswdMD5 = tools.SHA1("Password123")
	sysuser.QQ = "12345678"
	sysuser.Wechat = "12345678"

	dao.CreateSysUser(sysuser)
}
