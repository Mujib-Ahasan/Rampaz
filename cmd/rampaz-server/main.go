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
	if err := api.StartGRPC(":50052", handler); err != nil {
		log.Fatal(err)
	}

}
