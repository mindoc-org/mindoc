package mail

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
)

var (
	imageRegex  = regexp.MustCompile(`(src|background)=["'](.*?)["']`)
	schemeRegxp = regexp.MustCompile(`^[a-zA-Z]+://`)
)

// Mail will represent a formatted email
type Mail struct {
	To         []string
	ToName     []string
	Subject    string
	HTML       string
	Text       string
	From       string
	Bcc        []string
	FromName   string
	ReplyTo    string
	Date       string
	Files      map[string]string
	Headers    string
	BaseDir    string //内容中图片路径
	Charset    string //编码
	RetReceipt string //回执地址，空白则禁用回执
}

// NewMail returns a new Mail
func NewMail() Mail {
	return Mail{}
}

// SMTPClient struct
type SMTPClient struct {
	smtpAuth smtp.Auth
	host     string
	port     string
	user     string
	secure   string
}

// SMTPConfig 配置结构体
type SMTPConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Secure   string
	Identity string
}

func (s *SMTPConfig) Address() string {
	if s.Port == 0 {
		s.Port = 25
	}
	return s.Host + `:` + strconv.Itoa(s.Port)
}

func (s *SMTPConfig) Auth() smtp.Auth {
	var auth smtp.Auth
	s.Secure = strings.ToUpper(s.Secure)
	switch s.Secure {
	case "NONE":
		auth = unencryptedAuth{smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host)}
	case "LOGIN":
		auth = LoginAuth(s.Username, s.Password)
	case "SSL":
		fallthrough
	default:
		//auth = smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host)
		auth = unencryptedAuth{smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host)}
	}
	return auth
}

func NewSMTPClient(conf *SMTPConfig) SMTPClient {
	return SMTPClient{
		smtpAuth: conf.Auth(),
		host:     conf.Host,
		port:     strconv.Itoa(conf.Port),
		user:     conf.Username,
		secure:   conf.Secure,
	}
}

// NewMail returns a new Mail
func (c *SMTPClient) NewMail() Mail {
	return NewMail()
}

// Send - It can be used for generic SMTP stuff
func (c *SMTPClient) Send(m Mail) error {
	length := 0
	if len(m.Charset) == 0 {
		m.Charset = "utf-8"
	}
	boundary := "COSCMSBOUNDARYFORSMTPGOLIB"
	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("X-SMTPAPI: %s\r\n", m.Headers))
	//回执
	if len(m.RetReceipt) > 0 {
		message.WriteString(fmt.Sprintf("Return-Receipt-To: %s\r\n", m.RetReceipt))
		message.WriteString(fmt.Sprintf("Disposition-Notification-To: %s\r\n", m.RetReceipt))
	}
	message.WriteString(fmt.Sprintf("From: %s <%s>\r\n", m.FromName, m.From))
	if len(m.ReplyTo) > 0 {
		message.WriteString(fmt.Sprintf("Return-Path: %s\r\n", m.ReplyTo))
	}
	length = len(m.To)
	if length > 0 {
		nameLength := len(m.ToName)
		if nameLength > 0 {
			message.WriteString(fmt.Sprintf("To: %s <%s>", m.ToName[0], m.To[0]))
		} else {
			message.WriteString(fmt.Sprintf("To: <%s>", m.To[0]))
		}
		for i := 1; i < length; i++ {
			if nameLength > i {
				message.WriteString(fmt.Sprintf(", %s <%s>", m.ToName[i], m.To[i]))
			} else {
				message.WriteString(fmt.Sprintf(", <%s>", m.To[i]))
			}
		}
	}
	length = len(m.Bcc)
	if length > 0 {
		message.WriteString(fmt.Sprintf("Bcc: <%s>", m.Bcc[0]))
		for i := 1; i < length; i++ {
			message.WriteString(fmt.Sprintf(", <%s>", m.Bcc[i]))
		}
	}
	message.WriteString("\r\n")
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	if m.Files != nil {
		message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\n--%s\r\n", boundary, boundary))
	}
	if len(m.HTML) > 0 {
		//解析内容中的图片
		rs := imageRegex.FindAllStringSubmatch(m.HTML, -1)
		var embedImages string
		for _, v := range rs {
			surl := v[2]
			if v2 := schemeRegxp.FindStringIndex(surl); v2 == nil {
				filename := path.Base(surl)
				directory := path.Dir(surl)
				if directory == "." {
					directory = ""
				}
				h := md5.New()
				h.Write([]byte(surl + "@coscms.0"))
				cid := hex.EncodeToString(h.Sum(nil))
				if len(m.BaseDir) > 0 && !strings.HasSuffix(m.BaseDir, "/") {
					m.BaseDir += "/"
				}
				if len(directory) > 0 && !strings.HasSuffix(directory, "/") {
					directory += "/"
				}
				if str, err := m.ReadAttachment(m.BaseDir + directory + filename); err == nil {
					re3 := regexp.MustCompile(v[1] + `=["']` + regexp.QuoteMeta(surl) + `["']`)
					m.HTML = re3.ReplaceAllString(m.HTML, v[1]+`="cid:`+cid+`"`)

					embedImages += fmt.Sprintf("--%s\r\n", boundary)
					embedImages += fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"; charset=\"%s\"\r\n", filename, m.Charset)
					embedImages += fmt.Sprintf("Content-Description: %s\r\n", filename)
					embedImages += fmt.Sprintf("Content-Disposition: inline; filename=\"%s\"; charset=\"%s\"\r\n", filename, m.Charset)
					embedImages += fmt.Sprintf("Content-Transfer-Encoding: base64\r\nContent-ID: <%s>\r\n\r\n%s\r\n\n", cid, str)
				}
			}
		}
		part := fmt.Sprintf("Content-Type: text/html\r\n\n%s\r\n\n", m.HTML)
		message.WriteString(part)
		message.WriteString(embedImages)
	} else {
		part := fmt.Sprintf("Content-Type: text/plain\r\n\n%s\r\n\n", m.Text)
		message.WriteString(part)
	}
	if m.Files != nil {
		for key, value := range m.Files {
			message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
			message.WriteString("Content-Type: application/octect-stream\r\n")
			message.WriteString("Content-Transfer-Encoding:base64\r\n")
			message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"; charset=\"%s\"\r\n\r\n%s\r\n\n", key, m.Charset, value))
		}
		message.WriteString(fmt.Sprintf("--%s--", boundary))
	}
	if c.secure == "SSL" || c.secure == "TLS" {
		return c.SendTLS(m, message)
	}
	return smtp.SendMail(c.host+":"+c.port, c.smtpAuth, m.From, m.To, message.Bytes())
}

//SendTLS 通过TLS发送
func (c *SMTPClient) SendTLS(m Mail, message bytes.Buffer) error {

	var ct *smtp.Client
	var err error
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         c.host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", c.host+":"+c.port, tlsconfig)
	if err != nil {
		log.Println(err, c.host)
		return err
	}

	ct, err = smtp.NewClient(conn, c.host)
	if err != nil {
		log.Println(err)
		return err
	}

	//if err := ct.StartTLS(tlsconfig);err != nil {
	//	log.Println("StartTLS Error:",err,c.host,c.port)
	//	return err
	//}

	//if err := ct.StartTLS(tlsconfig);err != nil {
	//	fmt.Println(err)
	//	return err
	//}

	fmt.Println(c.smtpAuth)
	if ok, s := ct.Extension("AUTH"); ok {
		logs.Info(s)
		// Auth
		if err = ct.Auth(c.smtpAuth); err != nil {
			log.Println("Auth Error:",
				err,
				c.user,
			)
			return err
		}
	}
	// To && From
	if err = ct.Mail(m.From); err != nil {
		log.Println("Mail Error:", err, m.From)
		return err
	}

	for _, v := range m.To {
		if err := ct.Rcpt(v); err != nil {
			log.Println("Rcpt Error:", err, v)
			return err
		}
	}

	// Data
	w, err := ct.Data()
	if err != nil {
		log.Println("Data Object Error:", err)
		return err
	}

	_, err = w.Write(message.Bytes())
	if err != nil {
		log.Println("Write Data Object Error:", err)
		return err
	}

	err = w.Close()
	if err != nil {
		log.Println("Data Object Close Error:", err)
		return err
	}

	ct.Quit()
	return nil
}

// AddTo will take a valid email address and store it in the mail.
// It will return an error if the email is invalid.
func (m *Mail) AddTo(email string) error {
	//Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	parsedAddess, e := mail.ParseAddress(email)
	if e != nil {
		return e
	}
	m.AddRecipient(parsedAddess)
	return nil
}

// SetTos 设置收信人Email地址
func (m *Mail) SetTos(emails []string) {
	m.To = emails
}

// AddToName will add a new receipient name to mail
func (m *Mail) AddToName(name string) {
	m.ToName = append(m.ToName, name)
}

// AddRecipient will take an already parsed mail.Address
func (m *Mail) AddRecipient(receipient *mail.Address) {
	m.To = append(m.To, receipient.Address)
	if len(receipient.Name) > 0 {
		m.ToName = append(m.ToName, receipient.Name)
	}
}

// AddSubject will set the subject of the mail
func (m *Mail) AddSubject(s string) {
	m.Subject = s
}

// AddHTML will set the body of the mail
func (m *Mail) AddHTML(html string) {
	m.HTML = html
}

// AddText will set the body of the email
func (m *Mail) AddText(text string) {
	m.Text = text
}

// AddFrom will set the senders email
func (m *Mail) AddFrom(from string) error {
	//Parses a single RFC 5322 address, e.g. "Barry Gibbs <bg@example.com>"
	parsedAddess, e := mail.ParseAddress(from)
	if e != nil {
		return e
	}
	m.From = parsedAddess.Address
	m.FromName = parsedAddess.Name
	return nil
}

// AddBCC works like AddTo but for BCC
func (m *Mail) AddBCC(email string) error {
	parsedAddess, e := mail.ParseAddress(email)
	if e != nil {
		return e
	}
	m.Bcc = append(m.Bcc, parsedAddess.Address)
	return nil
}

// AddRecipientBCC works like AddRecipient but for BCC
func (m *Mail) AddRecipientBCC(email *mail.Address) {
	m.Bcc = append(m.Bcc, email.Address)
}

// AddFromName will set the senders name
func (m *Mail) AddFromName(name string) {
	m.FromName = name
}

// AddReplyTo will set the return address
func (m *Mail) AddReplyTo(reply string) {
	m.ReplyTo = reply
}

// AddDate specifies the date
func (m *Mail) AddDate(date string) {
	m.Date = date
}

// AddAttachment will include file/s in mail
func (m *Mail) AddAttachment(filePath string) error {
	if m.Files == nil {
		m.Files = make(map[string]string)
	}
	str, err := m.ReadAttachment(filePath)
	if err != nil {
		return err
	}
	_, filename := filepath.Split(filePath)
	m.Files[filename] = str
	return nil
}

// ReadAttachment reading attachment
func (m *Mail) ReadAttachment(filePath string) (string, error) {
	file, e := ioutil.ReadFile(filePath)
	if e != nil {
		return "", e
	}
	encoded := base64.StdEncoding.EncodeToString(file)
	totalChars := len(encoded)
	maxLength := 500 //每行最大长度
	totalLines := totalChars / maxLength
	var buf bytes.Buffer
	for i := 0; i < totalLines; i++ {
		buf.WriteString(encoded[i*maxLength:(i+1)*maxLength] + "\n")
	}
	buf.WriteString(encoded[totalLines*maxLength:])
	return buf.String(), nil
}

// AddHeaders addding header string
func (m *Mail) AddHeaders(headers string) {
	m.Headers = headers
}

// =======================================================
// unencryptedAuth
// =======================================================

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

// ======================================================
// loginAuth
// ======================================================

type loginAuth struct {
	username, password string
}

// LoginAuth loginAuth方式认证
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		return "", nil, errors.New("unencrypted connection")
	}
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
