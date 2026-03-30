package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	batchv1 "k8s.io/api/batch/v1"
)

type JobService struct {
	client *kubernetes.JobClient
}

func NewJobService(client *kubernetes.JobClient) *JobService {
	return &JobService{client: client}
}

func (s *JobService) List(ctx context.Context, namespace string, labelSelector string, health string) ([]*pb.Workload, error) {
	jobs, err := s.client.List(ctx, namespace, labelSelector)
	if err != nil {
		return nil, fmt.Errorf("fetch jobs data for workload response: %w", err)
	}

	var result []*pb.Workload

	for _, j := range jobs {
		w := transformJob(&j)
		if w.Health.String() == health || health == "" {
			result = append(result, w)
		}
	}

	return result, nil
}

func transformJob(job *batchv1.Job) *pb.Workload {
	conditions := extractJobConditions(job.Status.Conditions)

	return &pb.Workload{
		Name:       job.Name,
		Namespace:  job.Namespace,
		Active:     job.Status.Active,
		Succeeded:  job.Status.Succeeded,
		Failed:     job.Status.Failed,
		Labels:     job.Labels,
		Conditions: conditions,
		Owner:      extractOwner(job.OwnerReferences),
		Age:        calculateAge(job.CreationTimestamp.Time),
		Health:     computeJobHealth(job),
	}
}

func extractJobConditions(conds []batchv1.JobCondition) []string {
	var result []string
	for _, c := range conds {
		result = append(result, string(c.Type)+"="+string(c.Status))
	}
	return result
}
