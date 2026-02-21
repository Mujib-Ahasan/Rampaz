package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

type NodeService struct {
	nodeClient *kubernetes.NodeClient
}

func NewNodeService(nc *kubernetes.NodeClient) *NodeService {
	return &NodeService{nodeClient: nc}
}

func (s *NodeService) Getnode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	return s.nodeClient.ListNodes(ctx, nodeName)
}
