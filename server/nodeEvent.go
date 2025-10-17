package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) StreamEvents(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.EventResponse]) error {
	watcher, err := s.clientset.CoreV1().Events("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to start event watcher: %v", err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		if event.Type == watch.Error {
			log.Printf("⚠️  error event: %v", event)
			continue
		}

		k8sEvent, ok := event.Object.(*v1.Event)
		if !ok {
			continue
		}

		resp := &pb.EventResponse{
			Type:           string(event.Type),
			Reason:         k8sEvent.Reason,
			Message:        k8sEvent.Message,
			InvolvedObject: fmt.Sprintf("%s/%s", k8sEvent.InvolvedObject.Kind, k8sEvent.InvolvedObject.Name),
		}

		if err := stream.Send(resp); err != nil {
			log.Printf("error sending event to client: %v", err)
			return err
		}
	}

	return nil
}
