package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	appsv1 "k8s.io/api/apps/v1"
)

type DaemonSetService struct {
	client *kubernetes.DaemonSetClient
}

func NewDaemonSetService(client *kubernetes.DaemonSetClient) *DaemonSetService {
	return &DaemonSetService{client: client}
}

func (s *DaemonSetService) List(ctx context.Context, namespace string) ([]*pb.Workload, error) {
	daemonSets, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workload

	for _, ds := range daemonSets {
		result = append(result, transformDaemonSet(&ds))
	}

	return result, nil
}

func transformDaemonSet(ds *appsv1.DaemonSet) *pb.Workload {
	conditions := extractDaemonSetConditions(ds.Status.Conditions)

	return &pb.Workload{
		Name:              ds.Name,
		Namespace:         ds.Namespace,
		DesiredReplicas:   ds.Status.DesiredNumberScheduled,
		ReadyReplicas:     ds.Status.NumberReady,
		AvailableReplicas: ds.Status.NumberAvailable,
		UpdatedReplicas:   ds.Status.UpdatedNumberScheduled,
		Labels:            ds.Labels,
		Conditions:        conditions,
		Owner:             extractOwner(ds.OwnerReferences),
		Age:               calculateAge(ds.CreationTimestamp.Time),
	}
}

func extractDaemonSetConditions(conds []appsv1.DaemonSetCondition) []string {
	var result []string
	for _, c := range conds {
		result = append(result, string(c.Type)+"="+string(c.Status))
	}
	return result
}
