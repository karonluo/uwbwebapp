package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"uwbwebapp/conf"
	"uwbwebapp/pkg/dao"
	"uwbwebapp/pkg/entities"
	"uwbwebapp/pkg/message"
	"uwbwebapp/pkg/tools"

	"github.com/google/uuid"
)

// TODO: ACL FROM DATABASE OR REDIS
func GetUserACLByAuthorization() []string {
	/*
		acl := []string{
			"/swimmer/query",
			"/swimmer",
			"/login",
			"/sysuser",
			"/site",
			"/site/query",
			"/siteowners/list",
			"/siteowners/setowners",
			"/siteowners/list",
			"/sysuser/userinfo",
			"/athorizationinfo",
			"/company",
			"/sysuser/listall",
			"/company/relsites",
			"/role", "/sysfuncpage/list",
			"/company/query",
			"/companies",
			"/sysusers",
			"/sysuser/query",
			"/dict", "/dict/children"}
	*/

	return conf.WebConfiguration.UrlPathList
}
func LoginSystem(loginName string, password string) (string, bool) {
	var success bool = true
	var msg string

	sysuser, res := dao.GetUserFromRedisByLoginName(loginName)
	if !res {
		// 尝试从数据库中获取
		sysuser, res, _ = dao.GetUserFromDBByLoginName(loginName)
		if res {
			// 设置到缓存中
			dao.SetUserToRedis(sysuser)
		}
	}
	if res {

		if sysuser.PasswdMD5 == tools.SHA1(password) {
			rctx := context.Background()
			msg = tools.SHA1(uuid.New().String())
			var token_val map[string]interface{} = make(map[string]interface{})
			token_val["token"] = msg
			token_val["login_name"] = loginName
			err1 := dao.RedisDatabase.HSet(rctx, "token_"+msg, token_val).Err()
			if err1 != nil {
				tools.ProcessError(`biz.LoginSystem`, `dao.RedisDatabase.HSet(rctx, "token_"+msg, token_val)`, err1, "pkg/biz/sysuser.go")
				msg = err1.Error()
				success = false
			}
			// TODO: 放置 ACL (该用户可访问的URL和功能模块)
			acl := GetUserACLByAuthorization()
			res, _ := json.Marshal(acl)
			err2 := dao.RedisDatabase.HSet(rctx, "token_"+msg, "acl", res).Err()
			if err2 != nil {
				tools.ProcessError(`biz.LoginSystem`, `dao.RedisDatabase.HSet(rctx, "token_"+msg, "acl", res)`, err2, "pkg/biz/sysuser.go")
				msg = err2.Error()
				success = false
			}
			dao.RedisDatabase.Expire(rctx, "token_"+msg, time.Duration(conf.WebConfiguration.SessionExpireMinute)*time.Minute) // 设置20分钟后过期
		} else {
			// 密码不正确
			success = false
		}
	} else {
		success = false
	}
	return msg, success
}

func EnumSysUserFromDB() ([]entities.SysUser, int64) {
	return dao.EnumSysUserFromDB()
}

func CreateSysUser(user entities.SysUser) (string, error) {
	var err error
	if user.Id == "" {
		user.Id = uuid.New().String()
	}
	user.CreateDatetime = time.Now()
	user.ModifyDatetime = time.Now()
	if user.Creator != "" {
		user.Creator = "admin"
	}
	user.Modifier = user.Creator
	user.PasswdMD5 = tools.SHA1(user.PasswdMD5)
	if dao.CheckHaveReSysUser(user) {
		err = errors.New("有相同的系统用户信息, 请注意: 电子邮箱地址、手机号、登录名、身份证号 不能重复且必填。")
		user.Id = ""
	} else {
		dao.CreateSysUser(user)
	}
	return user.Id, err
}

func GetSysUserFromDBByLoginName(login_name string) entities.SysUser {
	result, _, _ := dao.GetUserFromDBByLoginName(login_name)
	return result
}

func CheckHaveReSysUserByLoginName(login_name string) bool {
	var user entities.SysUser
	user.LoginName = login_name
	return dao.CheckHaveReSysUser(user)
}

func CheckHaveReSysUserByEmail(email string) bool {
	var user entities.SysUser
	user.Email = email
	return dao.CheckHaveReSysUser(user)
}

func CheckHaveReSysUserByCellphone(cellphone string) bool {
	var user entities.SysUser
	user.Cellphone = cellphone
	return dao.CheckHaveReSysUser(user)
}

func DeleteSysUser(id string) (bool, error) {
	var user entities.SysUser
	user.Id = id
	return dao.DeleteSysUser(user)
}

func DeleteSysUsers(ids []string) error {
	return dao.DeleteSysUsers(ids)
}
func QuerySysUsers(queryCondition entities.QueryCondition, companyID string) ([]entities.SysUser, int64, int64, error) {
	var users []entities.SysUser
	dataRecordCount, err := dao.GetSysUserCount(queryCondition, companyID)
	pageCount := float64(dataRecordCount) / float64(queryCondition.PageSize)
	if err == nil {
		users, err = dao.QuerySysUsers(queryCondition, companyID)
	}
	return users, int64(math.Ceil(pageCount)), dataRecordCount, err
}

func UpdateSysUserPassword(sysUserId string, originPassword string, newPassword string) error {
	var tmpUser entities.SysUser
	var err error
	tmpUser, err = dao.GetSysUserFromDBById(sysUserId)
	if err == nil {
		if tmpUser.PasswdMD5 == tools.SHA1(originPassword) {
			err = dao.UpdateSysUserPassword(sysUserId, tools.SHA1(newPassword))
		}
	}
	return err
}

func UpdateSysUser(user entities.SysUser) error {
	var tmpUser entities.SysUser
	var err error
	tmpUser, err = dao.GetSysUserFromDBById(user.Id)
	if err == nil {
		// 防止以下字段被修改
		user.CreateDatetime = tmpUser.CreateDatetime
		user.Creator = tmpUser.Creator
		user.ModifyDatetime = time.Now()
		user.PasswdMD5 = tmpUser.PasswdMD5
		if user.Modifier == "" {
			user.Modifier = "admin"
		}
		err = dao.UpdateSysUser(user)
	}

	return err
}

func ResetSysUserPassword(password string, token string) error {
	var err error
	fmt.Println(password, token)
	//TODO: ResetSysUserPassword
	return err
}

func SendForgetPasswordEmail(email string) error {
	emailConf := conf.WebConfiguration.EmailSmtpServerConf
	var err error
	if CheckHaveReSysUserByEmail(email) {

		// 当检测到有该邮件地址时发送重置密码的邮件。
		token := tools.SHA1(uuid.New().String())
		rctx := context.Background()
		//将 token 存入 redis 用于重置密码的key中, 并根据配置进行超时设置。
		val := fmt.Sprintf(`{"email":"%s", "token":"%s"}`, email, token)
		dao.RedisDatabase.Set(rctx, "resetpwd_"+token, val, time.Minute*time.Duration(emailConf.ResetSysUserPasswordEmailTimeout))
		fmt.Println(emailConf.ResetSysUserPasswordEmailTemplate)
		html := strings.ReplaceAll(emailConf.ResetSysUserPasswordEmailTemplate, "{{token}}", token)
		html = strings.ReplaceAll(html, "{{timeout}}", fmt.Sprintf("%d", emailConf.ResetSysUserPasswordEmailTimeout))
		fmt.Println(html)
		err = message.SendEmail([]string{email}, html, "密码重置邮件")
	}
	return err
	// "ResetSysUserPasswordEmailTemplate1": "Click <a href='http://172.0.0.1/resetpwd?token={{token}}' target='_blank'> to reset password. </a>or enter reset password page and input code: {{token}} <p>, The code timeout is  {{timeout}} minute. <br> thanks!",
}

func EnumSysUsersFromSportsCompanyIds(siteIds []string) ([]entities.SysUser, error) {
	users, err := dao.EnumSysUsersFromSportsCompanyIds(siteIds)
	return users, err
}
