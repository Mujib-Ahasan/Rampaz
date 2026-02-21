package api

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/service"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type K8SServer struct {
	pb.UnimplementedK8SInfoServer
	PodService     *service.PodService
	NodeService    *service.NodeService
	EventService   *service.EventService
	PodMetService  *service.PodMetService
	NodeMetService *service.NodeMetService
}

func (s *K8SServer) GetPods(ctx context.Context, req *pb.PodRequest) (*pb.PodListResponse, error) {

	pods, err := s.PodService.GetPods(ctx, req.Namespace)
	if err != nil {
		return nil, err
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

	return &pb.PodListResponse{
		Pods: result,
	}, nil
}

func (s *K8SServer) GetNodeInfo(ctx context.Context, req *pb.NodeRequest) (*pb.NodeStatsResponse, error) {
	node, err := s.NodeService.Getnode(ctx, req.NodeName)
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %v", req.NodeName, err)
	}

	cpu := node.Status.Capacity.Cpu().String()
	memory := node.Status.Capacity.Memory().String()

	res := &pb.NodeStatsResponse{
		Name:   node.Name,
		Cpu:    cpu,
		Memory: memory,
	}

	return res, nil
}
