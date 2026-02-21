package service

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	v1 "k8s.io/api/core/v1"
)

type EventService struct {
	client *kubernetes.EventClient
}

func NewEventService(c *kubernetes.EventClient) *EventService {
	return &EventService{client: c}
}

func (s *EventService) StreamEvents(ctx context.Context) (<-chan *v1.Event, error) {
	return s.client.WatchEvents(ctx)
}
