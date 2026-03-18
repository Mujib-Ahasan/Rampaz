package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVClient struct {
	clientset kubernetes.Interface
}

func NewPVClient(clientset kubernetes.Interface) *PVClient {
	return &PVClient{
		clientset: clientset,
	}
}

func (c *PVClient) ListPVs(ctx context.Context) ([]corev1.PersistentVolume, error) {
	pvs, err := c.clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return pvs.Items, nil
}
