package kubernetes

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type EventClient struct {
	client *kubernetes.Clientset
}

func NewEventClient(c *kubernetes.Clientset) *EventClient {
	return &EventClient{client: c}
}

func (e *EventClient) WatchEvents(ctx context.Context, namespace string) (<-chan *v1.Event, error) {

	watcher, err := e.client.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("watch events across all namespaces: %w", err)
		}
		return nil, fmt.Errorf("watch events in namespace %q: %w", namespace, err)
	}

	out := make(chan *v1.Event)

	go func() {
		defer close(out)
		defer watcher.Stop()

		for {
			select {

			case <-ctx.Done():
				return

			case evt, ok := <-watcher.ResultChan():
				if !ok {
					return
				}

				k8sEvent, ok := evt.Object.(*v1.Event)
				if ok {
					out <- k8sEvent
				}
			}
		}
	}()

	return out, nil
}
