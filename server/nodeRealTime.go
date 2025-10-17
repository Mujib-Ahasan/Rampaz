package main

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) GetNodeRealTimeStats(ctx context.Context, req *pb.NodeRequest) (*pb.NodeStatsResponse, error) {
	metrics, err := s.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, req.NodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics for node %s: %v", req.NodeName, err)
	}

	cpu := (metrics.Usage.Cpu().MilliValue())
	memory := metrics.Usage.Memory().Value() / (1024 * 1024)

	res := &pb.NodeStatsResponse{
		Name:   req.NodeName,
		Cpu:    strconv.FormatInt(cpu, 10) + "m",
		Memory: strconv.FormatInt(memory, 10) + "Mi",
	}

	return res, nil
}
