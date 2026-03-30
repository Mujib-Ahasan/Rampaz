package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type NodeMetService struct {
	client *kubernetes.NodeMetricsClient
}

func NewNodeMetService(client *kubernetes.NodeMetricsClient) *NodeMetService {
	return &NodeMetService{client: client}
}

func (s *NodeMetService) StreamNodeStats(ctx context.Context, nodeName string, send func(*pb.NodeStatsResponse) error) error {

	if nodeName == "" {
		return fmt.Errorf("node metrics stream: node name is required")
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {

		case <-ctx.Done():
			return nil

		case <-ticker.C:
			metrics, err := s.client.GetNodeMetrics(ctx, nodeName)
			if err != nil {
				return fmt.Errorf("node metrics stream: fetch metrics for nodestatsresponse response: %w", err)
			}

			cpu := metrics.Usage.Cpu().MilliValue()
			mem := metrics.Usage.Memory().Value() / (1024 * 1024)

			resp := &pb.NodeStatsResponse{
				Name:   metrics.Name,
				Cpu:    fmt.Sprintf("%dm", cpu),
				Memory: fmt.Sprintf("%dMi", mem),
			}

			if err := send(resp); err != nil {
				return fmt.Errorf("node metrics stream: send response for node %q: %w", nodeName, err)
			}
		}
	}
}
