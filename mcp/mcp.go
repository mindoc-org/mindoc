package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// MCPServer MinDoc MCP Server
type MCPServer struct {
	server *server.MCPServer
}

// NewMCPServer creates a new MinDoc MCP Server
func NewMCPServer() *MCPServer {
	mcpServer := server.NewMCPServer(
		"MinDoc MCP Server",
		"1.0.0",
		server.WithRecovery(),
	)

	mcpServer.AddTool(GetGlobalSearchMcpTool(), GlobalSearchMcpHandler)

	return &MCPServer{
		server: mcpServer,
	}
}

// ServeHTTP Run starts the server
func (s *MCPServer) ServeHTTP() *server.StreamableHTTPServer {
	return server.NewStreamableHTTPServer(s.server)
}
