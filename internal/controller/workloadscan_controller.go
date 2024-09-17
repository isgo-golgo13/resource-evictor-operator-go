package controller

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WorkloadScanReconciler reconciles a WorkloadScan object
type WorkloadScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile function that watches for Deployments and StatefulSets
func (r *WorkloadScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// List Deployments
	deployments := &v1.DeploymentList{}
	if err := r.List(ctx, deployments, &client.ListOptions{}); err != nil {
		log.Error(err, "unable to list deployments")
		return ctrl.Result{}, err
	}

	// List StatefulSets
	statefulSets := &v1.StatefulSetList{}
	if err := r.List(ctx, statefulSets, &client.ListOptions{}); err != nil {
		log.Error(err, "unable to list statefulsets")
		return ctrl.Result{}, err
	}

	// Process Deployments
	for _, deployment := range deployments.Items {
		if !isValidResource(deployment.Spec.Template.Spec.Containers) {
			log.Info(fmt.Sprintf("Evicting deployment: %s", deployment.Name))
			if err := r.Delete(ctx, &deployment, &client.DeleteOptions{}); err != nil {
				log.Error(err, "unable to delete deployment")
				return ctrl.Result{}, err
			}
		}
	}

	// Process StatefulSets
	for _, statefulSet := range statefulSets.Items {
		if !isValidResource(statefulSet.Spec.Template.Spec.Containers) {
			log.Info(fmt.Sprintf("Evicting statefulset: %s", statefulSet.Name))
			if err := r.Delete(ctx, &statefulSet, &client.DeleteOptions{}); err != nil {
				log.Error(err, "unable to delete statefulset")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// Helper function to validate resources
func isValidResource(containers []corev1.Container) bool {
	for _, container := range containers {
		if container.Resources.Limits == nil || container.Resources.Limits.Cpu().IsZero() || container.Resources.Limits.Memory().IsZero() {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *WorkloadScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Deployment{}).
		For(&v1.StatefulSet{}).
		Complete(r)
}
