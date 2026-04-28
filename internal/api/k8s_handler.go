package api

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	"github.com/Mujib-Ahasan/Rampaz/internal/service"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type K8SServer struct {
	pb.UnimplementedK8SInfoServer
	PodService              *service.PodService
	NodeInfoService         *service.NodeInfoService
	EventService            *service.EventService
	PodMetService           *service.PodMetService
	NodeMetService          *service.NodeMetService
	DeploymentService       *service.DeploymentService
	ReplicaSetservice       *service.ReplicaSetService
	DaemonSetService        *service.DaemonSetService
	StatefulSetService      *service.StatefulSetService
	JobService              *service.JobService
	CronJobService          *service.CronJobService
	SVCService              *service.SVCService
	NamespaceService        *service.NamespaceService
	PVCService              *service.PVCService
	PVService               *service.PVService
	NodeService             *service.NodeService
	NetworkPolicyService    *service.NetworkPolicyService
	NamespaceSummaryService *service.SummaryService
	ClusterOverviewService  *service.SummaryService
	WorkloadService         *service.WorkloadService
	NamespaceMetricsService *service.NamespaceMetricsService
	NodeResourceService     *service.NodeResourceService
	PodLogService           *service.PodLogsService
	Logger                  *slog.Logger
}

func (s *K8SServer) ListPods(ctx context.Context, req *pb.NamespaceRequest) (*pb.PodListResponse, error) {
	endpoint := "list_pods"
	statusLabel := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)
	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, statusLabel).
			Inc()
	}()

	if req == nil {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "pod list request cannot be nil")
	}

	if req.Namespace == "" {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}

	pods, err := s.PodService.ListPods(ctx, req.Namespace)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("list pods failed", "namespace", req.Namespace, "err", err)
		return nil, errorHelper(err, fmt.Sprintf("failed to list pods for namespace %q", req.Namespace))
	}

	return pods, nil
}

func (s *K8SServer) GetNodeStats(ctx context.Context, req *pb.NodeRequest) (*pb.NodeStatsResponse, error) {
	endpoint := "get_node_stats"
	statusLabel := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, statusLabel).
			Inc()
	}()

	if req == nil {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "node stats request cannot be nil")
	}

	if req.NodeName == "" {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "node name is required")
	}

	node, err := s.NodeInfoService.GetNodeStats(ctx, req.NodeName)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("get node stats failed", "node", req.NodeName, "err", err)
		return nil, errorHelper(err, fmt.Sprintf("failed to get node stats for node %q", req.NodeName))

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
