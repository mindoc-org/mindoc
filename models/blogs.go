package models

import "time"

//博文表
type Blog struct {
	BlogId    int
	BlogTitle string
	MemberId  int
	//文字摘要
	BlogExcerpt string
	//文章内容
	BlogContent string

	//文章当前的状态，枚举enum(’publish’,’draft’,’private’,’static’,’object’)值，publish为已 发表，draft为草稿，private为私人内容(不会被公开) ，static(不详)，object(不详)。默认为publish。
	BlogStatus string
	//文章密码，varchar(20)值。文章编辑才可为文章设定一个密码，凭这个密码才能对文章进行重新强加或修改。
	Password string
	//最后修改时间
	Modified time.Time
	//修改人id
	ModifyAt int
	//创建时间
	Created time.Time
}
