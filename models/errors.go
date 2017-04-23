// Package models 为项目所需的模型对象定义.
package models

import "errors"

var(
	// ErrMemberNoExist 用户不存在.
	ErrMemberNoExist = errors.New("用户不存在")
	ErrMemberDisabled = errors.New("用户被禁用")
	// ErrorMemberPasswordError 密码错误.
	ErrorMemberPasswordError = errors.New("用户密码错误")
	// ErrDataNotExist 指定的服务已存在.
	ErrDataNotExist = errors.New("数据不存在")

	// ErrInvalidParameter 参数错误.
	ErrInvalidParameter = errors.New("Invalid parameter")

	ErrPermissionDenied = errors.New("Permission denied")

	ErrCommentClosed = errors.New("评论已关闭")
	ErrCommentContentNotEmpty = errors.New("评论内容不能为空")
)
