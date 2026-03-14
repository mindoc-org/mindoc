package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/segmenter"
)

func init() {
	//go InitializeMissingIndexes()
}

// ContentReverseIndex 倒排索引结构
type ContentReverseIndex struct {
	Id string `orm:"pk;column(id);size(64);description(唯一标识ID)" json:"id"`
	// Word 分词词汇，最长64个字
	Word string `orm:"column(word);size(64);index;description(分词词汇)" json:"word"`
	// ContentType 内容类型：1-Document 2-Blog
	ContentType int `orm:"column(content_type);type(int);index:idx_content_type_id,priority:1;description(内容类型：1-Document 2-Blog)" json:"content_type"`
	// ContentId 内容ID，对应DocumentId或BlogId
	ContentId int `orm:"column(content_id);type(int);index:idx_content_type_id,priority:2;description(内容ID)" json:"content_id"`
	// WordCount 词频数
	WordCount int `orm:"column(word_count);type(int);default(0);description(词频数)" json:"word_count"`
}

// TableName 获取对应数据库表名
func (c *ContentReverseIndex) TableName() string {
	return "t_content_reverse_index"
}

// TableEngine 获取数据使用的引擎
func (c *ContentReverseIndex) TableEngine() string {
	return "INNODB"
}

func (c *ContentReverseIndex) TableNameWithPrefix() string {
	return conf.GetDatabasePrefix() + c.TableName()
}

func NewContentReverseIndex() *ContentReverseIndex {
	return &ContentReverseIndex{}
}

// Insert 插入倒排索引记录
func (c *ContentReverseIndex) Insert() error {
	if c.Id == "" {
		return errors.New("id不能为空")
	}
	if c.Word == "" {
		return errors.New("分词词汇不能为空")
	}
	if c.ContentType != 1 && c.ContentType != 2 {
		return errors.New("内容类型必须是1(Document)或2(Blog)")
	}
	if c.ContentId <= 0 {
		return errors.New("内容ID必须大于0")
	}

	o := orm.NewOrm()
	_, err := o.Insert(c)
	return err
}

// DeleteByContentTypeAndContentId 根据内容类型和内容ID删除所有倒排索引记录
func (c *ContentReverseIndex) DeleteByContentTypeAndContentId(contentType, contentId int) error {
	if contentType != 1 && contentType != 2 {
		return errors.New("内容类型必须是1(Document)或2(Blog)")
	}
	if contentId <= 0 {
		return errors.New("内容ID必须大于0")
	}

	o := orm.NewOrm()
	_, err := o.QueryTable(c.TableNameWithPrefix()).Filter("content_type", contentType).Filter("content_id", contentId).Delete()
	return err
}

// BatchInsert 批量插入倒排索引记录
func (c *ContentReverseIndex) BatchInsert(indices []*ContentReverseIndex) error {
	if len(indices) == 0 {
		return nil
	}

	o := orm.NewOrm()
	_, err := o.InsertMulti(len(indices), indices)
	return err
}

// ContentReverseIndexResult 倒排索引查询结果结构
type ContentReverseIndexResult struct {
	ContentId   int     `json:"content_id"`
	ContentType int     `json:"content_type"`
	Score       float64 `json:"score"`       // TF-IDF分数
	WordCounts  []int   `json:"word_counts"` // 各个词的词频
}

// FindByWordsWithPagination 根据多个分词词汇分页批量查询结果，按IDF值排序
// words: 分词词汇列表
// pageIndex: 页码，从1开始
// pageSize: 每页数量
func (c *ContentReverseIndex) FindByWordsWithPagination(words []string, pageIndex, pageSize int) ([]*ContentReverseIndexResult, int, error) {
	if len(words) == 0 {
		return nil, 0, errors.New("分词词汇列表不能为空")
	}
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	o := orm.NewOrm()
	tableName := c.TableNameWithPrefix()

	// 计算总文档数
	totalDocsSql := "SELECT COUNT(DISTINCT CONCAT(content_type, '-', content_id)) FROM " + tableName
	var totalDocs int
	err := o.Raw(totalDocsSql).QueryRow(&totalDocs)
	if err != nil {
		return nil, 0, err
	}

	// 构建IN条件
	wordPlaceholders := ""
	wordArgs := make([]any, 0)
	for i, word := range words {
		if i > 0 {
			wordPlaceholders += ","
		}
		wordPlaceholders += "?"
		wordArgs = append(wordArgs, word)
	}
	sql := "SELECT word, content_type, content_id, word_count FROM " + tableName +
		" WHERE word IN (" + wordPlaceholders + ") ORDER BY content_type, content_id"
	type indexRecord struct {
		Word        string
		ContentType int
		ContentId   int
		WordCount   int
	}
	var records []indexRecord
	_, err = o.Raw(sql, wordArgs...).QueryRows(&records)
	if err != nil {
		return nil, 0, err
	}

	// 计算各文档的总词数
	sql = "SELECT content_type, content_id, count(word_count) total_word_count FROM " + tableName +
		" GROUP BY content_type, content_id"
	type docWordCountRecord struct {
		ContentType    int
		ContentId      int
		TotalWordCount int
	}
	var docWordCountRecords []docWordCountRecord
	_, err = o.Raw(sql).QueryRows(&docWordCountRecords)
	if err != nil {
		return nil, 0, err
	}

	docTotalWordCountMap := make(map[string]int)
	for _, record := range docWordCountRecords {
		key := fmt.Sprintf("%d-%d", record.ContentType, record.ContentId)
		docTotalWordCountMap[key] = record.TotalWordCount
	}

	// 聚合每个(content_type, content_id)的词频和计算TF-IDF
	contentMap := make(map[string]*ContentReverseIndexResult)
	for _, record := range records {
		key := fmt.Sprintf("%d-%d", record.ContentType, record.ContentId)
		if result, exists := contentMap[key]; exists {
			result.WordCounts = append(result.WordCounts, record.WordCount)
		} else {
			contentMap[key] = &ContentReverseIndexResult{
				ContentId:   record.ContentId,
				ContentType: record.ContentType,
				WordCounts:  []int{record.WordCount},
			}
		}
	}

	docMapWithWords := make(map[string]int) // 用于计算包含搜索词的文档数
	// 计算每个文档包含多少个查询词
	docWordCount := make(map[string]int)
	for _, record := range records {
		key := fmt.Sprintf("%d-%d", record.ContentType, record.ContentId)
		docWordCount[key] += record.WordCount
		docMapWithWords[key] += 1
	}

	// 计算IDF并生成结果
	results := make([]*ContentReverseIndexResult, 0, len(contentMap))
	for key := range contentMap {
		result := contentMap[key]
		// 计算TF：词频之和
		tf := float64(docWordCount[key]) / float64(docTotalWordCountMap[key]+1)
		// 计算DF：包含该词的文档数（简化处理，使用该文档包含的查询词数量）
		df := len(docMapWithWords)
		// 计算IDF
		idf := 0.0
		if df > 0 && totalDocs > 0 {
			idf = math.Log(float64(totalDocs+1) / float64(df))
		}

		// 用于根据文档总词数调整TF-IDF的权重，避免总词数过小的文档权重过高
		alpha := math.Log(1.0+float64(docTotalWordCountMap[key])*0.01) * 100
		// TF-IDF分数
		result.Score = float64(tf) * idf * float64(alpha)
		results = append(results, result)
	}

	// 按Score降序排序
	sortResultsByScore(results)
	totalCount := len(results)
	// 分页
	offset := (pageIndex - 1) * pageSize
	start := offset
	end := offset + pageSize
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}
	if start >= end {
		return nil, totalCount, nil
	}

	return results[start:end], totalCount, nil
}

func sortResultsByScore(results []*ContentReverseIndexResult) {
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

func generateIndexId(contentType, contentId int, word string) string {
	source := fmt.Sprintf("%d-%d-%s", contentType, contentId, word)
	hasher := md5.New()
	hasher.Write([]byte(source))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)[:32]
}

func BuildIndexForDocument(documentId int, content string) error {
	if documentId <= 0 {
		return errors.New("文档ID必须大于0")
	}

	index := NewContentReverseIndex()
	err := index.DeleteByContentTypeAndContentId(1, documentId)
	if err != nil {
		logs.Error("删除文档倒排索引失败 ->", documentId, err)
		return err
	}

	words := segmenter.Segment(content)
	if len(words) == 0 {
		return nil
	}

	wordCountMap := make(map[string]int)
	for _, word := range words {
		if len(word) > 64 {
			word = word[:64]
		}
		wordCountMap[word]++
	}

	indices := make([]*ContentReverseIndex, 0, len(wordCountMap))
	for word, count := range wordCountMap {
		id := generateIndexId(1, documentId, word)

		indexItem := &ContentReverseIndex{
			Id:          id,
			Word:        word,
			ContentType: 1,
			ContentId:   documentId,
			WordCount:   count,
		}
		indices = append(indices, indexItem)
	}

	if len(indices) > 0 {
		err = index.BatchInsert(indices)
		if err != nil {
			return fmt.Errorf("批量插入文档倒排索引失败 -> %d %v", documentId, err)
		}
	}

	return nil
}

func BuildIndexForBlog(blogId int, content string) error {
	if blogId <= 0 {
		return errors.New("BlogID必须大于0")
	}

	index := NewContentReverseIndex()
	err := index.DeleteByContentTypeAndContentId(2, blogId)
	if err != nil {
		logs.Error("删除Blog倒排索引失败 ->", blogId, err)
		return err
	}

	words := segmenter.Segment(content)
	if len(words) == 0 {
		return nil
	}

	wordCountMap := make(map[string]int)
	for _, word := range words {
		if len(word) > 64 {
			word = word[:64]
		}
		wordCountMap[word]++
	}

	indices := make([]*ContentReverseIndex, 0, len(wordCountMap))
	for word, count := range wordCountMap {
		id := generateIndexId(2, blogId, word)

		indexItem := &ContentReverseIndex{
			Id:          id,
			Word:        word,
			ContentType: 2,
			ContentId:   blogId,
			WordCount:   count,
		}
		indices = append(indices, indexItem)
	}

	if len(indices) > 0 {
		err = index.BatchInsert(indices)
		if err != nil {
			logs.Error("批量插入Blog倒排索引失败 ->", blogId, err)
			return err
		}
	}

	return nil
}

func CheckDocumentIndexed(documentId int) bool {
	if documentId <= 0 {
		return false
	}

	o := orm.NewOrm()
	index := NewContentReverseIndex()
	return o.QueryTable(index.TableNameWithPrefix()).Filter("content_type", 1).Filter("content_id", documentId).Exist()
}

func CheckBlogIndexed(blogId int) bool {
	if blogId <= 0 {
		return false
	}

	o := orm.NewOrm()
	index := NewContentReverseIndex()
	return o.QueryTable(index.TableNameWithPrefix()).Filter("content_type", 2).Filter("content_id", blogId).Exist()
}

func GetUnindexedDocuments(limit int) ([]*Document, error) {
	o := orm.NewOrm()
	var documents []*Document
	docTable := NewDocument().TableNameWithPrefix()
	indexTable := NewContentReverseIndex().TableNameWithPrefix()

	sql := "SELECT d.* FROM " + docTable + " d " +
		"LEFT JOIN " + indexTable + " i ON i.content_type = 1 AND i.content_id = d.document_id " +
		"WHERE i.id IS NULL " +
		"ORDER BY d.document_id DESC"

	if limit > 0 {
		sql += " LIMIT ?"
		_, err := o.Raw(sql, limit).QueryRows(&documents)
		return documents, err
	}

	_, err := o.Raw(sql).QueryRows(&documents)
	return documents, err
}

func GetUnindexedBlogs(limit int) ([]*Blog, error) {
	o := orm.NewOrm()
	var blogs []*Blog
	blogTable := NewBlog().TableNameWithPrefix()
	indexTable := NewContentReverseIndex().TableNameWithPrefix()

	sql := "SELECT b.* FROM " + blogTable + " b " +
		"LEFT JOIN " + indexTable + " i ON i.content_type = 2 AND i.content_id = b.blog_id " +
		"WHERE i.id IS NULL " +
		"ORDER BY b.blog_id DESC"

	if limit > 0 {
		sql += " LIMIT ?"
		_, err := o.Raw(sql, limit).QueryRows(&blogs)
		return blogs, err
	}

	_, err := o.Raw(sql).QueryRows(&blogs)
	return blogs, err
}

// InitializeMissingIndexes 初始化缺失的倒排索引
func InitializeMissingIndexes() {
	go func() {
		logs.Info("开始检查并初始化缺失的倒排索引...")
		InitializeMissingDocumentIndexes()
		InitializeMissingBlogIndexes()
		logs.Info("倒排索引初始化检查完成")
	}()
}

func InitializeMissingDocumentIndexes() {
	batchSize := 100
	for {
		documents, err := GetUnindexedDocuments(batchSize)
		if err != nil {
			logs.Error("获取未索引文档失败 ->", err)
			break
		}

		if len(documents) == 0 {
			break
		}

		for _, doc := range documents {
			indexed := CheckDocumentIndexed(doc.DocumentId)
			if !indexed {
				content := doc.Release
				if content == "" {
					content = doc.Markdown
				}
				for i := 0; i < 10; i++ { // 标题内容"十分"重要
					content = doc.DocumentName + "\n" + content
				}
				content = utils.StripTags(content)
				err := BuildIndexForDocument(doc.DocumentId, content)
				if err != nil {
					logs.Error("构建文档倒排索引失败 ->", doc.DocumentId, err)
				} else {
					logs.Info("文档倒排索引构建成功 ->", doc.DocumentId)
				}
			}
		}
	}
}

func InitializeMissingBlogIndexes() {
	batchSize := 100
	for {
		blogs, err := GetUnindexedBlogs(batchSize)
		if err != nil {
			logs.Error("获取未索引Blog失败 ->", err)
			break
		}

		if len(blogs) == 0 {
			break
		}

		for _, blog := range blogs {
			indexed := CheckBlogIndexed(blog.BlogId)
			if !indexed {
				content := blog.BlogRelease
				if content == "" {
					content = blog.BlogContent
				}
				for i := 0; i < 10; i++ { // 标题内容"十分"重要
					content = blog.BlogTitle + "\n" + content
				}
				content = utils.StripTags(content)

				err := BuildIndexForBlog(blog.BlogId, content)
				if err != nil {
					logs.Error("构建Blog倒排索引失败 ->", blog.BlogId, err)
				} else {
					logs.Info("Blog倒排索引构建成功 ->", blog.BlogId)
				}
			}
		}
	}
}
