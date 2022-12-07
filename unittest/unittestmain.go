package main

import (
	"fmt"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/biz"
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
	rev := []string{`lk139@126.com`, `karonsoft@126.com`, `wengkaiqin@163.com`}
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
	dao.InitRedisDatabase()
	fmt.Println(uuid.New().String())
	var site entities.Site
	fieldSize, err := tools.GetDatabaseTableFieldSize(site, "Users")
	if err == nil {
		fmt.Println(fieldSize)
	} else {
		fmt.Println(err.Error())
	}

	var strTest = "0123456789"
	fmt.Println(strTest[0:10])
	// site, err := dao.GetSiteById("847c2917-a3db-4f47-916e-ff3421eb64e1")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(site.Address)
	// }

	// testRefVariable(&site) // 传输的是对象内存地址而非对象实例，因此内存地址中的值，发生变化时 site 本身的值也会跟着变化，这个和直接传输实例有根本性的不同，实例通过传输形成的时新的实例。
	// fmt.Println(site.Address)

	// testRelVariable(site)
	// fmt.Println(site.Address)
	// TestResetPasswordEmail()
	//TestUpdateSysUser()
	//TestGetSysUserFromDBById()
	// companyIds := []string{"3efc5bfa-e039-d83b-e9fa-2ce3c9605498", "a883fd35-0041-9a64-1f7e-26a7a312b5b4"}

	// er1 := dao.DeleteSportsCompanies(companyIds)
	// if er1 != nil {
	// 	fmt.Println(er1.Error())
	// } else {
	// 	fmt.Println("delete sports_companies is success")
	// }

	// dictvaleus, err := dao.GetSystemDictValues("job_title")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// for _, dv := range dictvaleus {
	// 	fmt.Println(dv.Value)
	// }
	// dictValue, _ := dao.GetSystemDictValue("1001")
	// fmt.Println(dictValue.Value)

	// var condition entities.QueryCondition
	// condition.LikeValue = "98"
	// condition.PageIndex = 1
	// condition.PageSize = 5
	// // cos, err := dao.QuerySportsCompanies(condition)

	// cos, pageCount, err := biz.QueryCompanies(condition)
	// if err != nil {
	// 	tools.ProcessError(`unittest.main`, `cos, err := dao.QuerySportsCompanies(condition)`, err)
	// } else {
	// 	fmt.Printf("总共: %d 页数据\r\n", pageCount)
	// 	btmp, _ := json.Marshal(cos)
	// 	// for _, co := range cos {
	// 	// 	btmp, _ := json.Marshal(co)
	// 	// 	fmt.Println(string(btmp))
	// 	// }
	// 	fmt.Println(string(btmp))
	// }

	// var swimmer entities.Swimmer
	// swimmer.Id = uuid.New().String()
	// swimmer.ModifyDatetime = time.Now()
	// swimmer.CreateDatetime = time.Now()

	// result, _ := json.Marshal(&swimmer)

	// fmt.Println(string(result))

	// funcName := flag.String("func", "ver", "测试功能")
	// flag.Parse()
	// switch *funcName {
	// case "DoTestPublish":
	// 	{
	// 		fmt.Println(*funcName)
	// 		DoTestPublish()
	// 		break
	// 	}

	// }

	// objs, _ := json.Marshal(EnumSiteOwners("847c2917-a3db-4f47-916e-ff3421eb64e1"))
	// fmt.Println(string(objs))

	// var u entities.SysUser
	// u.LoginName = "admin"
	// res := biz.CheckHaveReSysUserByCellphone("18615768209")
	// var message web.WebMessage
	// message.StatusCode = 200
	// message.Message = res
	// smessage, _ := json.Marshal(message)
	// fmt.Println(string(smessage))

	// mqtt.InitMQTTClient()

	// sysuser, success := dao.GetUserFromRedisByLoginName("yangqin")
	// if success {
	// 	fmt.Println(sysuser.DisplayName)
	// } else {
	// 	fmt.Println("未找到该用户!")
	// }
	// biz.ClearSiteOwners("847c2917-a3db-4f47-916e-ff3421eb64e1")
	// biz.SetSiteOwners("847c2917-a3db-4f47-916e-ff3421eb64e1", "ecb22d21-e346-44f7-bfd1-4b58f66ed226,dec9fd43-f95b-4888-95d5-cbe9f30c2d58")
	// // TestMQTTPulishMessage()
	// //result, _ := json.Marshal(EnumSiteOwners("847c2917-a3db-4f47-916e-ff3421eb64e1"))
	// //fmt.Println(string(result))
	// //CreateUser()

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
