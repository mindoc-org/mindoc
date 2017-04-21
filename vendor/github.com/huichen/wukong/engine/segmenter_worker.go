package engine

import (
	"github.com/huichen/wukong/types"
)

type segmenterRequest struct {
	docId       uint64
	hash        uint32
	data        types.DocumentIndexData
	forceUpdate bool
}

func (engine *Engine) segmenterWorker() {
	for {
		request := <-engine.segmenterChannel
		if request.docId == 0 {
			if request.forceUpdate {
				for i := 0; i < engine.initOptions.NumShards; i++ {
					engine.indexerAddDocChannels[i] <- indexerAddDocumentRequest{forceUpdate: true}
				}
			}
			continue
		}

		shard := engine.getShard(request.hash)
		tokensMap := make(map[string][]int)
		numTokens := 0
		if !engine.initOptions.NotUsingSegmenter && request.data.Content != "" {
			// 当文档正文不为空时，优先从内容分词中得到关键词
			segments := engine.segmenter.Segment([]byte(request.data.Content))
			for _, segment := range segments {
				token := segment.Token().Text()
				if !engine.stopTokens.IsStopToken(token) {
					tokensMap[token] = append(tokensMap[token], segment.Start())
				}
			}
			numTokens = len(segments)
		} else {
			// 否则载入用户输入的关键词
			for _, t := range request.data.Tokens {
				if !engine.stopTokens.IsStopToken(t.Text) {
					tokensMap[t.Text] = t.Locations
				}
			}
			numTokens = len(request.data.Tokens)
		}

		// 加入非分词的文档标签
		for _, label := range request.data.Labels {
			if !engine.initOptions.NotUsingSegmenter {
				if !engine.stopTokens.IsStopToken(label) {
					//当正文中已存在关键字时，若不判断，位置信息将会丢失
					if _, ok := tokensMap[label]; !ok {
						tokensMap[label] = []int{}
					}
				}
			} else {
				//当正文中已存在关键字时，若不判断，位置信息将会丢失
				if _, ok := tokensMap[label]; !ok {
					tokensMap[label] = []int{}
				}
			}
		}

		indexerRequest := indexerAddDocumentRequest{
			document: &types.DocumentIndex{
				DocId:       request.docId,
				TokenLength: float32(numTokens),
				Keywords:    make([]types.KeywordIndex, len(tokensMap)),
			},
			forceUpdate: request.forceUpdate,
		}
		iTokens := 0
		for k, v := range tokensMap {
			indexerRequest.document.Keywords[iTokens] = types.KeywordIndex{
				Text: k,
				// 非分词标注的词频设置为0，不参与tf-idf计算
				Frequency: float32(len(v)),
				Starts:    v}
			iTokens++
		}

		engine.indexerAddDocChannels[shard] <- indexerRequest
		if request.forceUpdate {
			for i := 0; i < engine.initOptions.NumShards; i++ {
				if i == shard {
					continue
				}
				engine.indexerAddDocChannels[i] <- indexerAddDocumentRequest{forceUpdate: true}
			}
		}
		rankerRequest := rankerAddDocRequest{
			docId: request.docId, fields: request.data.Fields}
		engine.rankerAddDocChannels[shard] <- rankerRequest
	}
}
