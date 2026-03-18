package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	corev1 "k8s.io/api/core/v1"
)

type NodeService struct {
	client *kubernetes.NodeClient
}

func NewNodeService(client *kubernetes.NodeClient) *NodeService {
	return &NodeService{client: client}
}

func (s *NodeService) ListNodes(ctx context.Context) ([]*pb.NodeInfo, error) {
	nodes, err := s.client.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	var result []*pb.NodeInfo

	for _, node := range nodes {
		result = append(result, transformNode(&node))
	}

	return result, nil
}

func transformNode(node *corev1.Node) *pb.NodeInfo {
	var internalIP string

	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			internalIP = addr.Address
		}
	}

	return &pb.NodeInfo{
		Name:       node.Name,
		InternalIp: internalIP,
		Phase:      string(node.Status.Phase),
		Age:        calculateAge(node.CreationTimestamp.Time),
	}
}
