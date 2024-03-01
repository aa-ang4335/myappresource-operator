package controller

import (
	"context"
	"fmt"

	"errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	myapigroupv1alpha1 "github.com/aa-ang4335/myappresource-operator/api/v1alpha1"
	"github.com/aa-ang4335/myappresource-operator/internal/service/podinfo"
	"github.com/aa-ang4335/myappresource-operator/internal/service/redis"
)

const controllerName = "controller.MyAppResource"

// MyAppResourceReconciler reconciles a MyAppResource object
type MyAppResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.api.group,resources=myappresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.api.group,resources=myappresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.api.group,resources=myappresources/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyAppResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *MyAppResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	o := &myapigroupv1alpha1.MyAppResource{}

	if err := r.Client.Get(ctx, req.NamespacedName, o); err != nil {
		if apierrors.IsNotFound(err) {

			// MyAppResource object not found. attempt to clean up any existing managed objects
			if err := cleanK8sObjects(r.Client, ctx, getAllManagedObjects(req.Name, req.Namespace, &o.Spec)); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to cleanup all managed objects: %w", err)
			}
			return ctrl.Result{}, nil
		}
		logger.Error(err, "failed to get myappresource object", "name", o.Name)
		return ctrl.Result{}, err
	}
	logger = logger.WithValues("namespace", o.Namespace, "name", o.Name, "controller", controllerName)
	// assume that status is always invalid
	o.Status.Valid = false

	var errs error
	// fetch objects to manage from the request
	redisStatefulSet := redis.GetStatefulset(req.Name, req.Namespace)
	redisService := redis.GetService(req.Name, req.Namespace)
	podinfoDeployment := podinfo.GetDeployment(req.Name, req.Namespace, redis.GetServiceAddr(req.Name, req.Namespace), &o.Spec)
	podinfoService := podinfo.GetService(req.Name, req.Namespace)

	// syncs redis objects if redis is enabled
	if o.Spec.Redis != nil && o.Spec.Redis.Enabled {
		logger.Info("initiating a sync for redis backend")
		if err := syncK8sStatefulset(r.Client, ctx, redisStatefulSet); err != nil {
			logger.Error(err, "failed to sync k8 statefulset", "name", redisStatefulSet.GetName())
			errs = errors.Join(errs, err)
		}

		if err := syncK8sService(r.Client, ctx, redisService); err != nil {
			logger.Error(err, "failed to sync k8 service", "name", redisService.GetName())
			errs = errors.Join(errs, err)

		}

	} else {
		// attempt to cleanup redis objects if the flag is unset
		if err := cleanK8sObjects(r.Client, ctx, []client.Object{
			redisStatefulSet,
			redisService,
		}); err != nil {
			logger.Error(err, "failed to cleanup redis objects")
			errs = errors.Join(errs, err)
		}
	}

	// syncs podinfo deployment object
	if err := syncK8sDeployment(r.Client, ctx, podinfoDeployment); err != nil {
		logger.Error(err, "failed to sync k8 deployment", "name", podinfoDeployment.GetName())
		errs = errors.Join(errs, err)
	}

	// syncs podinfor service object
	if err := syncK8sService(r.Client, ctx, podinfoService); err != nil {
		logger.Error(err, "failed to sync k8 service", "name", podinfoService.GetName())
		errs = errors.Join(errs, err)

	}

	if errs == nil {
		o.Status.Error = ""
		o.Status.Valid = true
	} else {
		o.Status.Error = errs.Error()
		o.Status.Valid = false
	}

	if err := r.Status().Update(ctx, o); err != nil {
		logger.Error(err, "failed to update the resource's status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyAppResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&myapigroupv1alpha1.MyAppResource{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}

// getAllManagedObjects returns all the k8s object directly managed by the operator.
func getAllManagedObjects(name string, namespace string, spec *myapigroupv1alpha1.MyAppResourceSpec) []client.Object {
	return []client.Object{
		redis.GetStatefulset(name, namespace),
		redis.GetService(name, namespace),
		podinfo.GetDeployment(name, namespace, redis.GetServiceAddr(name, namespace), spec),
		podinfo.GetService(name, namespace),
	}
}
