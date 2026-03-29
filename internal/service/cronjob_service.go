package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	batchv1 "k8s.io/api/batch/v1"
)

type CronJobService struct {
	client *kubernetes.CronJobClient
}

func NewCronJobService(client *kubernetes.CronJobClient) *CronJobService {
	return &CronJobService{client: client}
}

func (s *CronJobService) List(ctx context.Context, namespace string, labelSelector string, health string) ([]*pb.Workload, error) {
	cronJobs, err := s.client.List(ctx, namespace, labelSelector)
	if err != nil {
		return nil, err
	}

	var result []*pb.Workload

	for _, cj := range cronJobs {
		w := transformCronJob(&cj)
		if w.Health.String() == health || health == "" {
			result = append(result, w)
		}
	}

	return result, nil
}

func transformCronJob(cj *batchv1.CronJob) *pb.Workload {
	var lastSchedule string
	if cj.Status.LastScheduleTime != nil {
		lastSchedule = cj.Status.LastScheduleTime.Time.String()
	}

	return &pb.Workload{
		Name:             cj.Name,
		Namespace:        cj.Namespace,
		Schedule:         cj.Spec.Schedule,
		LastScheduleTime: lastSchedule,
		Active:           int32(len(cj.Status.Active)),
		Labels:           cj.Labels,
		Owner:            extractOwner(cj.OwnerReferences),
		Age:              calculateAge(cj.CreationTimestamp.Time),
		Health:           computeCronJobHealth(cj),
	}
}
