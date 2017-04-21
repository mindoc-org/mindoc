package search

import (
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
)

var searcher engine.Engine

func init()  {
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: "data/dictionary.txt",
		StopTokenFile:         "data/stop_tokens.txt",
		IndexerInitOptions: &types.IndexerInitOptions{
			IndexType: types.LocationsIndex,
		},
	})
}