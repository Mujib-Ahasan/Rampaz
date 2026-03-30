package service

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type SummaryService struct {
	podClient           *kubernetes.PodClient
	deploymentClient    *kubernetes.DeploymentClient
	replicaSetclient    *kubernetes.ReplicaSetClient
	statefulStateClient *kubernetes.StatefulSetClient
	daemonSetClient     *kubernetes.DaemonSetClient
	jobClient           *kubernetes.JobClient
	cronJobClient       *kubernetes.CronJobClient
	serviceClient       *kubernetes.ServiceClient
	pvcClinet           *kubernetes.PVCClient
	networkPolicyClient *kubernetes.NetworkPolicyClient
	nodesClient         *kubernetes.NodeClient
	namespaceClient     *kubernetes.NamespaceClient
	identityClient      *kubernetes.IdentityClient
}

func NewSummaryService(
	podClient *kubernetes.PodClient,
	deploymentClient *kubernetes.DeploymentClient,
	replicaSetclient *kubernetes.ReplicaSetClient,
	statefulStateClient *kubernetes.StatefulSetClient,
	daemonSetClient *kubernetes.DaemonSetClient,
	jobClient *kubernetes.JobClient,
	cronJobClient *kubernetes.CronJobClient,
	serviceClient *kubernetes.ServiceClient,
	pvcClinet *kubernetes.PVCClient,
	networkPolicyClient *kubernetes.NetworkPolicyClient,
	nodeclinet *kubernetes.NodeClient,
	namespaceClient *kubernetes.NamespaceClient,
	identityClient *kubernetes.IdentityClient,
) *SummaryService {
	return &SummaryService{
		podClient:           podClient,
		deploymentClient:    deploymentClient,
		replicaSetclient:    replicaSetclient,
		statefulStateClient: statefulStateClient,
		daemonSetClient:     daemonSetClient,
		jobClient:           jobClient,
		cronJobClient:       cronJobClient,
		serviceClient:       serviceClient,
		pvcClinet:           pvcClinet,
		networkPolicyClient: networkPolicyClient,
		nodesClient:         nodeclinet,
		namespaceClient:     namespaceClient,
		identityClient:      identityClient,
	}

}

func (s *SummaryService) GetNamespaceSummary(ctx context.Context, namespace string) (*pb.NamespaceSummaryResponse, error) {
	pods, err := s.podClient.ListPods(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch pods data for namespacesummary response: %w", err)
	}

	deployments, err := s.deploymentClient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch deployments data for namespacesummary response: %w", err)
	}

	replicasets, err := s.replicaSetclient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch replicasets data for namespacesummary response: %w", err)
	}

	statefulsets, err := s.statefulStateClient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch statefulsets data for namespacesummary response: %w", err)
	}

	daemonsets, err := s.daemonSetClient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch daemonsets data for namespacesummary response: %w", err)
	}

	jobs, err := s.jobClient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch jobs data for namespacesummary response: %w", err)
	}

	cronjobs, err := s.cronJobClient.List(ctx, namespace, "")
	if err != nil {
		return nil, fmt.Errorf("fetch cronjobs data for namespacesummary response: %w", err)
	}

	svcs, err := s.serviceClient.List(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch services data for namespacesummary response: %w", err)
	}

	pvcs, err := s.pvcClinet.List(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch pvcs data for namespacesummary response: %w", err)
	}

	nps, err := s.networkPolicyClient.List(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("fetch networkpolicies data for namespacesummary response: %w", err)
	}

	return &pb.NamespaceSummaryResponse{
		Namespace:              namespace,
		Pods:                   int32(len(pods.Items)),
		Deployments:            int32(len(deployments)),
		Replicasets:            int32(len(replicasets)),
		Statefulsets:           int32(len(statefulsets)),
		Daemonsets:             int32(len(daemonsets)),
		Jobs:                   int32(len(jobs)),
		Cronjobs:               int32(len(cronjobs)),
		Services:               int32(len(svcs)),
		PersistentVolumeClaims: int32(len(pvcs)),
		NetworkPolicies:        int32(len(nps)),
	}, nil
}
