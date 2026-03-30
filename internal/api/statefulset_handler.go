package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) ListStatefulSets(ctx context.Context, req *pb.NamespaceRequest) (*pb.WorkloadListResponse, error) {
	endpoint := "list_statefulsets"
	reqStatus := "success"
	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)
	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, reqStatus).
			Inc()
	}()

	if req == nil {
		reqStatus = "error"
		return nil, status.Error(codes.InvalidArgument, "statefulset list request cannot be nil")
	}

	workloads, err := s.StatefulSetService.List(ctx, req.Namespace, req.LabelSelector, "")
	if err != nil {
		reqStatus = "error"
		s.Logger.Error("list statefulsets failed", "namespace", req.Namespace, "labelSelector", req.LabelSelector, "err", err)
		return nil, errorHelper(err, "statefulset list")

	}

	return &pb.WorkloadListResponse{
		Workloads: workloads,
	}, nil
}
