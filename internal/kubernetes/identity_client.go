package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IdentityClient struct {
	client kubernetes.Interface
}

func NewIdentityClientClient(client kubernetes.Interface) *IdentityClient {
	return &IdentityClient{client: client}
}

func (c *IdentityClient) ListIngresses(ctx context.Context, namespace string) ([]networkingv1.Ingress, error) {
	ingresses, err := c.client.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingresses.Items, nil
}

func (c *IdentityClient) ListConfigMaps(ctx context.Context, namespace string) ([]corev1.ConfigMap, error) {
	configMaps, err := c.client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return configMaps.Items, nil
}

func (c *IdentityClient) ListSecrets(ctx context.Context, namespace string) ([]corev1.Secret, error) {
	secrets, err := c.client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return secrets.Items, nil
}

func (c *IdentityClient) ListServiceAccounts(ctx context.Context, namespace string) ([]corev1.ServiceAccount, error) {
	serviceAccounts, err := c.client.CoreV1().ServiceAccounts(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return serviceAccounts.Items, nil
}
