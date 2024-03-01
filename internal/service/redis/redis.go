package redis

import (
	"fmt"

	"github.com/aa-ang4335/myappresource-operator/internal/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const version = "7.2.4"
const servicePort = 6379

// GetStatefulset retrieves a redis StatefulSet object based on the provided parameters.
//
// Parameters:
//
//	baseName: The base name of the statefulSet.
//	namespace: The namespace in which the StatefulSet resides.
//
// Returns:
//
//	*apps.StatefulSet: A pointer to the Kubernetes StatefulSet object or nil
func GetStatefulset(baseName string, namespace string) *appsv1.StatefulSet {

	replicas := int32(1)
	out := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(baseName),
			Namespace: namespace,
			Labels: utils.MergeLabels(utils.GenerateDefaultLabels(getName(baseName), namespace), map[string]string{
				"app.kubernetes.io/version": version,
			}),
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{

				MatchLabels: utils.GenerateDefaultLabels(getName(baseName), namespace),
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: utils.GenerateDefaultLabels(getName(baseName), namespace),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  getName(baseName),
							Image: fmt.Sprintf("docker.io/redis:%s", version),
							Ports: []corev1.ContainerPort{
								{
									Name:          "redis",
									ContainerPort: servicePort,
									Protocol:      "TCP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "redis-data",
									ReadOnly:  false,
									MountPath: "/data",
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"redis-cli ping",
										},
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"redis-cli ping",
										},
									},
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "File",
							ImagePullPolicy:          "IfNotPresent",
						},
					},
					RestartPolicy: "Always",
					DNSPolicy:     "ClusterFirst",
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					TypeMeta: metav1.TypeMeta{
						APIVersion: "v1",
						Kind:       "PersistentVolumeClaim",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:   "redis-data",
						Labels: utils.GenerateDefaultLabels(getName(baseName), namespace),
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						StorageClassName: utils.Ptr[string]("standard"),
						AccessModes: []corev1.PersistentVolumeAccessMode{
							"ReadWriteOnce",
						},
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("1Gi"),
							},
						},
					},
				},
			},
		},
	}

	return out
}

// GetService retrieves a redis k8s service object based on the provided parameters.
//
// Parameters:
//
//	baseName: The base name of the service.
//	namespace: The namespace in which the service lives.
//
// Returns:
//
//	*corev1.Service: A pointer to the Kubernetes Service object or nil.
func GetService(baseName string, namespace string) *corev1.Service {

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(baseName),
			Namespace: namespace,
			Labels:    utils.GenerateDefaultLabels(getName(baseName), namespace),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,

			Ports: []corev1.ServicePort{
				{Name: "redis", Protocol: "TCP", Port: servicePort},
			},
			Selector: utils.GenerateDefaultLabels(getName(baseName), namespace),
		},
	}

}

// GetServiceAddr returns the Kubernetes service address for Redis.
// The value returned will be passed as an environment variable for the podinfo cache server.
func GetServiceAddr(baseName string, namespace string) string {
	return fmt.Sprintf("tcp://%s.%s.svc.cluster.local:%d", getName(baseName), namespace, servicePort)
}

func getName(baseName string) string {
	return fmt.Sprintf("%s-redis", baseName)
}
