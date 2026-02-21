package main

import (
	"fmt"
	"log"

	"github.com/Mujib-Ahasan/Rampaz/internal/api"
	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	"github.com/Mujib-Ahasan/Rampaz/internal/service"
)

func main() {
	clients, err := kubernetes.NewClients()
	if err != nil {
		log.Fatal(err)
	}

	podClient := kubernetes.NewPodClient(clients.Kube)
	podService := service.NewPodService(podClient)

	nodeClent := kubernetes.NewNodeClient(clients.Kube)
	nodeService := service.NewNodeService(nodeClent)

	eventClient := kubernetes.NewEventClient(clients.Kube)
	eventService := service.NewEventService(eventClient)

	podMetricsClient := kubernetes.NewPodMetricsClient(clients.Metrics)
	podMetService := service.NewPodMetService(podMetricsClient)

	nodeMetricsClient := kubernetes.NewNodeMetricsClient(clients.Metrics)
	nodeMetService := service.NewNodeMetService(nodeMetricsClient)

	handler := &api.K8SServer{
		PodService:     podService,
		NodeService:    nodeService,
		EventService:   eventService,
		PodMetService:  podMetService,
		NodeMetService: nodeMetService,
	}

	fmt.Printf("server is running on port: 50052 \n")

	// temporary comments:
	// res, err := handler.GetPods(context.Background(), &proto.PodRequest{Namespace: "default"})
	// if err != nil {
	// 	log.Fatalf("some error: please check it %v \n", err)
	// }

	// res1, err := handler.GetNodeInfo(context.Background(), &proto.NodeRequest{NodeName: "minikube"})
	// if err != nil {
	// 	log.Fatalf("some error: please check it %v\n", err)
	// }

	// fmt.Printf("%+v\n", res)
	// fmt.Printf("%+v\n", res1)
	// temporary end

	if err := api.StartGRPC(":50052", handler); err != nil {
		log.Fatal(err)
	}
}
