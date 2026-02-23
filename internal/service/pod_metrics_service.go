package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type PodMetService struct {
	client *kubernetes.PodMetricsClient
}

func NewPodMetService(client *kubernetes.PodMetricsClient) *PodMetService {
	return &PodMetService{client: client}
}

func (s *PodMetService) StreamPodStats(ctx context.Context, namespace string, send func(*pb.PodStatsResponse) error) error {

	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			metricsList, err := s.client.GetPodMetrics(ctx, namespace)
			if err != nil {
				return fmt.Errorf("failed to fetch pod metrics: %v", err)
			}

			for _, m := range metricsList.Items {

				var totalCPU int64
				var totalMem int64

				for _, c := range m.Containers {
					cpuQuantity := c.Usage["cpu"]
					memQuantity := c.Usage["memory"]
					totalCPU += cpuQuantity.MilliValue()
					totalMem += memQuantity.Value() / (1024 * 1024)
				}
				resp := &pb.PodStatsResponse{
					Name:      m.Name,
					Namespace: m.Namespace,
					Cpu:       fmt.Sprintf("%dm", totalCPU),
					Memory:    fmt.Sprintf("%dMi", totalMem),
				}

				if err := send(resp); err != nil {
					return err
				}
			}
		}
	}
}
