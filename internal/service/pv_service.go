package service

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type PVService struct {
	client *kubernetes.PVClient
}

func NewPVService(client *kubernetes.PVClient) *PVService {
	return &PVService{client: client}
}

func (c *PVService) ListPVs(ctx context.Context) ([]*pb.PVInfo, error) {
	pvs, err := c.client.ListPVs(ctx)
	if err != nil {
		return nil, err
	}

	var result []*pb.PVInfo

	for _, pv := range pvs {
		result = append(result, transformPV(&pv))
	}

	return result, nil
}

func transformPV(pv *corev1.PersistentVolume) *pb.PVInfo {
	var capacity string
	if qty, ok := pv.Spec.Capacity[corev1.ResourceStorage]; ok {
		capacity = qty.String()
	}

	var accessModes []string
	for _, mode := range pv.Spec.AccessModes {
		accessModes = append(accessModes, string(mode))
	}

	volumeMode := ""
	if pv.Spec.VolumeMode != nil {
		volumeMode = string(*pv.Spec.VolumeMode)
	}

	claimName := ""
	claimNamespace := ""
	if pv.Spec.ClaimRef != nil {
		claimName = pv.Spec.ClaimRef.Name
		claimNamespace = pv.Spec.ClaimRef.Namespace
	}

	return &pb.PVInfo{
		Name:           pv.Name,
		Phase:          string(pv.Status.Phase),
		StorageClass:   pv.Spec.StorageClassName,
		Capacity:       capacity,
		AccessModes:    accessModes,
		VolumeMode:     volumeMode,
		ClaimName:      claimName,
		ClaimNamespace: claimNamespace,
		Age:            calculateAge(pv.CreationTimestamp.Time),
	}
}
