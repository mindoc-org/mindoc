package conf

import (
	"github.com/astaxie/beego"
	"strings"
)

type SmtpConf struct {
	EnableMail bool
	MailNumber int
	SmtpUserName string
	SmtpHost string
	SmtpPassword string
	SmtpPort int
	FormUserName string
	MailExpired int
}

func GetMailConfig() *SmtpConf {
	user_name := beego.AppConfig.String("smtp_user_name")
	password := beego.AppConfig.String("smtp_password")
	smtp_host := beego.AppConfig.String("smtp_host")
	smtp_port := beego.AppConfig.DefaultInt("smtp_port",25)
	form_user_name := beego.AppConfig.String("form_user_name")
	enable_mail := beego.AppConfig.String("enable_mail")
	mail_number := beego.AppConfig.DefaultInt("mail_number",5)

	c := &SmtpConf{
		EnableMail : strings.EqualFold(enable_mail,"true"),
		MailNumber: mail_number,
		SmtpUserName:user_name,
		SmtpHost:smtp_host,
		SmtpPassword:password,
		FormUserName:form_user_name,
		SmtpPort:smtp_port,
	}
	return c
}