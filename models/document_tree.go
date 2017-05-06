package models

import (
	"github.com/astaxie/beego/orm"
	"bytes"
	"strconv"
	"github.com/astaxie/beego"
	"html/template"
)

type DocumentTree struct {
	DocumentId int               	`json:"id"`
	DocumentName string 		`json:"text"`
	ParentId interface{}            `json:"parent"`
	Identify string 		`json:"identify"`
	BookIdentify string 		`json:"-"`
	Version int64			`json:"version"`
	State *DocumentSelected         `json:"state,omitempty"`
}
type DocumentSelected struct {
	Selected bool        	`json:"selected"`
	Opened bool        	`json:"opened"`
}

//获取项目的文档树状结构
func (m *Document) FindDocumentTree(book_id int) ([]*DocumentTree,error){
	o := orm.NewOrm()

	trees := make([]*DocumentTree,0)

	var docs []*Document

	count ,err := o.QueryTable(m).Filter("book_id",book_id).OrderBy("order_sort","document_id").All(&docs,"document_id","version","document_name","parent_id","identify")

	if err != nil {
		return trees,err
	}
	book,_ := NewBook().Find(book_id)

	trees = make([]*DocumentTree,count)

	for index,item := range docs {
		tree := &DocumentTree{}
		if index == 0{
			tree.State = &DocumentSelected{ Selected: true, Opened: true }
		}
		tree.DocumentId = item.DocumentId
		tree.Identify = item.Identify
		tree.Version = item.Version
		tree.BookIdentify = book.Identify
		if item.ParentId > 0 {
			tree.ParentId = item.ParentId
		}else{
			tree.ParentId = "#"
		}

		tree.DocumentName = item.DocumentName

		trees[index] = tree
	}

	return trees,nil
}

func (m *Document) CreateDocumentTreeForHtml(book_id, selected_id int) (string,error) {
	trees,err := m.FindDocumentTree(book_id)
	if err != nil {
		return "",err
	}
	parent_id := getSelectedNode(trees,selected_id)

	buf := bytes.NewBufferString("")

	getDocumentTree(trees,0,selected_id,parent_id,buf)

	return buf.String(),nil
	
}

//使用递归的方式获取指定ID的顶级ID
func getSelectedNode(array []*DocumentTree, parent_id int) int {

	for _,item := range array {
		if _,ok := item.ParentId.(string); ok && item.DocumentId == parent_id  {
			return item.DocumentId
		}else if pid,ok := item.ParentId.(int); ok  && item.DocumentId == parent_id{
			return getSelectedNode(array,pid)
		}
	}
	return  0
}

func getDocumentTree(array []*DocumentTree,parent_id int,selected_id int,selected_parent_id int,buf *bytes.Buffer)   {
	buf.WriteString("<ul>")

	for _,item := range array {
		pid := 0

		if p, ok := item.ParentId.(int); ok {
			pid = p
		}
		if pid == parent_id {

			selected := ""
			if item.DocumentId == selected_id {
				selected = ` class="jstree-clicked"`
			}
			selected_li := ""
			if item.DocumentId == selected_parent_id {
				selected_li = ` class="jstree-open"`
			}
			buf.WriteString("<li id=\"")
			buf.WriteString(strconv.Itoa(item.DocumentId))
			buf.WriteString("\"")
			buf.WriteString(selected_li)
			buf.WriteString("><a href=\"")
			if item.Identify != "" {
				uri := beego.URLFor("DocumentController.Read", ":key", item.BookIdentify, ":id", item.Identify)
				buf.WriteString(uri)
			} else {
				uri := beego.URLFor("DocumentController.Read", ":key", item.BookIdentify, ":id", item.DocumentId)
				buf.WriteString(uri)
			}
			buf.WriteString("\" title=\"")
			buf.WriteString(template.HTMLEscapeString(item.DocumentName) + "\"")
			buf.WriteString(selected + ">")
			buf.WriteString(template.HTMLEscapeString(item.DocumentName) + "</a>")

			for _, sub := range array {
				if p, ok := sub.ParentId.(int); ok && p == item.DocumentId {
					getDocumentTree(array, p, selected_id, selected_parent_id, buf)
					break
				}
			}
			buf.WriteString("</li>")

		}
	}
	buf.WriteString("</ul>")
}











