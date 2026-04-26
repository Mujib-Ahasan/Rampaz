package ai

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Answer    string `json:"answer"`
	SessionID string `json:"sessionId"`
}

type ToolCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}

type ToolResult struct {
	Name  string `json:"name"`
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}
