package models

import (
	"github.com/astaxie/beego/orm"
)

type DocumentTree struct {
	DocumentId int               	`json:"id"`
	DocumentName string 		`json:"text"`
	ParentId interface{}            `json:"parent"`
	Identify string 		`json:"identify"`
	Version int64			`json:"version"`
	State *DocumentSelected         `json:"state,omitempty"`
}
type DocumentSelected struct {
	Selected bool        	`json:"selected"`
	Opened bool        	`json:"opened"`
}


func (m *Document) FindDocumentTree(book_id int) ([]*DocumentTree,error){
	o := orm.NewOrm()

	trees := make([]*DocumentTree,0)

	var docs []*Document

	count ,err := o.QueryTable(m).Filter("book_id",book_id).OrderBy("order_sort","document_id").All(&docs,"document_id","version","document_name","parent_id","identify")

	if err != nil {
		return trees,err
	}

	trees = make([]*DocumentTree,count)

	for index,item := range docs {
		tree := &DocumentTree{}
		if index == 0{
			tree.State = &DocumentSelected{ Selected: true, Opened: true }
		}
		tree.DocumentId = item.DocumentId
		tree.Identify = item.Identify
		tree.Version = item.Version
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
