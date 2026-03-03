package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mindoc-org/mindoc/conf"
	"github.com/mindoc-org/mindoc/models"
	"github.com/mindoc-org/mindoc/utils/sqltil"
)

// GetGlobalSearchMcpTool 获取全局搜索的mcp工具
func GetGlobalSearchMcpTool() mcp.Tool {
	return mcp.NewTool("MinDocGlobalSearch",
		mcp.WithDescription("MinDoc全局文档内容搜索"),
		mcp.WithString("keyword",
			mcp.Required(),
			mcp.Description("要执行全局搜索的关键词，多个搜索关键词请用空格分割，请使用最少的关键词来检索，结果中只会出现包含全部关键词的结果，过多的无关词会导致更少的检索结果"),
		),
		mcp.WithNumber("pageIndex",
			mcp.Required(),
			mcp.Description("全局搜索时指定分页的顺序下标，每页最多有10条结果，建议只查看1-10页文档内容的搜索结果"),
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

func globalSearchFunction(keyword string, pageIndex int) (int, []*models.DocumentSearchResult) {
	memberId := 0
	searchResult, totalCount, err := models.NewDocumentSearchResult().FindToPager(sqltil.EscapeLike(keyword),
		pageIndex, conf.PageSize, memberId)
	if err != nil {
		return 0, make([]*models.DocumentSearchResult, 0)
	}
	return totalCount, searchResult
}
