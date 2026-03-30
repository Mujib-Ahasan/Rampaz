package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type PodService struct {
	podClient *kubernetes.PodClient
}

func NewPodService(pc *kubernetes.PodClient) *PodService {
	return &PodService{podClient: pc}
}

func (s *PodService) ListPods(ctx context.Context, namespace string) (*pb.PodListResponse, error) {
	pods, err := s.podClient.ListPods(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch pod data for podlistresponse: %w", err)
	}

	var result []*pb.Pod

	for _, pod := range pods.Items {
		result = append(result, &pb.Pod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			NodeName:  pod.Spec.NodeName,
			Status:    string(pod.Status.Phase),
		})
	}

	return &pb.PodListResponse{Pods: result}, nil

}
