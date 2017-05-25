// Package models 为项目所需的模型对象定义.
package models

import "errors"

var (
	// ErrMemberNoExist 用户不存在.
	ErrMemberNoExist  = errors.New("用户不存在")
	ErrMemberDisabled = errors.New("用户被禁用")
	// ErrorMemberPasswordError 密码错误.
	ErrorMemberPasswordError = errors.New("用户密码错误")
	//ErrorMemberAuthMethodInvalid 不支持此认证方式
	ErrMemberAuthMethodInvalid = errors.New("不支持此认证方式")
	//ErrLDAPConnect 无法连接到LDAP服务器
	ErrLDAPConnect = errors.New("无法连接到LDAP服务器")
	//ErrLDAPFirstBind 第一次LDAP绑定失败
	ErrLDAPFirstBind = errors.New("第一次LDAP绑定失败")
	//ErrLDAPSearch LDAP搜索失败
	ErrLDAPSearch = errors.New("LDAP搜索失败")
	//ErrLDAPUserNotFoundOrTooMany
	ErrLDAPUserNotFoundOrTooMany = errors.New("LDAP用户不存在或者多于一个")

	// ErrDataNotExist 指定的服务已存在.
	ErrDataNotExist = errors.New("数据不存在")

	// ErrInvalidParameter 参数错误.
	ErrInvalidParameter = errors.New("Invalid parameter")

	ErrPermissionDenied = errors.New("Permission denied")

	ErrCommentClosed          = errors.New("评论已关闭")
	ErrCommentContentNotEmpty = errors.New("评论内容不能为空")
)
