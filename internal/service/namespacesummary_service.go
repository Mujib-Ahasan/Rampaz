package service

import (
	"context"

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
		return nil, err
	}

	deployments, err := s.deploymentClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	replicasets, err := s.replicaSetclient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	statefulsets, err := s.statefulStateClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	daemonsets, err := s.daemonSetClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	jobs, err := s.jobClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	cronjobs, err := s.cronJobClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	svcs, err := s.serviceClient.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	pvcs, err := s.pvcClinet.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	nps, err := s.networkPolicyClient.List(ctx, namespace)
	if err != nil {
		return nil, err
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
