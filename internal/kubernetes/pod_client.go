package kubernetes

import (
	"context"

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
	return p.client.CoreV1().
		Pods(namespace).
		List(ctx, metav1.ListOptions{})
}
