# Fork Notes — SafeIdea-LLC/claude-agent-sdk-go

**Upstream:** [schlunsen/claude-agent-sdk-go](https://github.com/schlunsen/claude-agent-sdk-go)
**Fork:** [SafeIdea-LLC/claude-agent-sdk-go](https://github.com/SafeIdea-LLC/claude-agent-sdk-go)

## Last upstream merge

**Date:** 2026-03-27
**Upstream commit:** `24eac3b` (Merge PR #38 — skills + missing hook events)
**Branch:** `merge/upstream-2026-03-27`

## Our additions (delta from upstream)

### 1. MCP initialize handshake (`types/mcp.go`)
Upstream's `SDKMCPServer.HandleMessage` didn't handle the `initialize` or
`notifications/initialized` JSON-RPC methods. Without this, Claude CLI's
control protocol hangs during MCP server setup.

### 2. OnToolResult callback (`types/mcp.go`)
`ToolResultCallback` type + `SetOnToolResult()` on SDKMCPServer. After a tool
executes, the callback fires with the `toolUseId` (from `_meta.claudecode/toolUseId`),
letting the host emit `EventToolResult` to the SSE stream so the frontend can
match results to their corresponding tool_call events.

### 3. WithMCPServer helper (`types/options.go`)
`ClaudeAgentOptions.WithMCPServer(name, server)` — convenience method to register
an in-process SDK MCP server. Declares it as `type: "sdk"` so MCP messages route
through the control protocol.

### 4. CLI resolution (`internal/transport/cli_resolve.go`)
Robust node/claude binary discovery with platform-specific handling.
Searches PATH, common install locations, and bundled runtime paths.

### 5. MCP server wiring in Query (`internal/query.go`)
Registers `MCPServer` instances from options into the query handler so
control protocol requests get routed to the correct server.

### 6. Platform-specific process handling
- `platform_args_windows.go` / `platform_args_default.go` — CLI arg adjustments
- `process_windows.go` / `process_default.go` — process group / signal handling

## Merge schedule

Audit upstream quarterly. Next check: 2026-06-27.
