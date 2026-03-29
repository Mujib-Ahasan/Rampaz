package service

import (
	"context"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

func (s *SummaryService) GetClusterOverview(ctx context.Context) (*pb.ClusterOverviewResponse, error) {
	pods, err := s.podClient.ListPods(ctx, "")
	if err != nil {
		return nil, err
	}

	deployments, err := s.deploymentClient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	replicasets, err := s.replicaSetclient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	statefulsets, err := s.statefulStateClient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	daemonsets, err := s.daemonSetClient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	jobs, err := s.jobClient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	cronjobs, err := s.cronJobClient.List(ctx, "", "")
	if err != nil {
		return nil, err
	}

	svcs, err := s.serviceClient.List(ctx, "")
	if err != nil {
		return nil, err
	}

	pvcs, err := s.pvcClinet.List(ctx, "")
	if err != nil {
		return nil, err
	}

	nps, err := s.networkPolicyClient.List(ctx, "")
	if err != nil {
		return nil, err
	}

	nodes, err := s.nodesClient.List(ctx)
	if err != nil {
		return nil, err
	}

	namespaces, err := s.namespaceClient.List(ctx)
	if err != nil {
		return nil, err
	}

	secrets, err := s.identityClient.ListSecrets(ctx, "")
	if err != nil {
		return nil, err
	}

	ingresses, err := s.identityClient.ListIngresses(ctx, "")
	if err != nil {
		return nil, err
	}

	configMaps, err := s.identityClient.ListConfigMaps(ctx, "")
	if err != nil {
		return nil, err
	}

	sas, err := s.identityClient.ListServiceAccounts(ctx, "")
	if err != nil {
		return nil, err
	}
	return &pb.ClusterOverviewResponse{
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
		Nodes:                  int32(len(nodes)),
		Namespaces:             int32(len(namespaces)),
		Secrets:                int32(len(secrets)),
		Serviceaccounts:        int32(len(sas)),
		Ingresses:              int32(len(ingresses)),
		Configmaps:             int32(len(configMaps)),
	}, nil
}
