package api

import (
	"context"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

func (h *K8SServer) ListDaemonSets(ctx context.Context, req *pb.NamespaceRequest) (*pb.WorkloadListResponse, error) {

	workloads, err := h.DaemonSetService.List(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}

	return &pb.WorkloadListResponse{
		Workloads: workloads,
	}, nil
}
