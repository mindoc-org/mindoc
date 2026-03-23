package segmenter

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/beego/beego/v2/core/logs"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/yanyiwu/gojieba"
)

var (
	// jieba 分词器实例
	segmenterOnce sync.Once
	jiebaCut      *gojieba.Jieba
)

// getDictDir 获取词典目录
func getDictDir() string {
	// 使用项目根目录下的 lib/jieba 目录
	return filepath.Join(conf.WorkingDirectory, "lib", "jieba")
}

// ensureDictDir 确保词典目录存在
func ensureDictDir() error {
	dictDir := getDictDir()
	if _, err := os.Stat(dictDir); os.IsNotExist(err) {
		return os.MkdirAll(dictDir, 0755)
	}
	return nil
}

// initJieba 初始化 jieba 分词器
func initJieba() {
	segmenterOnce.Do(func() {
		dictDir := getDictDir()

		// 显式传入词典路径
		jiebaDict := filepath.Join(dictDir, "jieba.dict.utf8")
		hmmDict := filepath.Join(dictDir, "hmm_model.utf8")
		userDict := filepath.Join(dictDir, "user.dict.utf8")
		idfDict := filepath.Join(dictDir, "idf.utf8")
		stopWordsDict := filepath.Join(dictDir, "stop_words.utf8")
		// 确保词典目录存在
		if err := ensureDictDir(); err != nil {
			logs.Error("创建词典目录失败 ->", err)
		}
		// 创建分词器
		jiebaCut = gojieba.NewJieba(jiebaDict, hmmDict, userDict, idfDict, stopWordsDict)
		logs.Info("jieba分词器初始化完成")
	})
}

// Segment 中文分词器
// 使用 jieba 分词库的搜索引擎模式进行分词
func Segment(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	// 初始化分词器
	initJieba()
	// 使用 jieba 分词，搜索引擎模式
	// CutForSearch 第二个参数为 true 表示使用 HMM 模型
	words := jiebaCut.CutForSearch(text, true)
	// 过滤结果
	result := make([]string, 0)
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		// 转小写（英文）
		word = strings.ToLower(word)
		// 过滤单字符标点符号/特殊字符，避免匹配大量无关文档
		runes := []rune(word)
		if len(runes) == 1 && !unicode.IsLetter(runes[0]) && !unicode.IsDigit(runes[0]) {
			continue
		}
		result = append(result, word)
	}

	return result
}
