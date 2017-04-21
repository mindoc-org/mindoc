package types

// 这些常数定义了反向索引表存储的数据类型
const (
	// 仅存储文档的docId
	DocIdsIndex = 0

	// 存储关键词的词频，用于计算BM25
	FrequenciesIndex = 1

	// 存储关键词在文档中出现的具体字节位置（可能有多个）
	// 如果你希望得到关键词紧邻度数据，必须使用LocationsIndex类型的索引
	LocationsIndex = 2

	// 默认插入索引表文档 CACHE SIZE
	defaultDocCacheSize = 300000
)

// 初始化索引器选项
type IndexerInitOptions struct {
	// 索引表的类型，见上面的常数
	IndexType int

	// 待插入索引表文档 CACHE SIZE
	DocCacheSize int

	// BM25参数
	BM25Parameters *BM25Parameters
}

// 见http://en.wikipedia.org/wiki/Okapi_BM25
// 默认值见engine_init_options.go
type BM25Parameters struct {
	K1 float32
	B  float32
}

func (options *IndexerInitOptions) Init() {
	if options.DocCacheSize == 0 {
		options.DocCacheSize = defaultDocCacheSize
	}
}
