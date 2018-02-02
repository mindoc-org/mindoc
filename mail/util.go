package mail

import (
	"net/mail"
)

func MailAddr(name string, address string) *mail.Address {
	return &mail.Address{
		Name:    name,
		Address: address,
	}
}

type Attachments struct {
	Files   []string
	BaseDir string
}

//SendMail 发送电邮
func SendMail(subject string, content string, receiver, sender string,
	bcc []string, smtpConfig *SMTPConfig, attachments *Attachments) error {
	c := NewSMTPClient(smtpConfig)
	m := NewMail()
	err := m.AddTo(receiver) //receiver e.g. "Barry Gibbs <bg@example.com>"
	if err != nil {
		return err
	}
	err = m.AddFrom(sender)
	if err != nil {
		return err
	}
	m.AddSubject(subject)
	//m.AddText("Some text :)")
	m.AddHTML(content)
	if attachments != nil {
		m.BaseDir = attachments.BaseDir
		for _, v := range attachments.Files {
			err = m.AddAttachment(v)
			if err != nil {
				return err
			}
		}
	}
	for _, addr := range bcc {
		err = m.AddBCC(addr)
		if err != nil {
			return err
		}
	}
	return c.Send(m)
}
