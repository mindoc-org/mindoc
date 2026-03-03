package mcp

import (
	"context"
	"net/http"

	"github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
)

// AuthMiddleware 返回一个中间件函数，用于验证MCP请求中的认证令牌
func AuthMiddleware(ctx *beegoContext.Context) {
	presetMcpApiKey := web.AppConfig.DefaultString("mcp_api_key", "")
	mcpApiKeyParamValue := ctx.Request.URL.Query().Get("api_key")
	if presetMcpApiKey != mcpApiKeyParamValue {
		http.Error(ctx.ResponseWriter, "Missing or invalid mcp authorization key", http.StatusUnauthorized)
		return
	}

	// Add mcp_api_key to request context
	ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "mcp_api_key", mcpApiKeyParamValue))
}
