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
		return fmt.Errorf("node name required")
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
				return fmt.Errorf("metrics fetch failed: %w", err)
			}

			cpu := metrics.Usage.Cpu().MilliValue()
			mem := metrics.Usage.Memory().Value() / (1024 * 1024)

			resp := &pb.NodeStatsResponse{
				Name:   metrics.Name,
				Cpu:    fmt.Sprintf("%dm", cpu),
				Memory: fmt.Sprintf("%dMi", mem),
			}

			if err := send(resp); err != nil {
				return err
			}
		}
	}
}
