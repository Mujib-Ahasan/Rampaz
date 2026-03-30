package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) ListServices(ctx context.Context, req *pb.NamespaceRequest) (*pb.ServiceListResponse, error) {

	endpoint := "list_service"
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

	services, err := s.SVCService.ListServices(ctx, req.Namespace)
	if err != nil {
		reqStatus = "error"
		s.Logger.Error("list service failed", "namespace", req.Namespace, "err", err)
		return nil, errorHelper(err, "service list")

	}

	return &pb.ServiceListResponse{
		Services: services,
	}, nil
}
