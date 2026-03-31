package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NamespaceMetricsService struct {
	podClient        *kubernetes.PodClient
	podMetricsClient *kubernetes.PodMetricsClient
}

func NewNamespaceMetricsService(podClient *kubernetes.PodClient, podMetricsClient *kubernetes.PodMetricsClient) *NamespaceMetricsService {
	return &NamespaceMetricsService{podClient: podClient, podMetricsClient: podMetricsClient}
}

func (s *NamespaceMetricsService) GetNamespaceMetrics(ctx context.Context, namespace string) (*pb.NamespaceMetricsResponse, error) {
	pods, err := s.podClient.ListPods(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch pod data for namespacemetricsresponse: %w", err)
	}

	var total, running, pending, failed, succeeded, unknown int32

	for _, pod := range pods.Items {
		total++

		switch pod.Status.Phase {
		case corev1.PodRunning:
			running++
		case corev1.PodPending:
			pending++
		case corev1.PodFailed:
			failed++
		case corev1.PodSucceeded:
			succeeded++
		default:
			unknown++
		}
	}

	var totalCPURequests, totalMemoryRequests, totalCPULimits, totalMemoryLimits resource.Quantity

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if cpuReq, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
				totalCPURequests.Add(cpuReq)
			}
			if memReq, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
				totalMemoryRequests.Add(memReq)
			}
			if cpuLim, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
				totalCPULimits.Add(cpuLim)
			}
			if memLim, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
				totalMemoryLimits.Add(memLim)
			}
		}
	}

	var totalCPUUsage, totalMemoryUsage resource.Quantity
	podMetrics, err := s.podMetricsClient.GetPodMetrics(ctx, namespace)

	if err == nil {
		for _, pm := range podMetrics.Items {
			for _, c := range pm.Containers {
				if cpu, ok := c.Usage[corev1.ResourceCPU]; ok {
					totalCPUUsage.Add(cpu)
				}
				if mem, ok := c.Usage[corev1.ResourceMemory]; ok {
					totalMemoryUsage.Add(mem)
				}
			}
		}
	}
	//becareful with totalCPUUsage (1m = 1,000,000n)

	return &pb.NamespaceMetricsResponse{
		Namespace:     namespace,
		TotalPods:     total,
		RunningPods:   running,
		PendingPods:   pending,
		FailedPods:    failed,
		SucceededPods: succeeded,
		UnknownPods:   unknown,
		Usage: &pb.ResourceQuantity{
			Cpu:    totalCPUUsage.String(),
			Memory: totalMemoryUsage.String(),
		},
		Requests: &pb.ResourceQuantity{
			Cpu:    totalCPURequests.String(),
			Memory: totalMemoryRequests.String(),
		},
		Limits: &pb.ResourceQuantity{
			Cpu:    totalCPULimits.String(),
			Memory: totalMemoryLimits.String(),
		},
	}, nil
}
