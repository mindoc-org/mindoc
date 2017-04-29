package types

import (
	"log"
	"runtime"
)

var (
	// EngineInitOptions的默认值
	defaultNumSegmenterThreads       = runtime.NumCPU()
	defaultNumShards                 = 2
	defaultIndexerBufferLength       = runtime.NumCPU()
	defaultNumIndexerThreadsPerShard = runtime.NumCPU()
	defaultRankerBufferLength        = runtime.NumCPU()
	defaultNumRankerThreadsPerShard  = runtime.NumCPU()
	defaultDefaultRankOptions        = RankOptions{
		ScoringCriteria: RankByBM25{},
	}
	defaultIndexerInitOptions = IndexerInitOptions{
		IndexType:      FrequenciesIndex,
		BM25Parameters: &defaultBM25Parameters,
	}
	defaultBM25Parameters = BM25Parameters{
		K1: 2.0,
		B:  0.75,
	}
	defaultPersistentStorageShards = 8
)

type EngineInitOptions struct {
	// 是否使用分词器
	// 默认使用，否则在启动阶段跳过SegmenterDictionaries和StopTokenFile设置
	// 如果你不需要在引擎内分词，可以将这个选项设为true
	// 注意，如果你不用分词器，那么在调用IndexDocument时DocumentIndexData中的Content会被忽略
	NotUsingSegmenter bool

	// 半角逗号分隔的字典文件，具体用法见
	// sego.Segmenter.LoadDictionary函数的注释
	SegmenterDictionaries string

	// 停用词文件
	StopTokenFile string

	// 分词器线程数
	NumSegmenterThreads int

	// 索引器和排序器的shard数目
	// 被检索/排序的文档会被均匀分配到各个shard中
	NumShards int

	// 索引器的信道缓冲长度
	IndexerBufferLength int

	// 索引器每个shard分配的线程数
	NumIndexerThreadsPerShard int

	// 排序器的信道缓冲长度
	RankerBufferLength int

	// 排序器每个shard分配的线程数
	NumRankerThreadsPerShard int

	// 索引器初始化选项
	IndexerInitOptions *IndexerInitOptions

	// 默认的搜索选项
	DefaultRankOptions *RankOptions

	// 是否使用持久数据库，以及数据库文件保存的目录和裂分数目
	UsePersistentStorage    bool
	PersistentStorageFolder string
	PersistentStorageShards int
}

// 初始化EngineInitOptions，当用户未设定某个选项的值时用默认值取代
func (options *EngineInitOptions) Init() {
	if !options.NotUsingSegmenter {
		if options.SegmenterDictionaries == "" {
			log.Fatal("字典文件不能为空")
		}
	}

	if options.NumSegmenterThreads == 0 {
		options.NumSegmenterThreads = defaultNumSegmenterThreads
	}

	if options.NumShards == 0 {
		options.NumShards = defaultNumShards
	}

	if options.IndexerBufferLength == 0 {
		options.IndexerBufferLength = defaultIndexerBufferLength
	}

	if options.NumIndexerThreadsPerShard == 0 {
		options.NumIndexerThreadsPerShard = defaultNumIndexerThreadsPerShard
	}

	if options.RankerBufferLength == 0 {
		options.RankerBufferLength = defaultRankerBufferLength
	}

	if options.NumRankerThreadsPerShard == 0 {
		options.NumRankerThreadsPerShard = defaultNumRankerThreadsPerShard
	}

	if options.IndexerInitOptions == nil {
		options.IndexerInitOptions = &defaultIndexerInitOptions
	}

	if options.IndexerInitOptions.BM25Parameters == nil {
		options.IndexerInitOptions.BM25Parameters = &defaultBM25Parameters
	}

	if options.DefaultRankOptions == nil {
		options.DefaultRankOptions = &defaultDefaultRankOptions
	}

	if options.DefaultRankOptions.ScoringCriteria == nil {
		options.DefaultRankOptions.ScoringCriteria = defaultDefaultRankOptions.ScoringCriteria
	}

	if options.PersistentStorageShards == 0 {
		options.PersistentStorageShards = defaultPersistentStorageShards
	}
}
