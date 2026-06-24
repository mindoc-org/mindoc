package models

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/utils"
	"github.com/mindoc-org/mindoc/utils/segmenter"
)

func init() {
	//go InitializeMissingIndexes()
}

const emptyIndexWord = "__mindoc_empty_index__"

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

// FindByWords 根据多个分词词汇查询结果，按TF-IDF值排序，返回全部匹配结果（不分页）
// words: 分词词汇列表
// 返回值: 结果列表, 匹配的总文档数(截断前), error
func (c *ContentReverseIndex) FindByWords(words []string) ([]*ContentReverseIndexResult, int, error) {
	if len(words) == 0 {
		return nil, 0, errors.New("分词词汇列表不能为空")
	}
	words = normalizeIndexWords(words)
	if len(words) == 0 {
		return nil, 0, errors.New("分词词汇列表不能为空")
	}
	// 限制查询词数量，防止恶意超长关键词生成巨大IN子句
	const maxWords = 50
	if len(words) > maxWords {
		words = words[:maxWords]
	}

	o := orm.NewOrm()
	tableName := c.TableNameWithPrefix()
	if !validTableName.MatchString(tableName) {
		return nil, 0, errors.New("非法表名: " + tableName)
	}

	// 计算总文档数
	totalDocsSql := "SELECT COUNT(*) FROM (SELECT DISTINCT content_type, content_id FROM " + tableName + ") AS t"
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

	// 计算每个词的文档频率（DF）：每个词出现在多少个文档中
	wordDocFreq := make(map[string]map[string]bool)
	for _, record := range records {
		key := fmt.Sprintf("%d-%d", record.ContentType, record.ContentId)
		if wordDocFreq[record.Word] == nil {
			wordDocFreq[record.Word] = make(map[string]bool)
		}
		wordDocFreq[record.Word][key] = true
	}

	// 聚合每个文档的匹配词信息
	type docWordInfo struct {
		Word      string
		WordCount int
	}
	docWords := make(map[string][]docWordInfo)
	for _, record := range records {
		key := fmt.Sprintf("%d-%d", record.ContentType, record.ContentId)
		docWords[key] = append(docWords[key], docWordInfo{
			Word:      record.Word,
			WordCount: record.WordCount,
		})
	}

	// 计算每个文档的TF-IDF分数（使用正确的per-word IDF）
	results := make([]*ContentReverseIndexResult, 0, len(docWords))
	for key, wordInfos := range docWords {
		var contentType, contentId int
		if _, err := fmt.Sscanf(key, "%d-%d", &contentType, &contentId); err != nil {
			logs.Error("解析文档key失败 ->", key, err)
			continue
		}

		score := 0.0
		wordCounts := make([]int, 0, len(wordInfos))

		for _, wi := range wordInfos {
			wordCounts = append(wordCounts, wi.WordCount)
			// TF: 使用对数TF（sublinear TF），避免长文档被不合理惩罚
			tf := 1.0 + math.Log(float64(wi.WordCount)+1)
			// IDF: 每个词独立计算，稀有词权重更高
			df := len(wordDocFreq[wi.Word])
			idf := 0.0
			if df > 0 && totalDocs > 0 {
				idf = math.Log(float64(totalDocs+1) / float64(df+1))
			}
			// 词长权重：长词（更具体的词）贡献更大
			wordLen := float64(len([]rune(wi.Word)))
			lengthWeight := math.Log2(1.0 + wordLen)
			score += tf * idf * lengthWeight
		}

		// 查询词覆盖率加成：匹配的查询词越多，分数越高
		coverage := float64(len(wordInfos)) / float64(len(words))
		score *= (1.0 + coverage)

		results = append(results, &ContentReverseIndexResult{
			ContentId:   contentId,
			ContentType: contentType,
			Score:       score,
			WordCounts:  wordCounts,
		})
	}

	// 按Score降序排序
	sortResultsByScore(results)

	return results, len(results), nil
}

func normalizeIndexWords(words []string) []string {
	result := make([]string, 0, len(words))
	seen := make(map[string]struct{}, len(words))
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" || word == emptyIndexWord {
			continue
		}
		if _, ok := seen[word]; ok {
			continue
		}
		seen[word] = struct{}{}
		result = append(result, word)
	}
	return result
}

func sortResultsByScore(results []*ContentReverseIndexResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		if results[i].ContentType != results[j].ContentType {
			return results[i].ContentType < results[j].ContentType
		}
		return results[i].ContentId < results[j].ContentId
	})
}

func generateIndexId(contentType, contentId int, word string) string {
	source := fmt.Sprintf("%d-%d-%s", contentType, contentId, word)
	hasher := md5.New()
	hasher.Write([]byte(source))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)[:32]
}

func buildEmptyIndexRecord(contentType, contentId int) *ContentReverseIndex {
	return &ContentReverseIndex{
		Id:          generateIndexId(contentType, contentId, emptyIndexWord),
		Word:        emptyIndexWord,
		ContentType: contentType,
		ContentId:   contentId,
		WordCount:   1,
	}
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
	wordCountMap := make(map[string]int)
	for _, word := range words {
		if len(word) > 64 {
			word = word[:64]
		}
		wordCountMap[word]++
	}

	indices := make([]*ContentReverseIndex, 0, len(wordCountMap))
	if len(wordCountMap) == 0 {
		indices = append(indices, buildEmptyIndexRecord(1, documentId))
	}
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
	wordCountMap := make(map[string]int)
	for _, word := range words {
		if len(word) > 64 {
			word = word[:64]
		}
		wordCountMap[word]++
	}

	indices := make([]*ContentReverseIndex, 0, len(wordCountMap))
	if len(wordCountMap) == 0 {
		indices = append(indices, buildEmptyIndexRecord(2, blogId))
	}
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

	if !validTableName.MatchString(docTable) || !validTableName.MatchString(indexTable) {
		return nil, errors.New("非法表名")
	}

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

	if !validTableName.MatchString(blogTable) || !validTableName.MatchString(indexTable) {
		return nil, errors.New("非法表名")
	}

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
				content = doc.DocumentName + "\n" + content
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
				content = blog.BlogTitle + "\n" + content
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

// validTableName 校验表名仅包含安全字符，防止 SQL 注入
var validTableName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// RebuildAllIndexes 全量重建倒排索引（先清空再重建）。
// 注意：MySQL/Postgres 使用 TRUNCATE，SQLite 使用 DELETE，若后续重建阶段发生错误，
// 索引表将处于部分重建状态，此时搜索结果可能不完整。
// 建议在业务低峰期执行，并在返回 error 时手动重新执行本命令。
func RebuildAllIndexes() error {
	logs.Info("开始全量重建倒排索引...")

	// 清空倒排索引表
	o := orm.NewOrm()
	tableName := NewContentReverseIndex().TableNameWithPrefix()
	if !validTableName.MatchString(tableName) {
		err := errors.New("非法表名，拒绝执行: " + tableName)
		logs.Error(err)
		return err
	}
	err := clearReverseIndexTable(o, tableName)
	if err != nil {
		logs.Error("清空倒排索引表失败 ->", err)
		return err
	}
	logs.Info("倒排索引表已清空")

	// 重建文档索引
	if err := rebuildDocumentIndexes(); err != nil {
		logs.Error("文档索引重建失败，索引表处于部分重建状态，请重新执行 reindex ->", err)
		return err
	}
	// 重建博客索引
	if err := rebuildBlogIndexes(); err != nil {
		logs.Error("博客索引重建失败，索引表处于部分重建状态，请重新执行 reindex ->", err)
		return err
	}

	logs.Info("全量重建倒排索引完成")
	return nil
}

func clearReverseIndexTable(o orm.Ormer, tableName string) error {
	dbadapter, _ := web.AppConfig.String("db_adapter")
	if strings.EqualFold(dbadapter, "sqlite3") {
		_, err := o.Raw("DELETE FROM " + tableName).Exec()
		return err
	}
	_, err := o.Raw("TRUNCATE TABLE " + tableName).Exec()
	return err
}

func rebuildDocumentIndexes() error {
	o := orm.NewOrm()
	batchSize := 100
	offset := 0
	total := 0
	failed := 0
	var firstErr error

	for {
		var documents []*Document
		_, err := o.QueryTable(NewDocument().TableNameWithPrefix()).
			OrderBy("document_id").
			Limit(batchSize, offset).
			All(&documents)
		if err != nil {
			logs.Error("查询文档失败 ->", err)
			return err
		}
		if len(documents) == 0 {
			break
		}

		for _, doc := range documents {
			content := doc.Release
			if content == "" {
				content = doc.Markdown
			}
			content = doc.DocumentName + "\n" + content
			content = utils.StripTags(content)
			if err := BuildIndexForDocument(doc.DocumentId, content); err != nil {
				logs.Error("重建文档倒排索引失败 ->", doc.DocumentId, err)
				failed++
				if firstErr == nil {
					firstErr = fmt.Errorf("document_id=%d: %w", doc.DocumentId, err)
				}
			} else {
				total++
			}
		}

		offset += batchSize
		logs.Info("已重建文档索引:", total, "失败:", failed)
	}
	logs.Info("文档索引重建完成, 成功:", total, "失败:", failed)
	if failed > 0 {
		return fmt.Errorf("文档索引重建存在 %d 条失败，首个错误: %w", failed, firstErr)
	}
	return nil
}

func rebuildBlogIndexes() error {
	o := orm.NewOrm()
	batchSize := 100
	offset := 0
	total := 0
	failed := 0
	var firstErr error

	for {
		var blogs []*Blog
		_, err := o.QueryTable(NewBlog().TableNameWithPrefix()).
			OrderBy("blog_id").
			Limit(batchSize, offset).
			All(&blogs)
		if err != nil {
			logs.Error("查询博客失败 ->", err)
			return err
		}
		if len(blogs) == 0 {
			break
		}

		for _, blog := range blogs {
			content := blog.BlogRelease
			if content == "" {
				content = blog.BlogContent
			}
			content = blog.BlogTitle + "\n" + content
			content = utils.StripTags(content)
			if err := BuildIndexForBlog(blog.BlogId, content); err != nil {
				logs.Error("重建Blog倒排索引失败 ->", blog.BlogId, err)
				failed++
				if firstErr == nil {
					firstErr = fmt.Errorf("blog_id=%d: %w", blog.BlogId, err)
				}
			} else {
				total++
			}
		}

		offset += batchSize
		logs.Info("已重建Blog索引:", total, "失败:", failed)
	}
	logs.Info("Blog索引重建完成, 成功:", total, "失败:", failed)
	if failed > 0 {
		return fmt.Errorf("博客索引重建存在 %d 条失败，首个错误: %w", failed, firstErr)
	}
	return nil
}
