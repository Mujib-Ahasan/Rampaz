package ai

import (
	"context"
)

type ChatService struct {
	llm      *LLMClient
	executor *ToolExecutor
}

func NewChatService(llm *LLMClient, executor *ToolExecutor) *ChatService {
	return &ChatService{
		llm:      llm,
		executor: executor,
	}
}

func (s *ChatService) Chat(ctx context.Context, message string) (string, error) {
	toolCalls, err := s.llm.Plan(ctx, message, AvailableTools())
	if err != nil {
		return "", err
	}

	if len(toolCalls) > 5 {
		toolCalls = toolCalls[:5]
	}

	results := make([]ToolResult, 0, len(toolCalls))

	for _, call := range toolCalls {
		result, err := s.executor.Execute(ctx, call)
		if err != nil {
			results = append(results, ToolResult{
				Name:  call.Name,
				Error: err.Error(),
			})
			continue
		}

		results = append(results, result)
	}

	clusterContext := BuildClusterContext(results)

	return s.llm.Answer(ctx, message, clusterContext)
}
