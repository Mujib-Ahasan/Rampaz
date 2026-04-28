package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
)

type PodLogsService struct {
	podLogClient *kubernetes.PodLogsClient
}

func NewPodLogsService(pc *kubernetes.PodLogsClient) *PodLogsService {
	return &PodLogsService{podLogClient: pc}
}

func (s *PodLogsService) GetPodLogs(ctx context.Context, namespace, podName string, tailLines int64) (string, error) {
	logs, err := s.podLogClient.GetPodLogs(ctx, namespace, podName, tailLines)

	if err != nil {
		return "", fmt.Errorf("fetch pod log for pod list response: %w", err)
	}
	return logs, err
}
