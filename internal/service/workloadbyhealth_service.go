package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"golang.org/x/sync/errgroup"
)

type WorkloadService struct {
	DeploymentService  *DeploymentService
	ReplicaSetService  *ReplicaSetService
	StatefulSetService *StatefulSetService
	DaemonSetService   *DaemonSetService
	JobService         *JobService
	CronJobService     *CronJobService
}

func NewWorkLoadService(deploymentService *DeploymentService,
	replicaSetService *ReplicaSetService,
	statefulSetService *StatefulSetService,
	daemonSetService *DaemonSetService,
	jobService *JobService,
	cronJobService *CronJobService,
) *WorkloadService {
	return &WorkloadService{
		DeploymentService:  deploymentService,
		StatefulSetService: statefulSetService,
		ReplicaSetService:  replicaSetService,
		DaemonSetService:   daemonSetService,
		JobService:         jobService,
		CronJobService:     cronJobService,
	}
}

func (s *WorkloadService) GetWorkloadsByHealth(ctx context.Context, namespace string, health string) ([]*pb.Workload, error) {
	var (
		mu        sync.Mutex
		workloads []*pb.Workload
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		items, err := s.DeploymentService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list deployments by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		items, err := s.ReplicaSetService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list replicasets by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		items, err := s.StatefulSetService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list statefulsets by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		items, err := s.DaemonSetService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list daemonsets by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		items, err := s.JobService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list jobs by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	g.Go(func() error {
		items, err := s.CronJobService.List(ctx, namespace, "", health)
		if err != nil {
			return fmt.Errorf("list cronjobs by health: %w", err)
		}

		mu.Lock()
		workloads = append(workloads, items...)
		mu.Unlock()
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return workloads, nil
}
