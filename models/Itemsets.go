package models

import (
	"time"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"errors"
	"github.com/lifei6671/mindoc/utils/cryptil"
)

//项目集
type Itemsets struct {
	ItemId      int       `orm:"column(item_id);pk;auto;unique" json:"item_id"`
	ItemName    string    `orm:"column(item_name);size(500)" json:"item_name"`
	ItemKey     string    `orm:"column(item_key);size(200);unique" json:"item_key"`
	Description string    `orm:"column(description);type(text);null" json:"description"`
	MemberId    int       `orm:"column(member_id);size(100)" json:"member_id"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	ModifyTime  time.Time `orm:"column(modify_time);type(datetime);null;auto_now" json:"modify_time"`
	ModifyAt    int       `orm:"column(modify_at);type(int)" json:"modify_at"`

	BookNumber       int    `orm:"-" json:"book_number"`
	CreateTimeString string `orm:"-" json:"create_time_string"`
}

// TableName 获取对应数据库表名.
func (item *Itemsets) TableName() string {
	return "itemsets"
}

// TableEngine 获取数据使用的引擎.
func (item *Itemsets) TableEngine() string {
	return "INNODB"
}
func (item *Itemsets) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + item.TableName()
}

func (item *Itemsets) QueryTable() orm.QuerySeter {
	return orm.NewOrm().QueryTable(item.TableNameWithPrefix())
}

func NewItemsets() *Itemsets {
	return &Itemsets{}
}

func (item *Itemsets) First(itemId int) (*Itemsets, error) {
	if itemId <= 0 {
		return nil, ErrInvalidParameter
	}
	err := item.QueryTable().Filter("item_id", itemId).One(item)
	if err != nil {
		beego.Error("查询项目集失败 -> item_id=", itemId, err)
	} else {
		item.Include()
	}
	return item, err
}

func (item *Itemsets) FindFirst(itemKey string) (*Itemsets,error) {
	err := item.QueryTable().Filter("item_key",itemKey).One(item)
	if err != nil {
		beego.Error("查询项目集失败 -> itemKey=", itemKey, err)
	} else {
		item.Include()
	}
	return item,err
}

func (item *Itemsets) Exist(itemId int) bool {
	return item.QueryTable().Filter("item_id", itemId).Exist()
}

//保存
func (item *Itemsets) Save() (err error) {

	if item.ItemName == "" {
		return errors.New("项目集名称不能为空")
	}
	if item.ItemKey == "" {
		item.ItemKey = cryptil.NewRandChars(16)
	}

	if item.QueryTable().Filter("item_id__ne", item.ItemId).Filter("item_key", item.ItemKey).Exist() {
		return errors.New("项目集标识已存在")
	}
	if item.ItemId > 0 {
		_, err = orm.NewOrm().Update(item)
	} else {
		_, err = orm.NewOrm().Insert(item)
	}
	return
}

//删除.
func (item *Itemsets) Delete(itemId int) (err error) {
	if itemId <= 0 {
		return ErrInvalidParameter
	}
	if itemId == 1 {
		return errors.New("默认项目集不能删除")
	}
	if !item.Exist(itemId) {
		return errors.New("项目集不存在")
	}
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		beego.Error("开启事物失败 ->", err)
		return err
	}
	_, err = o.QueryTable(item.TableNameWithPrefix()).Filter("item_id", itemId).Delete()
	if err != nil {
		beego.Error("删除项目集失败 -> item_id=", itemId, err)
		o.Rollback()
	}
	_, err = o.Raw("update md_books set item_id=1 where item_id=?;", itemId).Exec()
	if err != nil {
		beego.Error("删除项目集失败 -> item_id=", itemId, err)
		o.Rollback()
	}

	return o.Commit()
}

func (item *Itemsets) Include() (*Itemsets, error) {

	item.CreateTimeString = item.CreateTime.Format("2006-01-02 15:04:05")

	i, err := NewBook().QueryTable().Filter("item_id", item.ItemId).Count()
	if err != nil && err != orm.ErrNoRows {
		return item, err
	}
	item.BookNumber = int(i)

	return item, nil
}

//分页查询.
func (item *Itemsets) FindToPager(pageIndex, pageSize int) (list []*Itemsets, totalCount int, err error) {

	offset := (pageIndex - 1) * pageSize

	_, err = item.QueryTable().OrderBy("-item_id").Offset(offset).Limit(pageSize).All(&list)

	if err != nil {
		return
	}

	c, err := item.QueryTable().Count()
	if err != nil {
		return
	}
	totalCount = int(c)

	for _, item := range list {
		item.Include()
	}
	return
}

//根据项目集名称查询.
func (item *Itemsets) FindItemsetsByName(name string, limit int) (*SelectMemberResult, error) {
	result := SelectMemberResult{}

	var itemsets []*Itemsets
	var err error
	if name == "" {
		_, err = item.QueryTable().Limit(limit).All(&itemsets)

	} else {
		_, err = item.QueryTable().Filter("item_name__icontains", name).Limit(limit).All(&itemsets)
	}
	if err != nil {
		beego.Error("查询项目集失败 ->", err)
		return &result, err
	}

	items := make([]KeyValueItem, 0)

	for _, m := range itemsets {
		item := KeyValueItem{}
		item.Id = m.ItemId
		item.Text = m.ItemName
		items = append(items, item)
	}
	result.Result = items

	return &result, err
}

//根据项目集标识查询项目集的项目列表.
func (item *Itemsets) FindItemsetsByItemKey(key string, pageIndex, pageSize, memberId int) (books []*BookResult, totalCount int, err error){
	o := orm.NewOrm()

	err = item.QueryTable().Filter("item_key",key).One(item)

	if err != nil {
		return nil,0,err
	}
	offset := (pageIndex - 1) * pageSize
	//如果是登录用户
	if memberId > 0 {
		sql1 := `SELECT COUNT(*)
FROM md_books AS book
  LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
  left join (select *
             from (select book_id,team_member_id,role_id
                   from md_team_relationship as mtr
                     left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.book_id) as team on team.book_id = book.book_id
WHERE book.item_id = ? AND (relationship_id > 0 OR book.privately_owned = 0 or team.team_member_id > 0)`

		err = o.Raw(sql1, memberId, memberId, item.ItemId).QueryRow(&totalCount)
		if err != nil {
			return
		}
		sql2 := `SELECT book.*,rel1.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.member_id = ?
			left join (select * from (select book_id,team_member_id,role_id
                   	from md_team_relationship as mtr
					left join md_team_member as mtm on mtm.team_id=mtr.team_id and mtm.member_id=? order by role_id desc )as t group by t.book_id) as team 
					on team.book_id = book.book_id
			LEFT JOIN md_relationship AS rel1 ON rel1.book_id = book.book_id AND rel1.role_id = 0
			LEFT JOIN md_members AS member ON rel1.member_id = member.member_id
			WHERE book.item_id = ? AND (rel.relationship_id > 0 OR book.privately_owned = 0 or team.team_member_id > 0) 
			ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql2, memberId, memberId, item.ItemId, offset, pageSize).QueryRows(&books)

		return

	} else {
		count, err1 := o.QueryTable(NewBook().TableNameWithPrefix()).Filter("privately_owned", 0).Filter("item_id", item.ItemId).Count()

		if err1 != nil {
			err = err1
			return
		}
		totalCount = int(count)

		sql := `SELECT book.*,rel.*,member.account AS create_name FROM md_books AS book
			LEFT JOIN md_relationship AS rel ON rel.book_id = book.book_id AND rel.role_id = 0
			LEFT JOIN md_members AS member ON rel.member_id = member.member_id
			WHERE book.item_id = ? AND book.privately_owned = 0 ORDER BY order_index DESC ,book.book_id DESC LIMIT ?,?`

		_, err = o.Raw(sql, item.ItemId, offset, pageSize).QueryRows(&books)

		return

	}
}