package controllers

import (
	"context"
	"fmt"
	"math"
	"time"

	"0xhub/operator/api/v1"
	"0xhub/operator/internal/backend"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// Base retry delay
	baseRetryDelay = 30 * time.Second
	// Maximum retry delay (5 minutes)
	maxRetryDelay = 5 * time.Minute
	// Maximum retry count before giving up (will still retry but with max delay)
	maxRetryCount = 10
)

// ProjectReconciler reconciles a Project object
type ProjectReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	BackendClient *backend.Client
}

//+kubebuilder:rbac:groups=hub.0xhub.io,resources=projects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hub.0xhub.io,resources=projects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hub.0xhub.io,resources=projects/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Project instance
	project := &v1.Project{}
	if err := r.Get(ctx, req.NamespacedName, project); err != nil {
		// Project not found, might have been deleted
		if client.IgnoreNotFound(err) == nil {
			logger.Info("Project resource not found, checking if it needs to be deleted from backend")
			// Try to delete from backend using the resource name as ID
			_ = r.BackendClient.DeleteProject(req.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Check if the project is being deleted
	if !project.DeletionTimestamp.IsZero() {
		logger.Info("Project is being deleted, removing from backend", "project", req.Name)
		if err := r.BackendClient.DeleteProject(req.Name); err != nil {
			logger.Error(err, "Failed to delete project from backend", "project", req.Name)
			// Update status with error and retry info
			retryDelay := r.calculateRetryDelay(project.Status.RetryCount)
			project.Status.Error = fmt.Sprintf("Failed to delete: %v", err)
			project.Status.Synced = false
			project.Status.RetryCount++
			now := time.Now()
			project.Status.LastRetryAt = &metav1.Time{Time: now}
			if updateErr := r.Status().Update(ctx, project); updateErr != nil {
				return ctrl.Result{}, updateErr
			}
			logger.Info("Will retry deletion", "project", req.Name, "retryCount", project.Status.RetryCount, "retryAfter", retryDelay)
			return ctrl.Result{RequeueAfter: retryDelay}, nil
		}
		// Reset retry count on success
		project.Status.RetryCount = 0
		project.Status.Error = ""
		if updateErr := r.Status().Update(ctx, project); updateErr != nil {
			logger.Error(updateErr, "Failed to update status after successful delete")
		}
		return ctrl.Result{}, nil
	}

	// Convert CRD Project to backend Project
	backendProject := &backend.Project{
		ID:          req.Name, // Use Kubernetes resource name as ID
		Name:        project.Spec.Name,
		Description: project.Spec.Description,
		URL:         project.Spec.URL,
		Icon:        project.Spec.Icon,
		Category:    project.Spec.Category,
		Status:      project.Spec.Status,
	}

	// Check if project exists in backend
	existingProject, err := r.BackendClient.GetProject(req.Name)
	if err != nil {
		// Project doesn't exist, create it
		logger.Info("Creating project in backend", "project", req.Name)
		if err := r.BackendClient.CreateProject(backendProject); err != nil {
			logger.Error(err, "Failed to create project in backend", "project", req.Name, "retryCount", project.Status.RetryCount)
			// Update status with error and retry info
			retryDelay := r.calculateRetryDelay(project.Status.RetryCount)
			project.Status.Error = fmt.Sprintf("Failed to create: %v", err)
			project.Status.Synced = false
			project.Status.RetryCount++
			now := time.Now()
			project.Status.LastRetryAt = &metav1.Time{Time: now}
			if updateErr := r.Status().Update(ctx, project); updateErr != nil {
				return ctrl.Result{}, updateErr
			}
			logger.Info("Will retry creation", "project", req.Name, "retryCount", project.Status.RetryCount, "retryAfter", retryDelay)
			return ctrl.Result{RequeueAfter: retryDelay}, nil
		}
	} else {
		// Project exists, check if update is needed
		needsUpdate := existingProject.Name != backendProject.Name ||
			existingProject.Description != backendProject.Description ||
			existingProject.URL != backendProject.URL ||
			existingProject.Icon != backendProject.Icon ||
			existingProject.Category != backendProject.Category ||
			existingProject.Status != backendProject.Status

		if needsUpdate {
			logger.Info("Updating project in backend", "project", req.Name)
			if err := r.BackendClient.UpdateProject(req.Name, backendProject); err != nil {
				logger.Error(err, "Failed to update project in backend", "project", req.Name, "retryCount", project.Status.RetryCount)
				// Update status with error and retry info
				retryDelay := r.calculateRetryDelay(project.Status.RetryCount)
				project.Status.Error = fmt.Sprintf("Failed to update: %v", err)
				project.Status.Synced = false
				project.Status.RetryCount++
				now := time.Now()
				project.Status.LastRetryAt = &metav1.Time{Time: now}
				if updateErr := r.Status().Update(ctx, project); updateErr != nil {
					return ctrl.Result{}, updateErr
				}
				logger.Info("Will retry update", "project", req.Name, "retryCount", project.Status.RetryCount, "retryAfter", retryDelay)
				return ctrl.Result{RequeueAfter: retryDelay}, nil
			}
		} else {
			logger.Info("Project already in sync", "project", req.Name)
		}
	}

	// Update status to indicate successful sync
	now := time.Now()
	project.Status.Synced = true
	project.Status.LastSyncedAt = &metav1.Time{Time: now}
	project.Status.Error = ""
	// Reset retry count on success
	project.Status.RetryCount = 0
	project.Status.LastRetryAt = nil

	if err := r.Status().Update(ctx, project); err != nil {
		logger.Error(err, "Failed to update project status", "project", req.Name)
		return ctrl.Result{}, err
	}

	logger.Info("Successfully synced project to backend", "project", req.Name)
	return ctrl.Result{}, nil
}

// calculateRetryDelay calculates exponential backoff delay with jitter
func (r *ProjectReconciler) calculateRetryDelay(retryCount int) time.Duration {
	if retryCount >= maxRetryCount {
		return maxRetryDelay
	}

	// Exponential backoff: baseDelay * 2^retryCount
	delay := float64(baseRetryDelay) * math.Pow(2, float64(retryCount))
	if delay > float64(maxRetryDelay) {
		delay = float64(maxRetryDelay)
	}

	return time.Duration(delay)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Project{}).
		Complete(r)
}

