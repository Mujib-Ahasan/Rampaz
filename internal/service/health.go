package service

import (
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

// Healthy
//1. readyReplicas == desiredReplicas
//2. generation is observed
//3. Available=True
//4. Progressing=True or rollout complete
// Degraded
//1. some replicas are ready, but not all
//2. or, rollout is still progressing
// Unhealthy
//1. zero ready replicas while replicas are desired
//2. or Available=False
//3. or progress deadline exceeded

func computeDeploymentHealth(d *appsv1.Deployment) pb.HealthStatus {
	desired := int32(1)
	if d.Spec.Replicas != nil {
		desired = *d.Spec.Replicas
	}

	ready := d.Status.ReadyReplicas

	var availableCond, progressingCond appsv1.DeploymentCondition
	var hasAvailable, hasProgressing bool

	for _, cond := range d.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			availableCond = cond
			hasAvailable = true
		}
		if cond.Type == appsv1.DeploymentProgressing {
			progressingCond = cond
			hasProgressing = true
		}
	}

	if desired == 0 {
		return pb.HealthStatus_HEALTHY
	}

	if ready == desired && d.Status.ObservedGeneration >= d.Generation && (!hasAvailable || availableCond.Status == corev1.ConditionTrue) {
		return pb.HealthStatus_HEALTHY
	}

	if hasProgressing && progressingCond.Reason == "ProgressDeadlineExceeded" {
		return pb.HealthStatus_UNHEALTHY
	}

	if ready == 0 {
		return pb.HealthStatus_UNHEALTHY
	}

	return pb.HealthStatus_DEGRADED
}

func computeReplicaSetHealth(rs *appsv1.ReplicaSet) pb.HealthStatus {
	desired := int32(1)
	if rs.Spec.Replicas != nil {
		desired = *rs.Spec.Replicas
	}

	ready := rs.Status.ReadyReplicas

	if desired == 0 {
		return pb.HealthStatus_HEALTHY
	}

	if ready == desired && rs.Status.ObservedGeneration >= rs.Generation {
		return pb.HealthStatus_HEALTHY
	}

	if ready == 0 {
		return pb.HealthStatus_UNHEALTHY
	}

	return pb.HealthStatus_DEGRADED
}

func computeStatefulSetHealth(sts *appsv1.StatefulSet) pb.HealthStatus {
	desired := int32(1)
	if sts.Spec.Replicas != nil {
		desired = *sts.Spec.Replicas
	}

	ready := sts.Status.ReadyReplicas
	updated := sts.Status.UpdatedReplicas

	if desired == 0 {
		return pb.HealthStatus_HEALTHY
	}

	if ready == desired && updated == desired && sts.Status.ObservedGeneration >= sts.Generation {
		return pb.HealthStatus_HEALTHY
	}

	if ready == 0 {
		return pb.HealthStatus_UNHEALTHY
	}

	return pb.HealthStatus_DEGRADED
}

func computeDaemonSetHealth(ds *appsv1.DaemonSet) pb.HealthStatus {
	desired := ds.Status.DesiredNumberScheduled
	ready := ds.Status.NumberReady
	available := ds.Status.NumberAvailable

	if desired == 0 {
		return pb.HealthStatus_HEALTHY
	}

	if ready == desired && available == desired && ds.Status.ObservedGeneration >= ds.Generation {
		return pb.HealthStatus_HEALTHY
	}

	if ready == 0 {
		return pb.HealthStatus_UNHEALTHY
	}

	return pb.HealthStatus_DEGRADED
}

func computeJobHealth(job *batchv1.Job) pb.HealthStatus {
	for _, cond := range job.Status.Conditions {
		if cond.Type == batchv1.JobComplete && cond.Status == corev1.ConditionTrue {
			return pb.HealthStatus_HEALTHY
		}
		if cond.Type == batchv1.JobFailed && cond.Status == corev1.ConditionTrue {
			return pb.HealthStatus_UNHEALTHY
		}
	}

	if job.Status.Active > 0 {
		return pb.HealthStatus_DEGRADED
	}

	if job.Status.Failed > 0 {
		return pb.HealthStatus_UNHEALTHY
	}

	return pb.HealthStatus_DEGRADED
}

func computeCronJobHealth(cj *batchv1.CronJob) pb.HealthStatus {
	if cj.Spec.Suspend != nil && *cj.Spec.Suspend {
		return pb.HealthStatus_DEGRADED
	}

	if cj.Status.LastScheduleTime == nil {
		return pb.HealthStatus_DEGRADED
	}

	return pb.HealthStatus_HEALTHY
}
