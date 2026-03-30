package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) ListCronJobs(ctx context.Context, req *pb.NamespaceRequest) (*pb.WorkloadListResponse, error) {
	endpoint := "list_cronjobs"
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
		return nil, status.Error(codes.InvalidArgument, "cronjobs list request cannot be nil")
	}

	workloads, err := s.CronJobService.List(ctx, req.Namespace, req.LabelSelector, "")
	if err != nil {
		reqStatus = "error"
		s.Logger.Error("list cronjob failed", "namespace", req.Namespace, "labelSelector", req.LabelSelector, "err", err)
		return nil, errorHelper(err, "cronjob list")
	}

	return &pb.WorkloadListResponse{
		Workloads: workloads,
	}, nil
}
