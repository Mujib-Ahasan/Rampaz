package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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

	result := make([]*pb.Pod, 0, len(pods.Items))

	for _, pod := range pods.Items {
		var cpuReq, memReq, cpuLimit, memLimit resource.Quantity

		for _, container := range pod.Spec.Containers {
			if q, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				cpuReq.Add(q)
			}
			if q, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				memReq.Add(q)
			}
			if q, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				cpuLimit.Add(q)
			}
			if q, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				memLimit.Add(q)
			}
		}

		result = append(result, &pb.Pod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			NodeName:  pod.Spec.NodeName,
			Status:    string(pod.Status.Phase),
			Resources: &pb.ResourceRequirements{
				CpuRequest:    cpuReq.String(),
				MemoryRequest: memReq.String(),
				CpuLimit:      cpuLimit.String(),
				MemoryLimit:   memLimit.String(),
			},
			QosClass: string(pod.Status.QOSClass),
		})
	}

	return &pb.PodListResponse{Pods: result}, nil
}
