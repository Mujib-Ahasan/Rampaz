package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) GetPodStats(req *pb.PodRequest, stream grpc.ServerStreamingServer[pb.PodStatsResponse]) error {

	for {
		podMetricsList, err := s.metricsClient.MetricsV1beta1().PodMetricses(req.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to fetch pod metrics: %v", err)
		}

		for _, m := range podMetricsList.Items {

			var totalCPU, totalMem int64
			// i := 0
			for _, c := range m.Containers { // if there are any multipod containers then also this code gonna run fine
				// fmt.Printf("containers: %v \n", m.Containers)
				// fmt.Println(i)
				// i++
				cpuQuantity := c.Usage["cpu"]
				memQuantity := c.Usage["memory"]
				totalCPU += cpuQuantity.MilliValue()
				totalMem += memQuantity.Value() / (1024 * 1024) // bytes → MiB
			}

			res := &pb.PodStatsResponse{
				Name:      m.Name,
				Namespace: m.Namespace,
				Cpu:       strconv.FormatInt(totalCPU, 10) + "m",
				Memory:    strconv.FormatInt(totalMem, 10) + "Mi",
			}

			if err := stream.Send(res); err != nil {
				return fmt.Errorf("error sending stream: %v", err)
			}
		}

		time.Sleep(60 * time.Second)
	}
	return nil
}
