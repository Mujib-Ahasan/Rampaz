package main

import (
	"context"
	"fmt"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *Server) GetNodeStats(ctx context.Context, req *pb.NodeRequest) (*pb.NodeStatsResponse, error) {
	node, err := s.clientset.CoreV1().Nodes().Get(ctx, req.NodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s: %v", req.NodeName, err)
	}

	cpu := node.Status.Capacity.Cpu().String()
	memory := node.Status.Capacity.Memory().String()

	res := &pb.NodeStatsResponse{
		Name:   node.Name,
		Cpu:    cpu,
		Memory: memory,
	}

	return res, nil

}
