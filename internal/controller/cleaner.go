package controller

import (
	"context"
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func cleanK8sObjects(k8sClient client.Client, ctx context.Context, objectsToClean []client.Object) error {
	var errs error
	logger := log.FromContext(ctx)

	for _, resource := range objectsToClean {
		logger.Info("deleting resource", "name", resource.GetName(), "namespace", resource.GetNamespace())
		if err := k8sClient.Delete(ctx, resource); err != nil {
			if apierrors.IsNotFound(err) {
				logger.Info("resource not found for deletion")
			} else {
				errs = errors.Join(errs, err)
			}
		}
	}

	return errs
}
