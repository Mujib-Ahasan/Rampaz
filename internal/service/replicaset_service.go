package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	appsv1 "k8s.io/api/apps/v1"
)

type ReplicaSetService struct {
	client *kubernetes.ReplicaSetClient
}

func NewReplicaSetService(client *kubernetes.ReplicaSetClient) *ReplicaSetService {
	return &ReplicaSetService{client: client}
}

func (s *ReplicaSetService) ListReplicaSet(ctx context.Context, namespace string) ([]*pb.Workload, error) {
	replicaSets, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workload

	for _, rs := range replicaSets {
		result = append(result, transformReplicaSet(&rs))
	}

	return result, nil
}

func transformReplicaSet(rs *appsv1.ReplicaSet) *pb.Workload {
	var desired int32
	if rs.Spec.Replicas != nil {
		desired = *rs.Spec.Replicas
	}

	conditions := extractReplicaSetConditions(rs.Status.Conditions)

	return &pb.Workload{
		Name:              rs.Name,
		Namespace:         rs.Namespace,
		DesiredReplicas:   desired,
		ReadyReplicas:     rs.Status.ReadyReplicas,
		AvailableReplicas: rs.Status.AvailableReplicas,
		Labels:            rs.Labels,
		Conditions:        conditions,
		Owner:             extractOwner(rs.OwnerReferences),
		Age:               calculateAge(rs.CreationTimestamp.Time),
	}
}

func extractReplicaSetConditions(conds []appsv1.ReplicaSetCondition) []string {
	var result []string
	for _, c := range conds {
		result = append(result, string(c.Type)+"="+string(c.Status))
	}
	return result
}
