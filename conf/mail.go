package conf

import (
	"strings"

	"github.com/beego/beego/v2/adapter"
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
	user_name := adapter.AppConfig.String("smtp_user_name")
	password := adapter.AppConfig.String("smtp_password")
	smtp_host := adapter.AppConfig.String("smtp_host")
	smtp_port := adapter.AppConfig.DefaultInt("smtp_port", 25)
	form_user_name := adapter.AppConfig.String("form_user_name")
	enable_mail := adapter.AppConfig.String("enable_mail")
	mail_number := adapter.AppConfig.DefaultInt("mail_number", 5)
	secure := adapter.AppConfig.DefaultString("secure", "NONE")

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
