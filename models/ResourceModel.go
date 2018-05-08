package models

import (
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"errors"
	"github.com/astaxie/beego"
)

type ResourceModel struct {
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
func (m *ResourceModel) TableName() string {
	return "resource"
}

// TableEngine 获取数据使用的引擎.
func (m *ResourceModel) TableEngine() string {
	return "INNODB"
}

// 多字段唯一键
func (m *ResourceModel) TableUnique() [][]string {
	return [][]string{{"resource_group_id", "resource_name","action_name","http_method"}}
}

func (m *ResourceModel) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewResource() *ResourceModel {
	return &ResourceModel{}
}
//添加或更新资源
func (m *ResourceModel) InsertOrUpdate(cols ...string) (err error) {
	if m.ControllerName == "" || m.ActionName == "" || m.ResourceGroupId <= 0 || m.ResourceName == ""{
		return errors.New("参数错误")
	}
	if m.HttpMethod == "" {
		m.HttpMethod = "*"
	}
	o := orm.NewOrm()
	resource := &ResourceModel{}
	//如果设置了资源id，需要先查询是否真实存在
	if m.ResourceId > 0 {
		err = o.QueryTable(m.TableNameWithPrefix()).Filter("resource_id",m.ResourceId).One(resource)
	}
	//如果资源不存在，需要查询是否存在相同的资源
	if err == nil {
		err = o.QueryTable(m.TableNameWithPrefix()).Filter("controller_name",m.ControllerName).Filter("action_name",m.ActionName).Filter("http_method__in",[]string{"*",m.HttpMethod}).One(resource)
		if err == nil {
			return errors.New("资源已存在")
		}
	}

	if err == orm.ErrNoRows {
		_,err = o.Insert(m)
	}else{
		_,err = o.Update(m,cols...)
	}
	if err != nil {
		beego.Error("添加或更新资源时出错 =>",err)
	}
	return
}

//删除资源
func (m *ResourceModel) Delete(resourceId int) (err error) {
	o := orm.NewOrm()

	_,err = o.QueryTable(m.TableNameWithPrefix()).Filter("resource_id",resourceId).Delete()
	if err != nil {
		beego.Error("删除资源时出错 =>",resourceId,err)
	}
	return
}

















