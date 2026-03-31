package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type NodeResourceService struct {
	client *kubernetes.NodeResourceClient
}

func NewNodeResourceService(client *kubernetes.NodeResourceClient) *NodeResourceService {
	return &NodeResourceService{client: client}
}

func (s *NodeResourceService) GetNodeResourceAllocation(ctx context.Context, nodeName string) (*pb.NodeResourceAllocationResponse, error) {
	node, err := s.client.GetNode(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("node resource allocation: get node %q: %w", nodeName, err)
	}

	capacityCPU := node.Status.Capacity[corev1.ResourceCPU]
	capacityMem := node.Status.Capacity[corev1.ResourceMemory]
	allocCPU := node.Status.Allocatable[corev1.ResourceCPU]
	allocMem := node.Status.Allocatable[corev1.ResourceMemory]

	pods, err := s.client.ListPodsByNode(ctx, nodeName)
	if err != nil {
		return nil, fmt.Errorf("node resource allocation: list pods for node %q: %w", nodeName, err)
	}

	var totalCPUReq, totalMemReq, totalCPULim, totalMemLim resource.Quantity
	var podCount int32

	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			continue
		}
		podCount++

		for _, c := range pod.Spec.Containers {
			if cpuReq, ok := c.Resources.Requests[corev1.ResourceCPU]; ok {
				totalCPUReq.Add(cpuReq)
			}
			if memReq, ok := c.Resources.Requests[corev1.ResourceMemory]; ok {
				totalMemReq.Add(memReq)
			}
			if cpuLim, ok := c.Resources.Limits[corev1.ResourceCPU]; ok {
				totalCPULim.Add(cpuLim)
			}
			if memLim, ok := c.Resources.Limits[corev1.ResourceMemory]; ok {
				totalMemLim.Add(memLim)
			}
		}
	}

	usageCPU := "0m"
	usageMem := "0Mi"

	nodeMetrics, err := s.client.GetNodeMetrics(ctx, nodeName)
	if err == nil {
		if cpu, ok := nodeMetrics.Usage[corev1.ResourceCPU]; ok {
			usageCPU = fmt.Sprintf("%dm", cpu.MilliValue())
		}
		if mem, ok := nodeMetrics.Usage[corev1.ResourceMemory]; ok {
			usageMem = fmt.Sprintf("%dMi", mem.Value()/(1024*1024))
		}
	}

	return &pb.NodeResourceAllocationResponse{
		NodeName: nodeName,
		Capacity: &pb.ResourceQuantity{
			Cpu:    capacityCPU.String(),
			Memory: capacityMem.String(),
		},
		Allocatable: &pb.ResourceQuantity{
			Cpu:    allocCPU.String(),
			Memory: allocMem.String(),
		},
		Requests: &pb.ResourceQuantity{
			Cpu:    totalCPUReq.String(),
			Memory: totalMemReq.String(),
		},
		Limits: &pb.ResourceQuantity{
			Cpu:    totalCPULim.String(),
			Memory: totalMemLim.String(),
		},
		Usage: &pb.ResourceQuantity{
			Cpu:    usageCPU,
			Memory: usageMem,
		},
		PodCount: podCount,
	}, nil
}
