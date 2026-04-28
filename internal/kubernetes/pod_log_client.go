package kubernetes

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type PodLogsClient struct {
	client *kubernetes.Clientset
}

func NewPodLogsClient(client *kubernetes.Clientset) *PodLogsClient {
	return &PodLogsClient{client: client}
}

func (c *PodLogsClient) GetPodLogs(ctx context.Context, namespace, podName string, tailLines int64) (string, error) {
	req := c.client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{TailLines: &tailLines})

	stream, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("stream error: get pod log for pod: %s, err: %w", podName, err)
	}
	defer stream.Close()

	data, err := io.ReadAll(stream)
	if err != nil {
		return "", fmt.Errorf("get pod log for pod: %s, err: %w", podName, err)
	}

	return string(data), nil
}
