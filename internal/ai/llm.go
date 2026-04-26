package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const deepSeekBaseURL = "https://api.deepseek.com/chat/completions"

type LLMClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewLLMClient() *LLMClient {
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-v4-flash"
	}

	return &LLMClient{
		apiKey: os.Getenv("DEEPSEEK_API_KEY"),
		model:  model,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type chatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

type deepSeekRequest struct {
	Model       string         `json:"model"`
	Messages    []chatMessage  `json:"messages"`
	Tools       []deepSeekTool `json:"tools,omitempty"`
	ToolChoice  any            `json:"tool_choice,omitempty"`
	Temperature float64        `json:"temperature,omitempty"`
	Stream      bool           `json:"stream"`
}

type deepSeekTool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (c *LLMClient) Plan(ctx context.Context, userMessage string, tools []ToolSchema) ([]ToolCall, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("DEEPSEEK_API_KEY is not set")
	}

	req := deepSeekRequest{
		Model: c.model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: PlannerPrompt,
			},
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Tools:       toDeepSeekTools(tools),
		ToolChoice:  "auto",
		Temperature: 0.1,
		Stream:      false,
	}

	resp, err := c.doChat(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("deepseek returned no choices")
	}

	msg := resp.Choices[0].Message

	calls := make([]ToolCall, 0, len(msg.ToolCalls))
	for _, tc := range msg.ToolCalls {
		var args map[string]any
		if tc.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				return nil, fmt.Errorf("invalid tool args for %s: %w", tc.Function.Name, err)
			}
		}

		calls = append(calls, ToolCall{
			Name: tc.Function.Name,
			Args: args,
		})
	}

	return calls, nil
}

func (c *LLMClient) Answer(ctx context.Context, userMessage string, clusterContext string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY is not set")
	}

	req := deepSeekRequest{
		Model: c.model,
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: AnswerPrompt,
			},
			{
				Role: "user",
				Content: fmt.Sprintf(
					"User question:\n%s\n\nCluster context:\n%s",
					userMessage,
					clusterContext,
				),
			},
		},
		Temperature: 0.2,
		Stream:      false,
	}

	resp, err := c.doChat(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("deepseek returned no choices")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *LLMClient) doChat(ctx context.Context, body deepSeekRequest) (*deepSeekResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, deepSeekBaseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out deepSeekResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode deepseek response: %w: %s", err, string(raw))
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if out.Error != nil {
			return nil, fmt.Errorf("deepseek api error: %s", out.Error.Message)
		}
		return nil, fmt.Errorf("deepseek api status %d: %s", res.StatusCode, string(raw))
	}

	return &out, nil
}

func toDeepSeekTools(tools []ToolSchema) []deepSeekTool {
	out := make([]deepSeekTool, 0, len(tools))

	for _, tool := range tools {
		out = append(out, deepSeekTool{
			Type: "function",
			Function: ToolFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
			},
		})
	}

	return out
}
