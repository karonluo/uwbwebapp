package message

import (
	"fmt"
	"net/smtp"
	"uwbwebapp/conf"
)

// 发送邮件
// receivers 接收人列表
// htmlBody HTML 内容
// subject 主题
func SendEmail(receivers []string, htmlBody string, subject string) error {
	emailConf := conf.WebConfiguration.EmailSmtpServerConf
	fmt.Println(emailConf.Identity, emailConf.UserName, emailConf.Password, emailConf.Host, emailConf.Port)
	auth := smtp.PlainAuth(emailConf.Identity, emailConf.UserName, emailConf.Password, emailConf.Host)
	strMsg := fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, htmlBody)
	fmt.Println(strMsg)
	msg := []byte(strMsg)
	err := smtp.SendMail(emailConf.Host+emailConf.Port, auth, emailConf.UserName, receivers, msg)
	return err
}
