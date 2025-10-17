package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

func (s *Server) ListPods(ctx context.Context, req *pb.NamespaceRequest) (*pb.PodListResponse, error) {
	pods, err := s.clientset.CoreV1().Pods(req.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Println("error in server1")
		return nil, err
	}

	var podList []*pb.Pod
	for _, pod := range pods.Items {
		podList = append(podList, &pb.Pod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			NodeName:  pod.Spec.NodeName,
		})
	}

	return &pb.PodListResponse{Pods: podList}, nil
}
