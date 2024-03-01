package podinfo

import (
	"fmt"

	myapigroupv1alpha1 "github.com/aa-ang4335/myappresource-operator/api/v1alpha1"
	"github.com/aa-ang4335/myappresource-operator/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const servicePort = 9898

// GetDeployment retrieves a k8s Deployment object based on the provided parameters.
//
// Parameters:
//
//	name: The name of the Deployment.
//	namespace: The namespace in which the Deployment lives.
//	redisServerAddr: The address of the Redis server.
//	spec: The MyAppResourceSpec containing specifications for the Deployment.
//
// Returns:
//
//	*appsv1.Deployment: A pointer to the k8s Deployment object or nil.
func GetDeployment(name string, namespace string, redisServerAddr string, spec *myapigroupv1alpha1.MyAppResourceSpec) *appsv1.Deployment {

	deployment := &appsv1.Deployment{}
	deployment.ObjectMeta = metav1.ObjectMeta{
		Name:      generateObjectName(name),
		Namespace: namespace,
		Labels:    utils.GenerateDefaultLabels(generateObjectName(name), namespace),
	}

	// return an empty deployment if image deployment not set.
	if spec.Image == nil || spec.Image.Repository == "" || spec.Image.Tag == "" {
		return deployment
	}

	envVarFromSpec := generateEnvVarForSpec(spec, redisServerAddr)
	containerResources := generateResourceRequirements(spec.Resources)

	deployment.Spec = appsv1.DeploymentSpec{

		Selector: &metav1.LabelSelector{
			MatchLabels: utils.GenerateDefaultLabels(generateObjectName(name), namespace),
		},
		Replicas: spec.ReplicaCount,
		Template: corev1.PodTemplateSpec{

			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-podinfo", name),
				Namespace: namespace,
				Labels:    utils.GenerateDefaultLabels(generateObjectName(name), namespace),
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:      fmt.Sprintf("%s-podinfo", name),
						Image:     fmt.Sprintf("%s:%s", spec.Image.Repository, spec.Image.Tag),
						Resources: containerResources,

						Ports: []corev1.ContainerPort{
							{
								Name:          "http",
								ContainerPort: servicePort,
								Protocol:      "TCP",
							},
						},

						Env: envVarFromSpec,
					},
				},
			},
		},
	}

	return deployment
}

// GetService retrieves podinfo k8s Service object based on the provided parameters.
//
// Parameters:
//
//	name: The name of the Service to retrieve.
//	namespace: The namespace in which the Service lives.
//
// Returns:
//
//	*corev1.Service: A pointer to the k8s Service object or nil.
func GetService(name string, namespace string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-podinfo", name),
			Namespace: namespace,
			Labels:    utils.GenerateDefaultLabels(generateObjectName(name), namespace),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,

			Selector: utils.GenerateDefaultLabels(generateObjectName(name), namespace),
			Ports: []corev1.ServicePort{
				{Name: "http", Protocol: "TCP", TargetPort: intstr.Parse("http"), Port: servicePort},
			},
		},
	}
}

// generateEnvVarForSpec generates container environment variables based on the provided MyAppResourceSpec and Redis server address.
// Environment variables are not set if there's no corresponding value in the spec.
// We also skip setting env var value for PODINFO_CACHE_SERVER if the Redis backend is disabled.
func generateEnvVarForSpec(spec *myapigroupv1alpha1.MyAppResourceSpec, redisServerAddr string) []corev1.EnvVar {
	var result []corev1.EnvVar

	if spec.UI != nil {
		if spec.UI.Color != "" {
			result = append(result, corev1.EnvVar{
				Name:  "PODINFO_UI_COLOR",
				Value: spec.UI.Color,
			})
		}
		if spec.UI.Message != "" {
			result = append(result, corev1.EnvVar{
				Name:  "PODINFO_UI_MESSAGE",
				Value: spec.UI.Message,
			})
		}

		if spec.Redis != nil {

			if spec.Redis.Enabled {
				result = append(result, corev1.EnvVar{
					Name:  "PODINFO_CACHE_SERVER",
					Value: redisServerAddr,
				})
			}
		}
	}

	return result
}

// generateResourceRequirements generates k8s resource requirements based on the provided resource specification.
// container resource requirements are set if CPURequest and MemoryLimit are non-empty.
func generateResourceRequirements(resourceSpec *myapigroupv1alpha1.Resources) corev1.ResourceRequirements {

	resourceRequirements := corev1.ResourceRequirements{}

	if resourceSpec == nil {
		return resourceRequirements
	}

	if resourceSpec.CPURequest != "" {
		resourceRequirements.Requests = corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse(resourceSpec.CPURequest),
		}
	}

	if resourceSpec.MemoryLimit != "" {
		resourceRequirements.Limits = corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse(resourceSpec.MemoryLimit),
		}
	}

	return resourceRequirements
}
func generateObjectName(baseName string) string {
	return fmt.Sprintf("%s-podinfo", baseName)
}
