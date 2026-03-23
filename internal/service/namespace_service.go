package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type NamespaceService struct {
	client *kubernetes.NamespaceClient
}

func NewNamespaceService(client *kubernetes.NamespaceClient) *NamespaceService {
	return &NamespaceService{client: client}
}

func (c *NamespaceService) ListNamespaces(ctx context.Context) ([]*pb.NamespaceInfo, error) {
	namespaces, err := c.client.List(ctx)
	if err != nil {
		return nil, err
	}

	var result []*pb.NamespaceInfo

	for _, ns := range namespaces {
		result = append(result, transformNamespace(&ns))
	}

	return result, nil

}

func transformNamespace(ns *corev1.Namespace) *pb.NamespaceInfo {
	return &pb.NamespaceInfo{
		Name:  ns.Name,
		Phase: string(ns.Status.Phase),
		Age:   calculateAge(ns.CreationTimestamp.Time),
	}
}
