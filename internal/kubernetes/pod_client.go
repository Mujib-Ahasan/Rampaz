package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodClient struct {
	client *kubernetes.Clientset
}

func NewPodClient(c *kubernetes.Clientset) *PodClient {
	return &PodClient{client: c}
}

func (p *PodClient) ListPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	ps, err := p.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list pods across all namespaces: %w", err)
		}
		return nil, fmt.Errorf("list pods in namespace %q: %w", namespace, err)
	}

	return ps, nil

}
