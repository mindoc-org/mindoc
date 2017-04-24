package models

import (
	"time"
	"github.com/lifei6671/godoc/conf"
)

const (
	Logger_Operate = "operate"
	Logger_System = "system"
	Logger_Exception = "exception"
)

// Logger struct .
type Logger struct {
	LoggerId int64		`orm:"pk;auto;unique;column(logger_id)" json:"logger_id"`
	MemberId int		`orm:"column(member_id);type(int)" json:"member_id"`
	// 日志类别：operate 操作日志/ system 系统日志/ exception 异常日志
	Category string		`orm:"column(category);size(255);default(operate)" json:"category"`
	Content string		`orm:"column(content);type(text)" json:"content"`
	OriginalData string	`orm:"column(original_data);type(text)" json:"original_data"`
	PresentData string	`orm:"column(present_data);type(text)" json:"present_data"`
	CreateTime time.Time	`orm:"type(datetime);column(create_time);auto_now_add" json:"create_time"`
	UserAgent string	`orm:"column(user_agent);size(500)" json:"user_agent"`
	IPAddress string	`orm:"column(ip_address);size(255)" json:"ip_address"`
}

// TableName 获取对应数据库表名.
func (m *Logger) TableName() string {
	return "logs"
}
// TableEngine 获取数据使用的引擎.
func (m *Logger) TableEngine() string {
	return "INNODB"
}
func (m *Logger) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

















