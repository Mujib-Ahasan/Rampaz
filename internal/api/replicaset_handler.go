package api

import (
	"context"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

func (s *K8SServer) ListReplicaSets(ctx context.Context, req *pb.NamespaceRequest) (*pb.WorkloadListResponse, error) {

	workloads, err := s.ReplicaSetservice.ListReplicaSet(ctx, req.Namespace)
	if err != nil {
		return nil, err
	}

	return &pb.WorkloadListResponse{
		Workloads: workloads,
	}, nil
}
