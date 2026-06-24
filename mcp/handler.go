package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/controllers"
)

// GetGlobalSearchMcpTool 获取全局搜索的mcp工具
func GetGlobalSearchMcpTool() mcp.Tool {
	return mcp.NewTool("MinDocGlobalSearch",
		mcp.WithDescription("MinDoc全局文档内容搜索"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("要执行全局搜索的关键词，多个搜索关键词请用空格分割，请尽量使用最精简的关键词来检索，不要输入无关词汇"),
		),
		mcp.WithNumber("pageIndex",
			mcp.Required(),
			mcp.Description("全局搜索时指定分页的顺序下标，每页最多有10条结果，建议只查看最相关的1-10页文档内容的搜索结果"),
			mcp.Enum("1", "2", "3", "4", "5", "6", "7", "8", "9", "10"),
		),
	)
}

// GlobalSearchMcpHandler 全局搜索的mcp处理函数
func GlobalSearchMcpHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paramMap := request.Params.Arguments.(map[string]any)
	pageIndex := 1
	if v, ok := paramMap["pageIndex"].(float64); ok {
		pageIndex = int(v)
	}
	totalCount, result := globalSearchFunction(paramMap["keyword"].(string), pageIndex)
	jsonContent, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultStructuredOnly(map[string]any{
			"totalCount": 0,
			"result":     make([]map[string]any, 0),
		}), err
	}

	structContent := make([]map[string]any, 0)
	err = json.Unmarshal(jsonContent, &structContent)
	if err != nil {
		return mcp.NewToolResultStructuredOnly(map[string]any{
			"totalCount": 0,
			"result":     make([]map[string]any, 0),
		}), err
	}

	return mcp.NewToolResultStructuredOnly(map[string]any{
		"totalCount": totalCount,
		"result":     structContent,
	}), nil
}

func globalSearchFunction(keyword string, pageIndex int) (int, []*controllers.SearchV2RawResult) {
	memberId := 0
	// 使用底层搜索函数
	searchResult, _, totalCount, err := controllers.PerformSearchV2Raw(keyword, pageIndex, conf.PageSize, memberId)
	if err != nil {
		return 0, make([]*controllers.SearchV2RawResult, 0)
	}
	return totalCount, searchResult
}
