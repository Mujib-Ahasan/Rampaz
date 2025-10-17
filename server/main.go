package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Server struct {
	pb.UnimplementedK8SInfoServer
	clientset     *kubernetes.Clientset
	metricsClient *metricsv.Clientset
}

func main() {

	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Fatalf("Failed to load in-cluster config: %v", err)
	// }
	// config, err := clientcmd.BuildConfigFromFlags("", "C:\\Users\\korim\\.kube\\config")
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("err1")
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("err2")
		log.Fatal(err)
	}
	metcicsv, err := metricsv.NewForConfig(config)

	if err != nil {
		fmt.Println("err5")
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Println("err3")

		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterK8SInfoServer(s, &Server{clientset: clientset, metricsClient: metcicsv})

	fmt.Println("gRPC server running on :50051")
	if err := s.Serve(lis); err != nil {
		fmt.Println("err4")
		log.Fatal(err)
	}
}
