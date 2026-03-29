package api

import (
	"context"
	"fmt"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

func (s *K8SServer) GetWorkloadsByHealth(ctx context.Context, req *pb.WorkloadHealthRequest) (*pb.WorkloadListResponse, error) {
	if !isValidHealth(req.Health) {
		return nil, fmt.Errorf(
			"invalid health value %q (allowed: Healthy, Degraded, Unhealthy)",
			req.Health,
		)
	}
	wls, err := s.WorkloadService.GetWorkloadsByHealth(ctx, req.Namespace, req.Health)
	if err != nil {
		return nil, err
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
