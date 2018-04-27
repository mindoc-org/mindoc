package models

import (
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"errors"
	"github.com/astaxie/beego"
)

type Resource struct {
	//主键
	ResourceId int		 		`orm:"column(resource_id);pk;auto;unique;" json:"resource_id"`
	//分组ID
	ResourceGroupId int  		`orm:"column(resource_group_id);index" json:"resource_group_id"`
	//分组名称
	ResourceGroupName string 	`orm:"-" json:"resource_group_name"`
	//资源名称
	ResourceName string  		`orm:"column(resource_name);size(255)" json:"resource_name"`
	ControllerName string 		`orm:"column(controller_name);size(255)" json:"controller_name"`
	ActionName string			`orm:"column(action_name);size(255)" json:"action_name"`
	HttpMethod string			`orm:"column(http_method);size(50)" json:"http_method"`
}

// TableName 获取对应数据库表名.
func (m *Resource) TableName() string {
	return "resource"
}

// TableEngine 获取数据使用的引擎.
func (m *Resource) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (m *Resource) TableUnique() [][]string {
	return [][]string{{"resource_group_id", "resource_name","action_name","http_method"}}
}

func (m *Resource) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewResource() *Resource {
	return &Resource{}
}
//添加或更新资源
func (m *Resource) InsertOrUpdate(cols ...string) (err error) {
	if m.ControllerName == "" || m.ActionName == "" || m.ResourceGroupId <= 0 || m.ResourceName == ""{
		return errors.New("参数错误")
	}
	if m.HttpMethod == "" {
		m.HttpMethod = "GET"
	}
	o := orm.NewOrm()

	if m.ResourceId > 0 {
		_,err = o.Update(m,cols...)
	}else{
		_,err = o.Insert(m)
	}
	if err != nil {
		beego.Error("添加或更新资源时出错 =>",err)
	}
	return
}

//删除资源
func (m *Resource) Delete(resourceId int) (err error) {
	o := orm.NewOrm()

	_,err = o.QueryTable(m.TableNameWithPrefix()).Filter("resource_id",resourceId).Delete()
	if err != nil {
		beego.Error("删除资源时出错 =>",resourceId,err)
	}
	return
}





















