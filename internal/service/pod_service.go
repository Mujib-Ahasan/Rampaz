package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

type PodService struct {
	podClient *kubernetes.PodClient
}

func NewPodService(pc *kubernetes.PodClient) *PodService {
	return &PodService{podClient: pc}
}

func (s *PodService) GetPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	return s.podClient.ListPods(ctx, namespace)
}
