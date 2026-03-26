package segmenter

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
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
	// 停用词集合
	stopWords map[string]bool
	technicalTermPattern = regexp.MustCompile(`(?i)[a-z0-9][a-z0-9+#._/-]{1,63}`)
)

// techTermWhitelist 技术术语白名单
// 这些词虽然是常见英语词汇，但同时也是 Linux/Unix 命令、编程语言
// 或重要技术术语，不应被停用词过滤，否则用户搜索相关命令时将无法找到文档
var techTermWhitelist = map[string]bool{
	// Linux/Unix 常用命令（同时也是英语常见词）
	"find": true, "top": true, "last": true, "more": true, "less": true,
	"who": true, "which": true, "done": true, "move": true, "give": true,
	"make": true, "take": true, "fill": true, "split": true, "cut": true,
	// 编程语言/框架名称
	"go": true, "net": true, "next": true,
	// HTTP 方法 / 数据库操作
	"get": true, "put": true, "call": true, "show": true, "describe": true,
	"like": true,
	// 系统/运维/容器/网络相关
	"system": true, "volume": true, "name": true, "save": true, "keep": true,
	"re": true, "mine": true, "near": true, "fire": true, "front": true,
	"full": true, "empty": true, "computer": true, "detail": true, "part": true,
	"back": true, "down": true, "up": true, "bar": true, "round": true,
	"side": true, "bottom": true,
	// 工具/软件名称（同时也是英语单词）
	"everything": true,
}

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
		// 加载停用词表
		stopWords = loadStopWords(stopWordsDict)
		logs.Info("jieba分词器初始化完成, 停用词数:", len(stopWords))
	})
}

// loadStopWords 从文件加载停用词集合
func loadStopWords(filePath string) map[string]bool {
	sw := make(map[string]bool)
	f, err := os.Open(filePath)
	if err != nil {
		logs.Error("加载停用词表失败 ->", err)
		return sw
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			sw[strings.ToLower(word)] = true
		}
	}
	return sw
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
		// 过滤停用词（白名单中的技术术语不过滤）
		if stopWords[word] && !techTermWhitelist[word] {
			continue
		}
		result = append(result, word)
	}

	for _, word := range extractTechnicalTerms(text) {
		if stopWords[word] && !techTermWhitelist[word] {
			continue
		}
		result = append(result, word)
	}

	return result
}

func extractTechnicalTerms(text string) []string {
	matches := technicalTermPattern.FindAllString(text, -1)
	if len(matches) == 0 {
		return nil
	}
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		word := strings.ToLower(strings.TrimSpace(match))
		if word == "" {
			continue
		}
		if len([]rune(word)) < 2 {
			continue
		}
		if !strings.ContainsAny(word, ".+#/_-") {
			continue
		}
		hasAlphaNumeric := false
		for _, r := range word {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				hasAlphaNumeric = true
				break
			}
		}
		if !hasAlphaNumeric {
			continue
		}
		result = append(result, word)
	}
	return result
}
