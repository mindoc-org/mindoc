package core

import (
	"log"
	"math"
	"sort"
	"sync"

	"github.com/huichen/wukong/types"
	"github.com/huichen/wukong/utils"
)

// 索引器
type Indexer struct {
	// 从搜索键到文档列表的反向索引
	// 加了读写锁以保证读写安全
	tableLock struct {
		sync.RWMutex
		table     map[string]*KeywordIndices
		docsState map[uint64]int // nil: 表示无状态记录，0: 存在于索引中，1: 等待删除，2: 等待加入
	}
	addCacheLock struct {
		sync.RWMutex
		addCachePointer int
		addCache        types.DocumentsIndex
	}
	removeCacheLock struct {
		sync.RWMutex
		removeCachePointer int
		removeCache        types.DocumentsId
	}

	initOptions types.IndexerInitOptions
	initialized bool

	// 这实际上是总文档数的一个近似
	numDocuments uint64

	// 所有被索引文本的总关键词数
	totalTokenLength float32

	// 每个文档的关键词长度
	docTokenLengths map[uint64]float32
}

// 反向索引表的一行，收集了一个搜索键出现的所有文档，按照DocId从小到大排序。
type KeywordIndices struct {
	// 下面的切片是否为空，取决于初始化时IndexType的值
	docIds      []uint64  // 全部类型都有
	frequencies []float32 // IndexType == FrequenciesIndex
	locations   [][]int   // IndexType == LocationsIndex
}

// 初始化索引器
func (indexer *Indexer) Init(options types.IndexerInitOptions) {
	if indexer.initialized == true {
		log.Fatal("索引器不能初始化两次")
	}
	options.Init()
	indexer.initOptions = options
	indexer.initialized = true

	indexer.tableLock.table = make(map[string]*KeywordIndices)
	indexer.tableLock.docsState = make(map[uint64]int)
	indexer.addCacheLock.addCache = make([]*types.DocumentIndex, indexer.initOptions.DocCacheSize)
	indexer.removeCacheLock.removeCache = make([]uint64, indexer.initOptions.DocCacheSize*2)
	indexer.docTokenLengths = make(map[uint64]float32)
}

// 从KeywordIndices中得到第i个文档的DocId
func (indexer *Indexer) getDocId(ti *KeywordIndices, i int) uint64 {
	return ti.docIds[i]
}

// 得到KeywordIndices中文档总数
func (indexer *Indexer) getIndexLength(ti *KeywordIndices) int {
	return len(ti.docIds)
}

// 向 ADDCACHE 中加入一个文档
func (indexer *Indexer) AddDocumentToCache(document *types.DocumentIndex, forceUpdate bool) {
	if indexer.initialized == false {
		log.Fatal("索引器尚未初始化")
	}

	indexer.addCacheLock.Lock()
	if document != nil {
		indexer.addCacheLock.addCache[indexer.addCacheLock.addCachePointer] = document
		indexer.addCacheLock.addCachePointer++
	}
	if indexer.addCacheLock.addCachePointer >= indexer.initOptions.DocCacheSize || forceUpdate {
		indexer.tableLock.Lock()
		position := 0
		for i := 0; i < indexer.addCacheLock.addCachePointer; i++ {
			docIndex := indexer.addCacheLock.addCache[i]
			if docState, ok := indexer.tableLock.docsState[docIndex.DocId]; ok && docState <= 1 {
				// ok && docState == 0 表示存在于索引中，需先删除再添加
				// ok && docState == 1 表示不一定存在于索引中，等待删除，需先删除再添加
				if position != i {
					indexer.addCacheLock.addCache[position], indexer.addCacheLock.addCache[i] =
						indexer.addCacheLock.addCache[i], indexer.addCacheLock.addCache[position]
				}
				if docState == 0 {
					indexer.removeCacheLock.Lock()
					indexer.removeCacheLock.removeCache[indexer.removeCacheLock.removeCachePointer] =
						docIndex.DocId
					indexer.removeCacheLock.removeCachePointer++
					indexer.removeCacheLock.Unlock()
					indexer.tableLock.docsState[docIndex.DocId] = 1
					indexer.numDocuments--
				}
				position++
			} else if !ok {
				indexer.tableLock.docsState[docIndex.DocId] = 2
			}
		}

		indexer.tableLock.Unlock()
		if indexer.RemoveDocumentToCache(0, forceUpdate) {
			// 只有当存在于索引表中的文档已被删除，其才可以重新加入到索引表中
			position = 0
		}

		addCachedDocuments := indexer.addCacheLock.addCache[position:indexer.addCacheLock.addCachePointer]
		indexer.addCacheLock.addCachePointer = position
		indexer.addCacheLock.Unlock()
		sort.Sort(addCachedDocuments)
		indexer.AddDocuments(&addCachedDocuments)
	} else {
		indexer.addCacheLock.Unlock()
	}
}

// 向反向索引表中加入 ADDCACHE 中所有文档
func (indexer *Indexer) AddDocuments(documents *types.DocumentsIndex) {
	if indexer.initialized == false {
		log.Fatal("索引器尚未初始化")
	}

	indexer.tableLock.Lock()
	defer indexer.tableLock.Unlock()
	indexPointers := make(map[string]int, len(indexer.tableLock.table))

	// DocId 递增顺序遍历插入文档保证索引移动次数最少
	for i, document := range *documents {
		if i < len(*documents)-1 && (*documents)[i].DocId == (*documents)[i+1].DocId {
			// 如果有重复文档加入，因为稳定排序，只加入最后一个
			continue
		}
		if docState, ok := indexer.tableLock.docsState[document.DocId]; ok && docState == 1 {
			// 如果此时 docState 仍为 1，说明该文档需被删除
			// docState 合法状态为 nil & 2，保证一定不会插入已经在索引表中的文档
			continue
		}

		// 更新文档关键词总长度
		if document.TokenLength != 0 {
			indexer.docTokenLengths[document.DocId] = float32(document.TokenLength)
			indexer.totalTokenLength += document.TokenLength
		}

		docIdIsNew := true
		for _, keyword := range document.Keywords {
			indices, foundKeyword := indexer.tableLock.table[keyword.Text]
			if !foundKeyword {
				// 如果没找到该搜索键则加入
				ti := KeywordIndices{}
				switch indexer.initOptions.IndexType {
				case types.LocationsIndex:
					ti.locations = [][]int{keyword.Starts}
				case types.FrequenciesIndex:
					ti.frequencies = []float32{keyword.Frequency}
				}
				ti.docIds = []uint64{document.DocId}
				indexer.tableLock.table[keyword.Text] = &ti
				continue
			}

			// 查找应该插入的位置，且索引一定不存在
			position, _ := indexer.searchIndex(
				indices, indexPointers[keyword.Text], indexer.getIndexLength(indices)-1, document.DocId)
			indexPointers[keyword.Text] = position
			switch indexer.initOptions.IndexType {
			case types.LocationsIndex:
				indices.locations = append(indices.locations, []int{})
				copy(indices.locations[position+1:], indices.locations[position:])
				indices.locations[position] = keyword.Starts
			case types.FrequenciesIndex:
				indices.frequencies = append(indices.frequencies, float32(0))
				copy(indices.frequencies[position+1:], indices.frequencies[position:])
				indices.frequencies[position] = keyword.Frequency
			}
			indices.docIds = append(indices.docIds, 0)
			copy(indices.docIds[position+1:], indices.docIds[position:])
			indices.docIds[position] = document.DocId
		}

		// 更新文章状态和总数
		if docIdIsNew {
			indexer.tableLock.docsState[document.DocId] = 0
			indexer.numDocuments++
		}
	}
}

// 向 REMOVECACHE 中加入一个待删除文档
// 返回值表示文档是否在索引表中被删除
func (indexer *Indexer) RemoveDocumentToCache(docId uint64, forceUpdate bool) bool {
	if indexer.initialized == false {
		log.Fatal("索引器尚未初始化")
	}

	indexer.removeCacheLock.Lock()
	if docId != 0 {
		indexer.tableLock.Lock()
		if docState, ok := indexer.tableLock.docsState[docId]; ok && docState == 0 {
			indexer.removeCacheLock.removeCache[indexer.removeCacheLock.removeCachePointer] = docId
			indexer.removeCacheLock.removeCachePointer++
			indexer.tableLock.docsState[docId] = 1
			indexer.numDocuments--
		} else if ok && docState == 2 {
			// 删除一个等待加入的文档
			indexer.tableLock.docsState[docId] = 1
		} else if !ok {
			// 若文档不存在，则无法判断其是否在 addCache 中，需避免这样的操作
		}
		indexer.tableLock.Unlock()
	}

	if indexer.removeCacheLock.removeCachePointer > 0 &&
		(indexer.removeCacheLock.removeCachePointer >= indexer.initOptions.DocCacheSize ||
			forceUpdate) {
		removeCachedDocuments := indexer.removeCacheLock.removeCache[:indexer.removeCacheLock.removeCachePointer]
		indexer.removeCacheLock.removeCachePointer = 0
		indexer.removeCacheLock.Unlock()
		sort.Sort(removeCachedDocuments)
		indexer.RemoveDocuments(&removeCachedDocuments)
		return true
	}
	indexer.removeCacheLock.Unlock()
	return false
}

// 向反向索引表中删除 REMOVECACHE 中所有文档
func (indexer *Indexer) RemoveDocuments(documents *types.DocumentsId) {
	if indexer.initialized == false {
		log.Fatal("索引器尚未初始化")
	}

	indexer.tableLock.Lock()
	defer indexer.tableLock.Unlock()

	// 更新文档关键词总长度，删除文档状态
	for _, docId := range *documents {
		indexer.totalTokenLength -= indexer.docTokenLengths[docId]
		delete(indexer.docTokenLengths, docId)
		delete(indexer.tableLock.docsState, docId)
	}

	for keyword, indices := range indexer.tableLock.table {
		indicesTop, indicesPointer := 0, 0
		documentsPointer := sort.Search(
			len(*documents), func(i int) bool { return (*documents)[i] >= indices.docIds[0] })
		// 双指针扫描，进行批量删除操作
		for documentsPointer < len(*documents) && indicesPointer < indexer.getIndexLength(indices) {
			if indices.docIds[indicesPointer] < (*documents)[documentsPointer] {
				if indicesTop != indicesPointer {
					switch indexer.initOptions.IndexType {
					case types.LocationsIndex:
						indices.locations[indicesTop] = indices.locations[indicesPointer]
					case types.FrequenciesIndex:
						indices.frequencies[indicesTop] = indices.frequencies[indicesPointer]
					}
					indices.docIds[indicesTop] = indices.docIds[indicesPointer]
				}
				indicesTop++
				indicesPointer++
			} else if indices.docIds[indicesPointer] == (*documents)[documentsPointer] {
				indicesPointer++
				documentsPointer++
			} else {
				documentsPointer++
			}
		}
		if indicesTop != indicesPointer {
			switch indexer.initOptions.IndexType {
			case types.LocationsIndex:
				indices.locations = append(
					indices.locations[:indicesTop], indices.locations[indicesPointer:]...)
			case types.FrequenciesIndex:
				indices.frequencies = append(
					indices.frequencies[:indicesTop], indices.frequencies[indicesPointer:]...)
			}
			indices.docIds = append(
				indices.docIds[:indicesTop], indices.docIds[indicesPointer:]...)
		}
		if len(indices.docIds) == 0 {
			delete(indexer.tableLock.table, keyword)
		}
	}
}

// 查找包含全部搜索键(AND操作)的文档
// 当docIds不为nil时仅从docIds指定的文档中查找
func (indexer *Indexer) Lookup(
	tokens []string, labels []string, docIds map[uint64]bool, countDocsOnly bool) (docs []types.IndexedDocument, numDocs int) {
	if indexer.initialized == false {
		log.Fatal("索引器尚未初始化")
	}

	if indexer.numDocuments == 0 {
		return
	}
	numDocs = 0

	// 合并关键词和标签为搜索键
	keywords := make([]string, len(tokens)+len(labels))
	copy(keywords, tokens)
	copy(keywords[len(tokens):], labels)

	indexer.tableLock.RLock()
	defer indexer.tableLock.RUnlock()
	table := make([]*KeywordIndices, len(keywords))
	for i, keyword := range keywords {
		indices, found := indexer.tableLock.table[keyword]
		if !found {
			// 当反向索引表中无此搜索键时直接返回
			return
		} else {
			// 否则加入反向表中
			table[i] = indices
		}
	}

	// 当没有找到时直接返回
	if len(table) == 0 {
		return
	}

	// 归并查找各个搜索键出现文档的交集
	// 从后向前查保证先输出DocId较大文档
	indexPointers := make([]int, len(table))
	for iTable := 0; iTable < len(table); iTable++ {
		indexPointers[iTable] = indexer.getIndexLength(table[iTable]) - 1
	}
	// 平均文本关键词长度，用于计算BM25
	avgDocLength := indexer.totalTokenLength / float32(indexer.numDocuments)
	for ; indexPointers[0] >= 0; indexPointers[0]-- {
		// 以第一个搜索键出现的文档作为基准，并遍历其他搜索键搜索同一文档
		baseDocId := indexer.getDocId(table[0], indexPointers[0])
		if docIds != nil {
			if _, found := docIds[baseDocId]; !found {
				continue
			}
		}
		iTable := 1
		found := true
		for ; iTable < len(table); iTable++ {
			// 二分法比简单的顺序归并效率高，也有更高效率的算法，
			// 但顺序归并也许是更好的选择，考虑到将来需要用链表重新实现
			// 以避免反向表添加新文档时的写锁。
			// TODO: 进一步研究不同求交集算法的速度和可扩展性。
			position, foundBaseDocId := indexer.searchIndex(table[iTable],
				0, indexPointers[iTable], baseDocId)
			if foundBaseDocId {
				indexPointers[iTable] = position
			} else {
				if position == 0 {
					// 该搜索键中所有的文档ID都比baseDocId大，因此已经没有
					// 继续查找的必要。
					return
				} else {
					// 继续下一indexPointers[0]的查找
					indexPointers[iTable] = position - 1
					found = false
					break
				}
			}
		}

		if found {
			if docState, ok := indexer.tableLock.docsState[baseDocId]; !ok || docState != 0 {
				continue
			}
			indexedDoc := types.IndexedDocument{}

			// 当为LocationsIndex时计算关键词紧邻距离
			if indexer.initOptions.IndexType == types.LocationsIndex {
				// 计算有多少关键词是带有距离信息的
				numTokensWithLocations := 0
				for i, t := range table[:len(tokens)] {
					if len(t.locations[indexPointers[i]]) > 0 {
						numTokensWithLocations++
					}
				}
				if numTokensWithLocations != len(tokens) {
					if !countDocsOnly {
						docs = append(docs, types.IndexedDocument{
							DocId: baseDocId,
						})
					}
					numDocs++
					//当某个关键字对应多个文档且有lable关键字存在时，若直接break,将会丢失相当一部分搜索结果
					continue
				}

				// 计算搜索键在文档中的紧邻距离
				tokenProximity, tokenLocations := computeTokenProximity(table[:len(tokens)], indexPointers, tokens)
				indexedDoc.TokenProximity = int32(tokenProximity)
				indexedDoc.TokenSnippetLocations = tokenLocations

				// 添加TokenLocations
				indexedDoc.TokenLocations = make([][]int, len(tokens))
				for i, t := range table[:len(tokens)] {
					indexedDoc.TokenLocations[i] = t.locations[indexPointers[i]]
				}
			}

			// 当为LocationsIndex或者FrequenciesIndex时计算BM25
			if indexer.initOptions.IndexType == types.LocationsIndex ||
				indexer.initOptions.IndexType == types.FrequenciesIndex {
				bm25 := float32(0)
				d := indexer.docTokenLengths[baseDocId]
				for i, t := range table[:len(tokens)] {
					var frequency float32
					if indexer.initOptions.IndexType == types.LocationsIndex {
						frequency = float32(len(t.locations[indexPointers[i]]))
					} else {
						frequency = t.frequencies[indexPointers[i]]
					}

					// 计算BM25
					if len(t.docIds) > 0 && frequency > 0 && indexer.initOptions.BM25Parameters != nil && avgDocLength != 0 {
						// 带平滑的idf
						idf := float32(math.Log2(float64(indexer.numDocuments)/float64(len(t.docIds)) + 1))
						k1 := indexer.initOptions.BM25Parameters.K1
						b := indexer.initOptions.BM25Parameters.B
						bm25 += idf * frequency * (k1 + 1) / (frequency + k1*(1-b+b*d/avgDocLength))
					}
				}
				indexedDoc.BM25 = float32(bm25)
			}

			indexedDoc.DocId = baseDocId
			if !countDocsOnly {
				docs = append(docs, indexedDoc)
			}
			numDocs++
		}
	}
	return
}

// 二分法查找indices中某文档的索引项
// 第一个返回参数为找到的位置或需要插入的位置
// 第二个返回参数标明是否找到
func (indexer *Indexer) searchIndex(
	indices *KeywordIndices, start int, end int, docId uint64) (int, bool) {
	// 特殊情况
	if indexer.getIndexLength(indices) == start {
		return start, false
	}
	if docId < indexer.getDocId(indices, start) {
		return start, false
	} else if docId == indexer.getDocId(indices, start) {
		return start, true
	}
	if docId > indexer.getDocId(indices, end) {
		return end + 1, false
	} else if docId == indexer.getDocId(indices, end) {
		return end, true
	}

	// 二分
	var middle int
	for end-start > 1 {
		middle = (start + end) / 2
		if docId == indexer.getDocId(indices, middle) {
			return middle, true
		} else if docId > indexer.getDocId(indices, middle) {
			start = middle
		} else {
			end = middle
		}
	}
	return end, false
}

// 计算搜索键在文本中的紧邻距离
//
// 假定第 i 个搜索键首字节出现在文本中的位置为 P_i，长度 L_i
// 紧邻距离计算公式为
//
// 	ArgMin(Sum(Abs(P_(i+1) - P_i - L_i)))
//
// 具体由动态规划实现，依次计算前 i 个 token 在每个出现位置的最优值。
// 选定的 P_i 通过 tokenLocations 参数传回。
func computeTokenProximity(table []*KeywordIndices, indexPointers []int, tokens []string) (
	minTokenProximity int, tokenLocations []int) {
	minTokenProximity = -1
	tokenLocations = make([]int, len(tokens))

	var (
		currentLocations, nextLocations []int
		currentMinValues, nextMinValues []int
		path                            [][]int
	)

	// 初始化路径数组
	path = make([][]int, len(tokens))
	for i := 1; i < len(path); i++ {
		path[i] = make([]int, len(table[i].locations[indexPointers[i]]))
	}

	// 动态规划
	currentLocations = table[0].locations[indexPointers[0]]
	currentMinValues = make([]int, len(currentLocations))
	for i := 1; i < len(tokens); i++ {
		nextLocations = table[i].locations[indexPointers[i]]
		nextMinValues = make([]int, len(nextLocations))
		for j, _ := range nextMinValues {
			nextMinValues[j] = -1
		}

		var iNext int
		for iCurrent, currentLocation := range currentLocations {
			if currentMinValues[iCurrent] == -1 {
				continue
			}
			for iNext+1 < len(nextLocations) && nextLocations[iNext+1] < currentLocation {
				iNext++
			}

			update := func(from int, to int) {
				if to >= len(nextLocations) {
					return
				}
				value := currentMinValues[from] + utils.AbsInt(nextLocations[to]-currentLocations[from]-len(tokens[i-1]))
				if nextMinValues[to] == -1 || value < nextMinValues[to] {
					nextMinValues[to] = value
					path[i][to] = from
				}
			}

			// 最优解的状态转移只发生在左右最接近的位置
			update(iCurrent, iNext)
			update(iCurrent, iNext+1)
		}

		currentLocations = nextLocations
		currentMinValues = nextMinValues
	}

	// 找出最优解
	var cursor int
	for i, value := range currentMinValues {
		if value == -1 {
			continue
		}
		if minTokenProximity == -1 || value < minTokenProximity {
			minTokenProximity = value
			cursor = i
		}
	}

	// 从路径倒推出最优解的位置
	for i := len(tokens) - 1; i >= 0; i-- {
		if i != len(tokens)-1 {
			cursor = path[i+1][cursor]
		}
		tokenLocations[i] = table[i].locations[indexPointers[i]][cursor]
	}
	return
}
