package controllers

import (
	"sort"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/i18n"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/pagination"
	"github.com/mindoc-org/mindoc/utils/segmenter"
	"github.com/mindoc-org/mindoc/utils/sqltil"
)

// SearchV2Result 用于 IndexV2 的搜索结果结构
type SearchV2Result struct {
	SearchType   string      `json:"search_type"`
	DocumentId   int         `json:"doc_id"`
	DocumentName string      `json:"doc_name"`
	Identify     string      `json:"identify"`
	Description  string      `json:"description"`
	Author       string      `json:"author"`
	ModifyTime   interface{} `json:"modify_time"`
	CreateTime   interface{} `json:"create_time"`
	BookId       int         `json:"book_id"`
	BookName     string      `json:"book_name"`
	BookIdentify string      `json:"book_identify"`
}

// SearchV2RawResult 底层搜索返回的原始结果，包含倒排索引的分数信息
type SearchV2RawResult struct {
	ContentType  int         // 1-Document 2-Blog
	ContentId    int         // 文档ID或博客ID
	Score        float64     // TF-IDF分数
	WordCounts   []int       // 各个词的词频
	SearchType   string      // "document" 或 "blog"
	DocumentId   int         // 文档ID
	DocumentName string      // 文档名称/博客标题
	Identify     string      // 文档标识/博客标识
	Description  string      // 描述
	Content      string      // 原始内容
	Author       string      // 作者
	ModifyTime   interface{} // 修改时间
	CreateTime   interface{} // 创建时间
	BookId       int         // 项目ID (仅Document)
	BookName     string      // 项目名称 (仅Document)
	BookIdentify string      // 项目标识 (仅Document)
	BlogId       int         // 博客ID (仅Blog)
	BlogTitle    string      // 博客标题 (仅Blog)
	BlogIdentify string      // 博客标识 (仅Blog)
	BlogExcerpt  string      // 博客摘要 (仅Blog)
}

// PerformSearchV2Raw 执行倒排索引搜索的底层函数，返回原始结果
func PerformSearchV2Raw(keyword string, pageIndex, pageSize int, memberId int) ([]*SearchV2RawResult, []string, int, error) {
	// 使用分词器对关键词进行分词
	words := segmenter.Segment(keyword)
	if len(words) == 0 {
		// 如果分词结果为空，直接使用原关键词
		words = []string{keyword}
	}

	// 将原始关键词（小写）加入搜索词列表，确保能匹配索引中存储的完整词条
	lowerKeyword := strings.ToLower(strings.TrimSpace(keyword))
	if lowerKeyword != "" {
		found := false
		for _, w := range words {
			if w == lowerKeyword {
				found = true
				break
			}
		}
		if !found {
			words = append(words, lowerKeyword)
		}
	}

	// 使用倒排索引模型进行搜索
	index := models.NewContentReverseIndex()
	allResults, err := index.FindByWords(words)
	if err != nil {
		return nil, words, 0, err
	}

	// 收集需要批量查询的ID
	docIds := make([]int, 0)
	blogIds := make([]int, 0)
	for _, result := range allResults {
		if result.ContentType == 1 {
			docIds = append(docIds, result.ContentId)
		} else if result.ContentType == 2 {
			blogIds = append(blogIds, result.ContentId)
		}
	}

	// 批量加载 Document
	docMap := make(map[int]*models.Document)
	if len(docIds) > 0 {
		var docs []*models.Document
		o := orm.NewOrm()
		_, err := o.QueryTable(models.NewDocument().TableNameWithPrefix()).Filter("document_id__in", docIds).All(&docs)
		if err == nil {
			for _, doc := range docs {
				docMap[doc.DocumentId] = doc
			}
		}
	}

	// 批量加载 Blog
	blogMap := make(map[int]*models.Blog)
	if len(blogIds) > 0 {
		var blogs []*models.Blog
		o := orm.NewOrm()
		_, err := o.QueryTable(models.NewBlog().TableNameWithPrefix()).Filter("blog_id__in", blogIds).All(&blogs)
		if err == nil {
			for _, blog := range blogs {
				blogMap[blog.BlogId] = blog
			}
		}
	}

	// 收集需要加载的 BookId 和 MemberId
	bookIds := make([]int, 0)
	memberIds := make([]int, 0)
	bookIdSet := make(map[int]bool)
	memberIdSet := make(map[int]bool)
	for _, doc := range docMap {
		if doc.BookId > 0 && !bookIdSet[doc.BookId] {
			bookIds = append(bookIds, doc.BookId)
			bookIdSet[doc.BookId] = true
		}
		if doc.MemberId > 0 && !memberIdSet[doc.MemberId] {
			memberIds = append(memberIds, doc.MemberId)
			memberIdSet[doc.MemberId] = true
		}
	}
	for _, blog := range blogMap {
		if blog.MemberId > 0 && !memberIdSet[blog.MemberId] {
			memberIds = append(memberIds, blog.MemberId)
			memberIdSet[blog.MemberId] = true
		}
	}

	// 批量加载 Book
	bookMap := make(map[int]*models.Book)
	if len(bookIds) > 0 {
		var books []*models.Book
		o := orm.NewOrm()
		_, err := o.QueryTable(models.NewBook().TableNameWithPrefix()).Filter("book_id__in", bookIds).All(&books)
		if err == nil {
			for _, book := range books {
				bookMap[book.BookId] = book
			}
		}
	}

	// 批量加载 Member
	memberMap := make(map[int]*models.Member)
	if len(memberIds) > 0 {
		var members []*models.Member
		o := orm.NewOrm()
		_, err := o.QueryTable(models.NewMember().TableNameWithPrefix()).Filter("member_id__in", memberIds).All(&members, "member_id", "account", "real_name")
		if err == nil {
			for _, member := range members {
				memberMap[member.MemberId] = member
			}
		}
	}

	// 构建返回结果
	searchResults := make([]*SearchV2RawResult, 0)
	for _, result := range allResults {
		item := &SearchV2RawResult{
			ContentType: result.ContentType,
			ContentId:   result.ContentId,
			Score:       result.Score,
			WordCounts:  result.WordCounts,
		}

		// 根据内容类型获取详细信息
		if result.ContentType == 1 {
			doc, ok := docMap[result.ContentId]
			if !ok {
				continue
			}
			book, ok := bookMap[doc.BookId]
			if !ok {
				continue
			}

			item.SearchType = "document"
			item.DocumentId = doc.DocumentId
			item.DocumentName = doc.DocumentName
			item.BookId = doc.BookId
			item.BookName = book.BookName
			item.Identify = doc.Identify
			item.BookIdentify = book.Identify
			item.CreateTime = doc.CreateTime
			item.ModifyTime = doc.ModifyTime
			item.Content = doc.Release

			// 获取作者信息
			if member, ok := memberMap[doc.MemberId]; ok {
				if member.RealName != "" {
					item.Author = member.RealName
				} else {
					item.Author = member.Account
				}
			}

			// 提取描述
			description := doc.Release
			if description == "" {
				description = doc.Markdown
			}
			description = utils.StripTags(description)
			if len([]rune(description)) > 100 {
				description = string([]rune(description)[:100]) + "..."
			}
			item.Description = description

			// 标题匹配加权：搜索词命中标题时提升分数
			titleLower := strings.ToLower(doc.DocumentName)
			for _, w := range words {
				if strings.Contains(titleLower, w) {
					item.Score *= 1.5
				}
			}

			// 精确匹配加权：文档内容包含原始关键词时大幅提分
			if lowerKeyword != "" {
				contentLower := strings.ToLower(utils.StripTags(doc.Release))
				if contentLower == "" {
					contentLower = strings.ToLower(utils.StripTags(doc.Markdown))
				}
				if strings.Contains(contentLower, lowerKeyword) {
					item.Score *= 5.0
				}
			}

			searchResults = append(searchResults, item)
		} else if result.ContentType == 2 {
			blog, ok := blogMap[result.ContentId]
			if !ok {
				continue
			}

			item.SearchType = "blog"
			item.BlogId = blog.BlogId
			item.BlogTitle = blog.BlogTitle
			item.DocumentId = blog.BlogId
			item.DocumentName = blog.BlogTitle
			item.BlogIdentify = blog.BlogIdentify
			item.Identify = blog.BlogIdentify
			item.BlogExcerpt = blog.BlogExcerpt
			item.CreateTime = blog.Created
			item.ModifyTime = blog.Modified
			item.Content = blog.BlogRelease

			// 获取作者信息
			if member, ok := memberMap[blog.MemberId]; ok {
				if member.RealName != "" {
					item.Author = member.RealName
				} else {
					item.Author = member.Account
				}
			}

			// 提取描述
			description := blog.BlogExcerpt
			if description == "" {
				description = blog.BlogRelease
				if description == "" {
					description = blog.BlogContent
				}
			}
			description = utils.StripTags(description)
			if len([]rune(description)) > 100 {
				description = string([]rune(description)[:100]) + "..."
			}
			item.Description = description

			// 标题匹配加权：搜索词命中标题时提升分数
			titleLower := strings.ToLower(blog.BlogTitle)
			for _, w := range words {
				if strings.Contains(titleLower, w) {
					item.Score *= 1.5
				}
			}

			// 精确匹配加权：博客内容包含原始关键词时大幅提分
			if lowerKeyword != "" {
				contentLower := strings.ToLower(utils.StripTags(blog.BlogRelease))
				if contentLower == "" {
					contentLower = strings.ToLower(utils.StripTags(blog.BlogContent))
				}
				if strings.Contains(contentLower, lowerKeyword) {
					item.Score *= 5.0
				}
			}

			searchResults = append(searchResults, item)
		}
	}

	// 按加权后的分数重新排序
	sort.Slice(searchResults, func(i, j int) bool {
		return searchResults[i].Score > searchResults[j].Score
	})

	// 分页
	totalCount := len(searchResults)
	offset := (pageIndex - 1) * pageSize
	end := offset + pageSize
	if offset > totalCount {
		offset = totalCount
	}
	if end > totalCount {
		end = totalCount
	}
	if offset >= end {
		return nil, words, totalCount, nil
	}

	return searchResults[offset:end], words, totalCount, nil
}

// performSearchV2 执行倒排索引搜索，返回 SearchV2Result 列表
func (c *SearchController) performSearchV2(keyword string, pageIndex, pageSize int) ([]*SearchV2Result, int, error) {
	memberId := 0
	if c.Member != nil {
		memberId = c.Member.MemberId
	}
	rawResults, words, totalCount, err := PerformSearchV2Raw(keyword, pageIndex, pageSize, memberId)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 SearchV2Result
	searchResults := make([]*SearchV2Result, 0, len(rawResults))
	for _, raw := range rawResults {
		item := &SearchV2Result{
			SearchType:   raw.SearchType,
			DocumentId:   raw.DocumentId,
			DocumentName: raw.DocumentName,
			Identify:     raw.Identify,
			Description:  raw.Description,
			Author:       raw.Author,
			ModifyTime:   raw.ModifyTime,
			CreateTime:   raw.CreateTime,
			BookId:       raw.BookId,
			BookName:     raw.BookName,
			BookIdentify: raw.BookIdentify,
		}
		searchResults = append(searchResults, item)
	}

	// 高亮关键词
	for _, item := range searchResults {
		for _, word := range words {
			if word != "" {
				item.DocumentName = strings.Replace(item.DocumentName, word, "<em>"+word+"</em>", -1)
				if item.Description != "" {
					item.Description = strings.Replace(item.Description, word, "<em>"+word+"</em>", -1)
				}
			}
		}
	}

	return searchResults, totalCount, nil
}

type SearchController struct {
	BaseController
}

// 搜索首页
func (c *SearchController) Index() {
	c.Prepare()
	c.TplName = "search/index.tpl"

	//如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		return
	}

	keyword := c.GetString("keyword")
	pageIndex, _ := c.GetInt("page", 1)

	c.Data["BaseUrl"] = c.BaseUrl()

	if keyword != "" {
		c.Data["Keyword"] = keyword
		memberId := 0
		if c.Member != nil {
			memberId = c.Member.MemberId
		}
		searchResult, totalCount, err := models.NewDocumentSearchResult().FindToPager(sqltil.EscapeLike(keyword), pageIndex, conf.PageSize, memberId)

		if err != nil {
			logs.Error("搜索失败 ->", err)
			return
		}
		if totalCount > 0 {
			pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
			c.Data["PageHtml"] = pager.HtmlPages()
		} else {
			c.Data["PageHtml"] = ""
		}
		if len(searchResult) > 0 {
			keywords := strings.Split(keyword, " ")

			for _, item := range searchResult {
				for _, word := range keywords {
					item.DocumentName = strings.Replace(item.DocumentName, word, "<em>"+word+"</em>", -1)
					if item.Description != "" {
						src := item.Description

						r := []rune(utils.StripTags(item.Description))

						if len(r) > 100 {
							src = string(r[:100])
						} else {
							src = string(r)
						}
						item.Description = strings.Replace(src, word, "<em>"+word+"</em>", -1)
					}
				}
				if item.Identify == "" {
					item.Identify = strconv.Itoa(item.DocumentId)
				}
				if item.ModifyTime.IsZero() {
					item.ModifyTime = item.CreateTime
				}
			}
		}
		c.Data["Lists"] = searchResult
	}
}

// 搜索用户
func (c *SearchController) User() {
	c.Prepare()
	key := c.Ctx.Input.Param(":key")
	keyword := strings.TrimSpace(c.GetString("q"))
	if key == "" || keyword == "" {
		c.JsonResult(404, i18n.Tr(c.Lang, "message.param_error"))
	}
	keyword = sqltil.EscapeLike(keyword)

	book, err := models.NewBookResult().FindByIdentify(key, c.Member.MemberId)
	if err != nil {
		if err == models.ErrPermissionDenied {
			c.JsonResult(403, i18n.Tr(c.Lang, "message.no_permission"))
		}
		c.JsonResult(500, i18n.Tr(c.Lang, "message.item_not_exist"))
	}

	//members, err := models.NewMemberRelationshipResult().FindNotJoinUsersByAccount(book.BookId, 10, "%"+keyword+"%")
	members, err := models.NewMemberRelationshipResult().FindNotJoinUsersByAccountOrRealName(book.BookId, 10, "%"+keyword+"%")
	if err != nil {
		logs.Error("查询用户列表出错：" + err.Error())
		c.JsonResult(500, err.Error())
	}
	result := models.SelectMemberResult{}
	items := make([]models.KeyValueItem, 0)

	for _, member := range members {
		item := models.KeyValueItem{}
		item.Id = member.MemberId
		item.Text = member.Account + "[" + member.RealName + "]"
		items = append(items, item)
	}

	result.Result = items

	c.JsonResult(0, "OK", result)
}

// IndexV2 使用倒排索引的搜索页面
func (c *SearchController) IndexV2() {
	c.Prepare()
	c.TplName = "search/index.tpl"

	// 如果没有开启你们访问则跳转到登录
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(conf.URLFor("AccountController.Login"), 302)
		return
	}

	keyword := strings.TrimSpace(c.GetString("keyword"))
	pageIndex, _ := c.GetInt("page", 1)
	pageSize := conf.PageSize

	c.Data["BaseUrl"] = c.BaseUrl()

	if keyword != "" {
		c.Data["Keyword"] = keyword

		searchResult, totalCount, err := c.performSearchV2(keyword, pageIndex, pageSize)

		if err != nil {
			logs.Error("搜索失败 ->", err)
			return
		}
		if totalCount > 0 {
			pager := pagination.NewPagination(c.Ctx.Request, totalCount, conf.PageSize, c.BaseUrl())
			c.Data["PageHtml"] = pager.HtmlPages()
		} else {
			c.Data["PageHtml"] = ""
		}

		// 处理结果中的一些字段
		for _, item := range searchResult {
			if item.Identify == "" {
				item.Identify = strconv.Itoa(item.DocumentId)
			}
		}

		c.Data["Lists"] = searchResult
	}
}

// SearchV2 使用倒排索引进行搜索
func (c *SearchController) SearchV2() {
	c.Prepare()
	memberId := 0
	if c.Member != nil {
		memberId = c.Member.MemberId
	}

	// 如果没有开启匿名访问且未登录则返回错误
	if !c.EnableAnonymous && c.Member == nil {
		c.JsonResult(401, "请先登录")
		return
	}

	keyword := strings.TrimSpace(c.GetString("keyword"))
	if keyword == "" {
		c.JsonResult(400, "搜索关键词不能为空")
		return
	}

	pageIndex, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)

	// 使用底层搜索函数
	rawResults, words, totalCount, err := PerformSearchV2Raw(keyword, pageIndex, pageSize, memberId)
	if err != nil {
		logs.Error("倒排索引搜索失败 ->", err)
		c.JsonResult(500, "搜索失败")
		return
	}

	// 构建返回结果
	searchResults := make([]map[string]interface{}, 0)
	for _, raw := range rawResults {
		item := make(map[string]interface{})
		item["content_type"] = raw.ContentType
		item["content_id"] = raw.ContentId
		item["score"] = raw.Score
		item["word_counts"] = raw.WordCounts
		item["content"] = raw.Content

		if raw.ContentType == 1 {
			// Document类型
			item["type"] = "document"
			item["document_id"] = raw.DocumentId
			item["document_name"] = raw.DocumentName
			item["book_id"] = raw.BookId
			item["book_name"] = raw.BookName
			item["identify"] = raw.Identify
			item["book_identify"] = raw.BookIdentify
			item["create_time"] = raw.CreateTime
			item["modify_time"] = raw.ModifyTime
			item["description"] = raw.Description
		} else if raw.ContentType == 2 {
			// Blog类型
			item["type"] = "blog"
			item["blog_id"] = raw.BlogId
			item["blog_title"] = raw.BlogTitle
			item["blog_identify"] = raw.BlogIdentify
			item["blog_excerpt"] = raw.BlogExcerpt
			item["create_time"] = raw.CreateTime
			item["modify_time"] = raw.ModifyTime
		}
		searchResults = append(searchResults, item)
	}

	// 高亮关键词
	for _, item := range searchResults {
		for _, word := range words {
			if word != "" {
				if title, ok := item["document_name"]; ok {
					item["document_name"] = strings.Replace(title.(string), word, "<em>"+word+"</em>", -1)
				}
				if title, ok := item["blog_title"]; ok {
					item["blog_title"] = strings.Replace(title.(string), word, "<em>"+word+"</em>", -1)
				}
				if desc, ok := item["description"]; ok {
					item["description"] = strings.Replace(desc.(string), word, "<em>"+word+"</em>", -1)
				}
				if excerpt, ok := item["blog_excerpt"]; ok {
					item["blog_excerpt"] = strings.Replace(excerpt.(string), word, "<em>"+word+"</em>", -1)
				}
			}
		}
	}

	responseData := map[string]interface{}{
		"lists":     searchResults,
		"total":     totalCount,
		"page":      pageIndex,
		"page_size": pageSize,
	}

	c.JsonResult(0, "OK", responseData)
}
