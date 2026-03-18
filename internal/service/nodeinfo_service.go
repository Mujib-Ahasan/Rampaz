package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

type NodeInfoService struct {
	nodeClient *kubernetes.NodeInfoClient
}

func NewNodeInfoService(nc *kubernetes.NodeInfoClient) *NodeInfoService {
	return &NodeInfoService{nodeClient: nc}
}

func (s *NodeInfoService) GetNodeStats(ctx context.Context, nodeName string) (*corev1.Node, error) {
	return s.nodeClient.NodeInfo(ctx, nodeName)
}
