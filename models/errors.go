// Package models 为项目所需的模型对象定义.
package models

import "errors"

var(
	// ErrMemberNoExist 用户不存在.
	ErrMemberNoExist = errors.New("用户不存在")
	// ErrorMemberPasswordError 密码错误.
	ErrorMemberPasswordError = errors.New("用户密码错误")
	// ErrServerAlreadyExist 指定的服务已存在.
	ErrServerAlreadyExist = errors.New("服务已存在")

	// ErrInvalidParameter 参数错误.
	ErrInvalidParameter = errors.New("Invalid parameter")
)
