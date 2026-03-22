package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type PVCService struct {
	client *kubernetes.PVCClient
}

func NewPVCService(client *kubernetes.PVCClient) *PVCService {
	return &PVCService{client: client}
}

func (c *PVCService) ListPVCs(ctx context.Context, namespace string) ([]*pb.PVCInfo, error) {
	pvcs, err := c.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.PVCInfo

	for _, pvc := range pvcs {
		result = append(result, transformPVC(&pvc))
	}

	return result, nil
}

func transformPVC(pvc *corev1.PersistentVolumeClaim) *pb.PVCInfo {
	var storage string
	if qty, ok := pvc.Spec.Resources.Requests[corev1.ResourceStorage]; ok {
		storage = qty.String()
	}

	var accessModes []string
	for _, mode := range pvc.Spec.AccessModes {
		accessModes = append(accessModes, string(mode))
	}

	var storageClass string
	if pvc.Spec.StorageClassName != nil {
		storageClass = *pvc.Spec.StorageClassName
	}

	return &pb.PVCInfo{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		Phase:            string(pvc.Status.Phase),
		StorageClass:     storageClass,
		AccessModes:      accessModes,
		RequestedStorage: storage,
		VolumeName:       pvc.Spec.VolumeName,
		Age:              calculateAge(pvc.CreationTimestamp.Time),
	}
}
