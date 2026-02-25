package api

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	"github.com/Mujib-Ahasan/Rampaz/internal/service"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

type K8SServer struct {
	pb.UnimplementedK8SInfoServer
	PodService        *service.PodService
	NodeService       *service.NodeService
	EventService      *service.EventService
	PodMetService     *service.PodMetService
	NodeMetService    *service.NodeMetService
	DeploymentService *service.DeploymentService
	ReplicaSetservice *service.ReplicaSetService
	DaemonSetService  *service.DaemonSetService
}

func (s *K8SServer) ListPods(ctx context.Context, req *pb.NamespaceRequest) (*pb.PodListResponse, error) {
	endpoint := "list_pods"
	status := "success"
	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)
	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, status).
			Inc()
	}()

	pods, err := s.PodService.ListPods(ctx, req.Namespace)
	if err != nil {
		status = "error"
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

	return &pb.PodListResponse{Pods: result}, nil
}

func (s *K8SServer) GetNodeStats(ctx context.Context, req *pb.NodeRequest) (*pb.NodeStatsResponse, error) {
	endpoint := "get_node_stats"
	status := "success"
	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, status).
			Inc()
	}()

	node, err := s.NodeService.GetNodeStats(ctx, req.NodeName)
	if err != nil {
		status = "error"
		return nil, fmt.Errorf("failed to get node %s: %v", req.NodeName, err)
	}

	cpu := node.Status.Capacity.Cpu().String()
	memory := node.Status.Capacity.Memory().String()

	result := &pb.NodeStatsResponse{
		Name:   node.Name,
		Cpu:    cpu,
		Memory: memory,
	}

	return result, nil
}
