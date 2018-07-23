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
	MemberId  int		`orm:"column(member_id);type(int);default(0):index" json:"member_id"`
	//用户头像
	MemberAvatar string		`orm:"-" json:"member_avatar"`
	//文章类型:0 普通文章/1 链接文章
	BlogType int		`orm:"column(blog_type);type(int);default(0)" json:"blog_type"`
	//链接到的项目中的文档ID
	DocumentId int		`orm:"column(document_id);type(int);default(0)" json:"document_id"`
	//文章的标识
	DocumentIdentify string `orm:"-" json:"document_identify"`
	//关联文档的项目标识
	BookIdentify string 	`orm:"-" json:"book_identify"`
	//关联文档的项目ID
	BookId int 				`orm:"-" json:"book_id"`
	//文章摘要
	BlogExcerpt string	`orm:"column(blog_excerpt);size(1500)" json:"blog_excerpt"`
	//文章内容
	BlogContent string	`orm:"column(blog_content);type(text);null" json:"blog_content"`
	//发布后的文章内容
	BlogRelease string 	`orm:"column(blog_release);type(text);null" json:"blog_release"`
	//文章当前的状态，枚举enum(’publish’,’draft’,’password’)值，publish为已 发表，draft为草稿，password 为私人内容(不会被公开) 。默认为publish。
	BlogStatus string	`orm:"column(blog_status);size(100);default(publish)" json:"blog_status"`
	//文章密码，varchar(100)值。文章编辑才可为文章设定一个密码，凭这个密码才能对文章进行重新强加或修改。
	Password string		`orm:"column(password);size(100)" json:"-"`
	//最后修改时间
	Modified time.Time	`orm:"column(modify_time);type(datetime);auto_now" json:"modify_time"`
	//修改人id
	ModifyAt int		`orm:"column(modify_at);type(int)" json:"-"`
	ModifyRealName string `orm:"-" json:"modify_real_name"`
	//创建时间
	Created time.Time	`orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	CreateName string 	`orm:"-" json:"create_name"`
	//版本号
	Version    int64    `orm:"type(bigint);column(version)" json:"version"`
	//附件列表
	AttachList []*Attachment `orm:"-" json:"attach_list"`
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
		BlogStatus: "public",
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


	return b.Link()
}
//查找指定用户的指定文章
func (b *Blog) FindByIdAndMemberId(blogId,memberId int) (*Blog,error) {
	o := orm.NewOrm()

	err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_id",blogId).Filter("member_id",memberId).One(b)
	if err != nil {
		beego.Error("查询文章时失败 -> ",err)
		return nil,err
	}

	return b.Link()
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

//获取指定文章的链接内容
func (b *Blog)Link() (*Blog,error)  {
	o := orm.NewOrm()
	//如果是链接文章，则需要从链接的项目中查找文章内容
	if b.BlogType == 1 && b.DocumentId > 0 {
		doc := NewDocument()
		if err := o.QueryTable(doc.TableNameWithPrefix()).Filter("document_id",b.DocumentId).One(doc,"release","markdown","identify","book_id");err != nil {
			beego.Error("查询文章链接对象时出错 -> ",err)
		}else{
			b.DocumentIdentify = doc.Identify
			b.BlogRelease = doc.Release
			//目前仅支持markdown文档进行链接
			b.BlogContent = doc.Markdown
			book := NewBook()
			if err := o.QueryTable(book.TableNameWithPrefix()).Filter("book_id",doc.BookId).One(book,"identify");err != nil {
				beego.Error("查询关联文档的项目时出错 ->",err)
			}else{
				b.BookIdentify = book.Identify
				b.BookId = doc.BookId
			}
		}
	}

	if b.ModifyAt > 0{
		member := NewMember()
		if err := o.QueryTable(member.TableNameWithPrefix()).Filter("member_id",b.ModifyAt).One(member,"real_name","account"); err == nil {
			if member.RealName != ""{
				b.ModifyRealName = member.RealName
			}else{
				b.ModifyRealName = member.Account
			}
		}
	}
	if b.MemberId > 0 {
		member := NewMember()
		if err := o.QueryTable(member.TableNameWithPrefix()).Filter("member_id",b.MemberId).One(member,"real_name","account","avatar"); err == nil {
			if member.RealName != ""{
				b.CreateName = member.RealName
			}else{
				b.CreateName = member.Account
			}
			b.MemberAvatar = member.Avatar
		}
	}

	return b,nil
}

//判断指定的文章标识是否存在
func (b *Blog) IsExist(identify string) bool {
	o := orm.NewOrm()

	return o.QueryTable(b.TableNameWithPrefix()).Filter("blog_identify",identify).Exist()
}

//保存文章
func (b *Blog) Save(cols ...string) error {
	o := orm.NewOrm()

	if b.OrderIndex <= 0 {
		blog := NewBlog()
		if err := o.QueryTable(blog.TableNameWithPrefix()).OrderBy("-blog_id").Limit(1).One(blog,"blog_id");err == nil{
			b.OrderIndex = blog.BlogId + 1;
		}else{
			c,_ := o.QueryTable(b.TableNameWithPrefix()).Count()
			b.OrderIndex = int(c) + 1
		}
	}
	var err error
	b.Version = time.Now().Unix()
	if b.BlogId > 0 {
		b.Modified = time.Now()
		_,err = o.Update(b,cols...)
	}else{

		b.Created = time.Now()
		_,err = o.Insert(b)
	}
	return err
}

//分页查询文章列表
func (b *Blog) FindToPager(pageIndex, pageSize int,memberId int,status string) (blogList []*Blog, totalCount int, err error) {

	o := orm.NewOrm()

	offset := (pageIndex - 1) * pageSize

	query := o.QueryTable(b.TableNameWithPrefix());

	if memberId > 0 {
		query = query.Filter("member_id",memberId)
	}
	if status != "" {
		query = query.Filter("blog_status",status)
	}


	_,err = query.OrderBy("-order_index","-blog_id").Offset(offset).Limit(pageSize).All(&blogList)

	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
		}
		beego.Error("获取文章列表时出错 ->",err)
		return
	}
	count,err := query.Count()

	if err != nil {
		beego.Error("获取文章数量时出错 ->",err)
		return nil,0,err
	}
	totalCount = int(count)
	for _,blog := range blogList {
		if blog.BlogType == 1 {
			blog.Link()
		}
	}

	return
}

//删除文章
func (b *Blog) Delete(blogId int) error {
	o := orm.NewOrm()

	_,err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_id",blogId).Delete()
	if err != nil {
		beego.Error("删除文章失败 ->",err)
	}
	return err
}

//查询下一篇文章
func (b *Blog) QueryNext(blogId int) (*Blog,error) {
	o := orm.NewOrm()
	blog := NewBlog()

	if err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_id",blogId).One(blog,"order_index"); err != nil {
		beego.Error("查询文章时出错 ->",err)
		return b,err
	}

	err := o.QueryTable(b.TableNameWithPrefix()).Filter("order_index__gte",blog.OrderIndex).Filter("blog_id__gt",blogId).OrderBy("-order_index","-blog_id").One(blog)

	if err != nil {
		beego.Error("查询文章时出错 ->",err)
	}
	return blog,err
}

//查询下一篇文章
func (b *Blog) QueryPrevious(blogId int) (*Blog,error) {
	o := orm.NewOrm()
	blog := NewBlog()

	if err := o.QueryTable(b.TableNameWithPrefix()).Filter("blog_id",blogId).One(blog,"order_index"); err != nil {
		beego.Error("查询文章时出错 ->",err)
		return b,err
	}

	err := o.QueryTable(b.TableNameWithPrefix()).Filter("order_index__lte",blog.OrderIndex).Filter("blog_id__lt",blogId).OrderBy("-order_index","-blog_id").One(blog)

	if err != nil {
		beego.Error("查询文章时出错 ->",err)
	}
	return blog,err
}

//关联文章附件
func (b *Blog) LinkAttach() (err error) {

	o := orm.NewOrm()

	var attachList []*Attachment
	//当不是关联文章时，用文章ID去查询附件
	if b.BlogType != 1 || b.DocumentId <= 0 {
		_, err = o.QueryTable(NewAttachment().TableNameWithPrefix()).Filter("document_id", b.BlogId).Filter("book_id",0).All(&attachList)
		if err != nil {
			beego.Error("查询文章附件时出错 ->", err)
		}
	}else {
		_, err = o.QueryTable(NewAttachment().TableNameWithPrefix()).Filter("document_id", b.DocumentId).Filter("book_id", b.BookId).All(&attachList)

		if err != nil {
			beego.Error("查询文章附件时出错 ->", err)
		}
	}
	b.AttachList = attachList
	return
}