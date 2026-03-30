package api

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) GetWorkloadsByHealth(ctx context.Context, req *pb.WorkloadHealthRequest) (*pb.WorkloadListResponse, error) {
	endpoint := "get_workloads_by_health"
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
		return nil, status.Error(codes.InvalidArgument, "workload health request cannot be nil")
	}

	if !isValidHealth(req.Health) {
		statusLabel = "error"
		return nil, status.Errorf(codes.InvalidArgument, "invalid health value %q (allowed: HEALTHY, DEGRADED, UNHEALTHY)", req.Health)
	}

	wls, err := s.WorkloadService.GetWorkloadsByHealth(ctx, req.Namespace, req.Health)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("get workloads by health failed", "namespace", req.Namespace, "health", req.Health, "err", err)
		return nil, errorHelper(err, fmt.Sprintf("get workloads by health %q in namespace %q", req.Health, req.Namespace))

	}

	return &pb.WorkloadListResponse{
		Workloads: wls,
	}, nil
}

func isValidHealth(h string) bool {
	switch h {
	case "HEALTHY", "DEGRADED", "UNHEALTHY":
		return true
	default:
		return false
	}
}
