package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/kmp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func syncK8sDeployment(k8sClient client.Client, ctx context.Context, local *appsv1.Deployment) error {
	logger := log.FromContext(ctx)
	logger.WithValues("name", local.GetName(), "kind", local.GetObjectKind())
	remote := &appsv1.Deployment{}
	if err := lookupDeployment(k8sClient, ctx, local, remote); err != nil {
		return fmt.Errorf("failed to lookup deployment %s: %w", local.GetName(), err)
	}

	// Candidate for create
	if remote.Name == "" {
		if err := k8sClient.Create(ctx, local); err != nil {
			return fmt.Errorf("failed to create resource: %w", err)

		}

		logger.Info("deployment created successfully")
		return nil
	}

	// Candidate for update
	if err := k8sClient.Update(ctx, local, client.DryRunAll); err != nil {
		return fmt.Errorf("failed to dry-run update resource: %w", err)
	}
	diff, err := kmp.SafeDiff(remote.Spec, local.Spec)
	if err != nil {
		return fmt.Errorf("failed to diff resurces: %w", err)
	}

	if diff != "" {
		logger.Info("submitted resource for update", "diff", diff)
		if err := k8sClient.Update(ctx, local); err != nil {
			return fmt.Errorf("failed to update resource: %w", err)
		}
	} else {
		logger.Info("no changes detected")
	}

	return nil
}

func syncK8sService(k8sClient client.Client, ctx context.Context, local *corev1.Service) error {
	logger := log.FromContext(ctx)
	logger.WithValues("name", local.GetName(), "kind", local.GetObjectKind())
	remote := &corev1.Service{}
	if err := lookupService(k8sClient, ctx, local, remote); err != nil {
		return fmt.Errorf("failed to lookup deployment %s: %w", local.GetName(), err)
	}

	// Candidate for create
	if remote.Name == "" {
		if err := k8sClient.Create(ctx, local); err != nil {
			return fmt.Errorf("failed to create resource: %w", err)

		}

		logger.Info("deployment created successfully")
		return nil
	}

	// Candidate for update
	if err := k8sClient.Update(ctx, local, client.DryRunAll); err != nil {
		return fmt.Errorf("failed to dry-run update resource: %w", err)
	}

	diff, err := kmp.SafeDiff(remote.Spec, local.Spec)
	if err != nil {
		return fmt.Errorf("failed to diff resurces: %w", err)
	}
	if diff != "" {
		logger.Info("submitted resource for update", "diff", diff)
		if err := k8sClient.Update(ctx, local); err != nil {
			return fmt.Errorf("failed to update resource: %w", err)
		}
	} else {
		logger.Info("no changes detected")
	}

	return nil
}

func syncK8sStatefulset(k8sClient client.Client, ctx context.Context, local *appsv1.StatefulSet) error {
	logger := log.FromContext(ctx)
	logger.WithValues("name", local.GetName(), "kind", local.GetObjectKind())
	remote := &appsv1.StatefulSet{}
	if err := lookupStatefulset(k8sClient, ctx, local, remote); err != nil {
		return fmt.Errorf("failed to lookup statefulset %s: %w", local.GetName(), err)
	}

	// Candidate for create
	if remote.Name == "" {
		if err := k8sClient.Create(ctx, local); err != nil {
			return fmt.Errorf("failed to create resource: %w", err)

		}

		logger.Info("statefulset created successfully")
		return nil
	}

	// Candidate for update
	if err := k8sClient.Update(ctx, local, client.DryRunAll); err != nil {
		return fmt.Errorf("failed to dry-run update resource: %w", err)
	}

	diff, err := kmp.SafeDiff(remote.Spec, local.Spec)
	if err != nil {
		return fmt.Errorf("failed to diff resurces: %w", err)
	}
	if diff != "" {
		logger.Info("submitted resource for update", "diff", diff)
		if err := k8sClient.Update(ctx, local); err != nil {
			return fmt.Errorf("failed to update resource: %w", err)
		}
	} else {
		logger.Info("no changes detected")
	}

	return nil
}

func lookupDeployment(k8sClient client.Client, ctx context.Context, local *appsv1.Deployment, remote *appsv1.Deployment) error {

	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(local), remote); err != nil {
		if apierrors.IsNotFound(err) {
			return nil

		} else {
			return fmt.Errorf("failed to lookup resource: %w", err)
		}

	}

	return nil
}

func lookupStatefulset(k8sClient client.Client, ctx context.Context, local *appsv1.StatefulSet, remote *appsv1.StatefulSet) error {

	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(local), remote); err != nil {
		if apierrors.IsNotFound(err) {
			return nil

		} else {
			return fmt.Errorf("failed to lookup resource: %w", err)
		}

	}

	return nil
}

func lookupService(k8sClient client.Client, ctx context.Context, local *corev1.Service, remote *corev1.Service) error {

	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(local), remote); err != nil {
		if apierrors.IsNotFound(err) {
			return nil

		} else {
			return fmt.Errorf("failed to lookup resource: %w", err)
		}

	}

	return nil
}
