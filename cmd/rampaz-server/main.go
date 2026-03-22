package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mujib-Ahasan/Rampaz/internal/api"
	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	"github.com/Mujib-Ahasan/Rampaz/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	metrics.Init()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9090", nil)
	}()

	clients, err := kubernetes.NewClients()
	if err != nil {
		log.Fatal(err)
	}

	podClient := kubernetes.NewPodClient(clients.Kube)
	podService := service.NewPodService(podClient)

	nodeInfoClent := kubernetes.NewNodeInfoClient(clients.Kube)
	nodeInfoService := service.NewNodeInfoService(nodeInfoClent)

	eventClient := kubernetes.NewEventClient(clients.Kube)
	eventService := service.NewEventService(eventClient)

	podMetricsClient := kubernetes.NewPodMetricsClient(clients.Metrics)
	podMetService := service.NewPodMetService(podMetricsClient)

	nodeMetricsClient := kubernetes.NewNodeMetricsClient(clients.Metrics)
	nodeMetService := service.NewNodeMetService(nodeMetricsClient)

	deploymentClient := kubernetes.NewDeploymentClient(clients.Kube)
	deploymentService := service.NewDeploymentService(deploymentClient)

	replicaSetClient := kubernetes.NewReplicaSetClient(clients.Kube)
	replicaSetService := service.NewReplicaSetService(replicaSetClient)

	daemonSetClient := kubernetes.NewDaemonSetClient(clients.Kube)
	daemonSetService := service.NewDaemonSetService(daemonSetClient)

	statefulSetClient := kubernetes.NewStatefulSetClient(clients.Kube)
	statefulSetService := service.NewStatefulSetService(statefulSetClient)

	jobClient := kubernetes.NewJobClient(clients.Kube)
	jobService := service.NewJobService(jobClient)

	cronJobClient := kubernetes.NewCronJobClient(clients.Kube)
	cronJobService := service.NewCronJobService(cronJobClient)

	serviceClient := kubernetes.NewServiceClient(clients.Kube)
	svcService := service.NewSVCService(serviceClient)

	nameSpaceClient := kubernetes.NewNamespaceClient(clients.Kube)
	namespaceService := service.NewNamespaceService(nameSpaceClient)

	pvcClient := kubernetes.NewPVCClient(clients.Kube)
	pvcService := service.NewPVCService(pvcClient)

	pvClient := kubernetes.NewPVClient(clients.Kube)
	pvService := service.NewPVService(pvClient)

	nodeClient := kubernetes.NewNodeClient(clients.Kube)
	nodeService := service.NewNodeService(nodeClient)

	networkPolicyClient := kubernetes.NewNetworkPolicyClient(clients.Kube)
	networkPolicyservice := service.NewNetworkPolicyService(networkPolicyClient)

	namespaceSummaryService := service.NewNamespaceSummaryService(
		podClient,
		deploymentClient,
		replicaSetClient,
		statefulSetClient,
		daemonSetClient,
		jobClient,
		cronJobClient,
		serviceClient,
		pvcClient,
		networkPolicyClient)

	handler := &api.K8SServer{
		PodService:              podService,
		NodeInfoService:         nodeInfoService,
		EventService:            eventService,
		PodMetService:           podMetService,
		NodeMetService:          nodeMetService,
		DeploymentService:       deploymentService,
		ReplicaSetservice:       replicaSetService,
		DaemonSetService:        daemonSetService,
		StatefulSetService:      statefulSetService,
		JobService:              jobService,
		CronJobService:          cronJobService,
		SVCService:              svcService,
		NamespaceService:        namespaceService,
		PVCService:              pvcService,
		PVService:               pvService,
		NodeService:             nodeService,
		NetworkPolicyService:    networkPolicyservice,
		NamespaceSummaryService: namespaceSummaryService,
	}

	fmt.Printf("server is running on port: 50052 \n")
	if err := api.StartGRPC(":50052", handler); err != nil {
		log.Fatal(err)
	}

}
