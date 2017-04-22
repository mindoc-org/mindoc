package models

import "github.com/astaxie/beego/orm"

type Model struct {

}

func (m *Model) Find(id int) error {
	o := orm.NewOrm()

	return o.Read(m)
}

func (m *Model) Insert() error  {
	o := orm.NewOrm()
	_,err := o.Insert(m)

	return err
}