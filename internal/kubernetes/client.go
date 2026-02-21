package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Clients struct {
	Kube    *kubernetes.Clientset
	Metrics *metricsv.Clientset
}

func NewClients() (*Clients, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", "/Users/mujibahasan/.kube/config")
	}

	kube, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	metrics, err := metricsv.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &Clients{
		Kube:    kube,
		Metrics: metrics,
	}, nil
}
