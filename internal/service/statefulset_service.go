package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	appsv1 "k8s.io/api/apps/v1"
)

type StatefulSetService struct {
	client *kubernetes.StatefulSetClient
}

func NewStatefulSetService(client *kubernetes.StatefulSetClient) *StatefulSetService {
	return &StatefulSetService{client: client}
}

func (s *StatefulSetService) List(ctx context.Context, namespace string) ([]*pb.Workload, error) {
	statefulSets, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workload

	for _, ss := range statefulSets {
		result = append(result, transformStatefulSet(&ss))
	}

	return result, nil
}

func transformStatefulSet(ss *appsv1.StatefulSet) *pb.Workload {
	var desired int32
	if ss.Spec.Replicas != nil {
		desired = *ss.Spec.Replicas
	}

	conditions := extractStatefulSetConditions(ss.Status.Conditions)

	return &pb.Workload{
		Name:              ss.Name,
		Namespace:         ss.Namespace,
		DesiredReplicas:   desired,
		ReadyReplicas:     ss.Status.ReadyReplicas,
		AvailableReplicas: ss.Status.CurrentReplicas,
		UpdatedReplicas:   ss.Status.UpdatedReplicas,
		Labels:            ss.Labels,
		Conditions:        conditions,
		Owner:             extractOwner(ss.OwnerReferences),
		Age:               calculateAge(ss.CreationTimestamp.Time),
	}
}

func extractStatefulSetConditions(conds []appsv1.StatefulSetCondition) []string {
	var result []string
	for _, c := range conds {
		result = append(result, string(c.Type)+"="+string(c.Status))
	}
	return result
}
