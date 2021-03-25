package conf

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
)

type SmtpConf struct {
	EnableMail   bool
	MailNumber   int
	SmtpUserName string
	SmtpHost     string
	SmtpPassword string
	SmtpPort     int
	FormUserName string
	MailExpired  int
	Secure       string
}

func GetMailConfig() *SmtpConf {
	user_name, _ := web.AppConfig.String("smtp_user_name")
	password, _ := web.AppConfig.String("smtp_password")
	smtp_host, _ := web.AppConfig.String("smtp_host")
	smtp_port := web.AppConfig.DefaultInt("smtp_port", 25)
	form_user_name, _ := web.AppConfig.String("form_user_name")
	enable_mail, _ := web.AppConfig.String("enable_mail")
	mail_number := web.AppConfig.DefaultInt("mail_number", 5)
	secure := web.AppConfig.DefaultString("secure", "NONE")

	if secure != "NONE" && secure != "LOGIN" && secure != "SSL" {
		secure = "NONE"
	}
	c := &SmtpConf{
		EnableMail:   strings.EqualFold(enable_mail, "true"),
		MailNumber:   mail_number,
		SmtpUserName: user_name,
		SmtpHost:     smtp_host,
		SmtpPassword: password,
		FormUserName: form_user_name,
		SmtpPort:     smtp_port,
		Secure:       secure,
	}
	return c
}
