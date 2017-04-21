package types

import (
	"github.com/huichen/wukong/utils"
)

type SearchResponse struct {
	// 搜索用到的关键词
	Tokens []string

	// 搜索到的文档，已排序
	Docs []ScoredDocument

	// 搜索是否超时。超时的情况下也可能会返回部分结果
	Timeout bool

	// 搜索到的文档个数。注意这是全部文档中满足条件的个数，可能比返回的文档数要大
	NumDocs int
}

type ScoredDocument struct {
	DocId uint64

	// 文档的打分值
	// 搜索结果按照Scores的值排序，先按照第一个数排，如果相同则按照第二个数排序，依次类推。
	Scores []float32

	// 用于生成摘要的关键词在文本中的字节位置，该切片长度和SearchResponse.Tokens的长度一样
	// 只有当IndexType == LocationsIndex时不为空
	TokenSnippetLocations []int

	// 关键词出现的位置
	// 只有当IndexType == LocationsIndex时不为空
	TokenLocations [][]int
}

// 为了方便排序

type ScoredDocuments []ScoredDocument

func (docs ScoredDocuments) Len() int {
	return len(docs)
}
func (docs ScoredDocuments) Swap(i, j int) {
	docs[i], docs[j] = docs[j], docs[i]
}
func (docs ScoredDocuments) Less(i, j int) bool {
	// 为了从大到小排序，这实际上实现的是More的功能
	for iScore := 0; iScore < utils.MinInt(len(docs[i].Scores), len(docs[j].Scores)); iScore++ {
		if docs[i].Scores[iScore] > docs[j].Scores[iScore] {
			return true
		} else if docs[i].Scores[iScore] < docs[j].Scores[iScore] {
			return false
		}
	}
	return len(docs[i].Scores) > len(docs[j].Scores)
}
