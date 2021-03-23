package utils

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/logs"
	"gopkg.in/ldap.v2"
)

/*
对应的config
ldap:
  host: hostname.yourdomain.com //ldap服务器地址
  port: 3268 //ldap服务器端口
  attribute: mail //用户名对应ldap object属性
  base: DC=yourdomain,DC=com //搜寻范围
  user: CN=ldap helper,OU=yourdomain.com,DC=yourdomain,DC=com //第一次绑定用户
  password: p@sswd //第一次绑定密码
  ssl: false //使用使用ssl
*/

func ValidLDAPLogin(password string) (result bool, err error) {
	result = false
	err = nil
	lc, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.3.104", 389))
	if err != nil {
		logs.Error("Dial => ", err)
		return
	}

	defer lc.Close()
	err = lc.Bind("cn=admin,dc=minho,dc=com", "123456")
	if err != nil {
		logs.Error("Bind => ", err)
		return
	}
	searchRequest := ldap.NewSearchRequest(
		"DC=minho,DC=com",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=User)(%s=%s))", "mail", "longfei6671@163.com"),
		[]string{"dn"},
		nil,
	)
	searchResult, err := lc.Search(searchRequest)
	if err != nil {
		logs.Error("Search => ", err)
		return
	}
	if len(searchResult.Entries) != 1 {
		err = errors.New("ldap.no_user_found_or_many_users_found")
		return
	}
	fmt.Printf("%+v = %d", searchResult.Entries, len(searchResult.Entries))

	userdn := searchResult.Entries[0].DN

	err = lc.Bind(userdn, password)
	if err == nil {
		result = true
	} else {
		logs.Error("Bind2 => ", err)
		err = nil
	}
	return
}

func AddMember(account, password string) error {
	lc, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.3.104", 389))
	if err != nil {
		logs.Error("Dial => ", err)
		return err
	}

	defer lc.Close()
	user := fmt.Sprintf("cn=%s,dc=minho,dc=com", account)

	member := ldap.NewAddRequest(user)

	member.Attribute("mail", []string{"longfei6671@163.com"})

	err = lc.Add(member)

	if err == nil {

		err = lc.Bind(user, "")
		if err != nil {
			logs.Error("Bind => ", err)
			return err
		}
		passwordModifyRequest := ldap.NewPasswordModifyRequest(user, "", "1q2w3e__ABC")
		_, err = lc.PasswordModify(passwordModifyRequest)

		if err != nil {
			logs.Error("PasswordModify => ", err)
			return err
		}
		return nil
	}
	logs.Error("Add => ", err)
	return err
}

func ModifyPassword(account, old_password, new_password string) error {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", "192.168.3.104", 389))
	if err != nil {
		logs.Error("Dial => ", err)
	}
	defer l.Close()

	user := fmt.Sprintf("cn=%s,dc=minho,dc=com", account)
	err = l.Bind(user, old_password)
	if err != nil {
		logs.Error("Bind => ", err)
		return err
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest(user, old_password, new_password)
	_, err = l.PasswordModify(passwordModifyRequest)

	if err != nil {
		logs.Error(fmt.Sprintf("Password could not be changed: %s", err.Error()))
		return err
	}
	return nil
}
