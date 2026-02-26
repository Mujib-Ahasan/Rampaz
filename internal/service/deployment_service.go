package service

import (
	"context"
	"time"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentService struct {
	client *kubernetes.DeploymentClient
}

func NewDeploymentService(client *kubernetes.DeploymentClient) *DeploymentService {
	return &DeploymentService{client: client}
}

func (s *DeploymentService) ListDeployment(ctx context.Context, namespace string) ([]*pb.Workload, error) {
	deployments, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workload

	for _, d := range deployments {
		result = append(result, transformDeployment(&d))
	}

	return result, nil
}

func transformDeployment(d *appsv1.Deployment) *pb.Workload {
	var desired int32
	if d.Spec.Replicas != nil {
		desired = *d.Spec.Replicas
	}

	conditions := extractConditions(d.Status.Conditions)

	return &pb.Workload{
		Name:              d.Name,
		Namespace:         d.Namespace,
		DesiredReplicas:   desired,
		ReadyReplicas:     d.Status.ReadyReplicas,
		AvailableReplicas: d.Status.AvailableReplicas,
		UpdatedReplicas:   d.Status.UpdatedReplicas,
		Labels:            d.Labels,
		Conditions:        conditions,
		Owner:             extractOwner(d.OwnerReferences),
		Age:               calculateAge(d.CreationTimestamp.Time),
	}
}

func extractConditions(conds []appsv1.DeploymentCondition) []string {
	var result []string
	for _, c := range conds {
		result = append(result, string(c.Type)+"="+string(c.Status))
	}
	return result
}

func extractOwner(refs []metav1.OwnerReference) string {
	if len(refs) == 0 {
		return ""
	}
	return refs[0].Kind + "/" + refs[0].Name
}

func calculateAge(t time.Time) string {
	return time.Since(t).Round(time.Second).String()
}
