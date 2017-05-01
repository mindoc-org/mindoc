package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/lifei6671/godoc/conf"
)


// Option struct .
type Option struct {
	OptionId int		`orm:"column(option_id);pk;auto;unique;" json:"option_id"`
	OptionTitle string	`orm:"column(option_title);size(500)" json:"option_title"`
	OptionName string	`orm:"column(option_name);unique;size(80)" json:"option_name"`
	OptionValue string	`orm:"column(option_value);type(text);null" json:"option_value"`
	Remark string		`orm:"column(remark);type(text);null" json:"remark"`
}

// TableName 获取对应数据库表名.
func (m *Option) TableName() string {
	return "options"
}
// TableEngine 获取数据使用的引擎.
func (m *Option) TableEngine() string {
	return "INNODB"
}

func (m *Option)TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() +  m.TableName()
}


func NewOption() *Option  {
	return &Option{}
}

func (p *Option) Find(id int) error {
	o := orm.NewOrm()

	p.OptionId = id

	if err := o.Read(p);err != nil {
		return err
	}
	return  nil
}

func (p *Option) FindByKey(key string) error {
	o := orm.NewOrm()

	p.OptionName = key
	if err := o.Read(p);err != nil {
		return err
	}
	return  nil
}

func GetOptionValue(key, def string) string {

	option := NewOption()

	if err := option.FindByKey(key); err == nil {
		return option.OptionValue
	}
	return def
}

func (p *Option) InsertOrUpdate() error  {

	o := orm.NewOrm()

	var err error
	if p.OptionId > 0 {
		_,err = o.Update(p)
	}else{
		_,err = o.Insert(p)
	}
	return err
}

func (p *Option) InsertMulti(option... Option )  (error){

	o := orm.NewOrm()

	_,err := o.InsertMulti(len(option),option)
	return err
}

func (p *Option) All() ([]*Option,error)  {
	o := orm.NewOrm()
	var options []*Option

	_,err := o.QueryTable(p.TableNameWithPrefix()).All(&options)

	if err != nil {
		return options,err
	}
	return options,nil
}