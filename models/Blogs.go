package models

import (
	"time"
	"github.com/lifei6671/mindoc/conf"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

//博文表
type Blog struct {
	BlogId    int		`orm:"pk;auto;unique;column(blog_id)" json:"blog_id"`
	//文章标题
	BlogTitle string	`orm:"column(blog_title);size(500)" json:"blog_title"`
	//文章标识
	BlogIdentify string	`orm:"column(blog_identify);size(100);unique" json:"blog_identify"`
	//排序序号
	OrderIndex int 		`orm:"column(order_index);type(int);default(0)" json:"order_index"`
	//所属用户
	MemberId  int		`orm:"column(member_id);type(int);default(0)" json:"member_id"`
	//文章类型:0 普通文章/1 链接文章
	BlogType int		`orm:"column(blog_type);type(int);default(0)" json:"blog_type"`
	//链接到的项目中的文档ID
	DocumentId int		`orm:"column(document_id);type(int);default(0)" json:"document_id"`
	//文章摘要
	BlogExcerpt string	`orm:"column(blog_excerpt);size(1500);unique" json:"blog_excerpt"`
	//文章内容
	BlogContent string	`orm:"column(blog_content);type(text);null" json:"blog_content"`
	//发布后的文章内容
	BlogRelease string 	`orm:"column(blog_release);type(text);null" json:"blog_release"`
	//文章当前的状态，枚举enum(’publish’,’draft’,’private’,’static’,’object’)值，publish为已 发表，draft为草稿，private为私人内容(不会被公开) ，static(不详)，object(不详)。默认为publish。
	BlogStatus string	`orm:"column(blog_status);size(100);default(publish)" json:"blog_status"`
	//文章密码，varchar(100)值。文章编辑才可为文章设定一个密码，凭这个密码才能对文章进行重新强加或修改。
	Password string		`orm:"column(password);size(100)" json:"-"`
	//最后修改时间
	Modified time.Time	`orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	//修改人id
	ModifyAt int		`orm:"column(modify_at);type(int)" json:"-"`
	//创建时间
	Created time.Time	`orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	Version    int64         `orm:"type(bigint);column(version)" json:"version"`
}

// 多字段唯一键
func (m *Blog) TableUnique() [][]string {
	return [][]string{
		{"blog_id", "blog_identify"},
	}
}

// TableName 获取对应数据库表名.
func (m *Blog) TableName() string {
	return "blogs"
}

// TableEngine 获取数据使用的引擎.
func (m *Blog) TableEngine() string {
	return "INNODB"
}

func (m *Blog) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + m.TableName()
}

func NewBlog() *Blog {
	return &Blog{
		Version: time.Now().Unix(),
	}
}

//根据文章ID查询文章
func (b *Blog) Find(blogId int) (*Blog,error) {
	o := orm.NewOrm()

	err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_id",blogId).One(b)
	if err != nil {
		beego.Error("查询文章时失败 -> ",err)
		return nil,err
	}

	return b,nil
}
//根据文章标识查询文章
func (b *Blog) FindByIdentify(identify string) (*Blog,error) {
	o := orm.NewOrm()

	err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_identify",identify).One(b)
	if err != nil {
		beego.Error("查询文章时失败 -> ",err)
		return nil,err
	}
	return b,nil
}

func (b *Blog)Link() (*Blog,error)  {
	o := orm.NewOrm()
	//如果是链接文章，则需要从链接的项目中查找文章内容
	if b.BlogType == 1 && b.DocumentId > 0{
		doc := NewDocument()
		if err := o.QueryTable(doc.TableNameWithPrefix()).Filter("document_id",b.DocumentId).One(doc,"");err != nil {
			beego.Error("查询文章链接对象时出错 -> ",err)
		}else{
			b.BlogRelease = doc.Release
			//目前仅支持markdown文档进行链接
			b.BlogContent = doc.Markdown
		}
	}

	return b,nil
}

