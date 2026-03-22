package service

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type SVCService struct {
	client *kubernetes.ServiceClient
}

func NewSVCService(client *kubernetes.ServiceClient) *SVCService {
	return &SVCService{client: client}
}
func (s *SVCService) ListServices(ctx context.Context, namespace string) ([]*pb.ServiceInfo, error) {
	services, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.ServiceInfo

	for _, svc := range services {
		result = append(result, transformService(&svc))
	}

	return result, nil
}

func transformService(svc *corev1.Service) *pb.ServiceInfo {
	var externalIPs []string
	externalIPs = append(externalIPs, svc.Spec.ExternalIPs...)

	for _, ingress := range svc.Status.LoadBalancer.Ingress {
		if ingress.IP != "" {
			externalIPs = append(externalIPs, ingress.IP)
		}
		if ingress.Hostname != "" {
			externalIPs = append(externalIPs, ingress.Hostname)
		}
	}

	return &pb.ServiceInfo{
		Name:        svc.Name,
		Namespace:   svc.Namespace,
		Type:        string(svc.Spec.Type),
		ClusterIp:   svc.Spec.ClusterIP,
		ExternalIps: externalIPs,
		Ports:       extractServicePorts(svc.Spec.Ports),
		Age:         calculateAge(svc.CreationTimestamp.Time),
	}
}

func extractServicePorts(ports []corev1.ServicePort) []string {
	var result []string

	for _, p := range ports {
		port := fmt.Sprintf("%d/%s", p.Port, string(p.Protocol))

		if p.TargetPort.String() != "" {
			port = fmt.Sprintf("%s - %s", port, p.TargetPort.String())
		}

		if p.NodePort != 0 {
			port = fmt.Sprintf("%s(node:%d)", port, p.NodePort)
		}

		result = append(result, port)
	}

	return result
}
