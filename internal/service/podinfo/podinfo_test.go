package podinfo

import (
	"testing"

	myapigroupv1alpha1 "github.com/aa-ang4335/myappresource-operator/api/v1alpha1"
	"github.com/aa-ang4335/myappresource-operator/internal/utils"
	"github.com/google/go-cmp/cmp"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type testCase struct {
	name string

	argName      string
	argNamespace string
	argSpec      *myapigroupv1alpha1.MyAppResourceSpec

	expected *appsv1.Deployment
}

func TestGetDeployment(t *testing.T) {

	for _, tc := range []testCase{
		{
			name:         "MyAppResourceSpec with default spec",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec: &myapigroupv1alpha1.MyAppResourceSpec{
				ReplicaCount: utils.Ptr[int32](3),
				Resources: &myapigroupv1alpha1.Resources{
					MemoryLimit: "160Mi",
					CPURequest:  "200m",
				},
				Image: &myapigroupv1alpha1.Image{
					Repository: "ghcr.io/stefanprodan/podinfo",
					Tag:        "latest",
				},
				UI: &myapigroupv1alpha1.UI{
					Color:   "#34577c",
					Message: "some string",
				},
				Redis: &myapigroupv1alpha1.Redis{
					Enabled: true,
				},
			},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: utils.Ptr[int32](3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "testName-podinfo",
							"app.kubernetes.io/namespace": "testNamespace",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testName-podinfo",
							Namespace: "testNamespace",
							Labels: map[string]string{
								"app.kubernetes.io/name":      "testName-podinfo",
								"app.kubernetes.io/namespace": "testNamespace",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "testName-podinfo",
									Image: "ghcr.io/stefanprodan/podinfo:latest",
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: 9898,
											Protocol:      "TCP",
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "PODINFO_UI_COLOR",
											Value: "#34577c",
										},
										{
											Name:  "PODINFO_UI_MESSAGE",
											Value: "some string",
										},
										{
											Name:  "PODINFO_CACHE_SERVER",
											Value: "redis.server.com:6379",
										},
									},
									Resources: corev1.ResourceRequirements{
										Limits: corev1.ResourceList{
											corev1.ResourceMemory: resource.MustParse("160Mi"),
										},
										Requests: corev1.ResourceList{
											corev1.ResourceCPU: resource.MustParse("200m"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:         "MyAppResourceSpec with empty Resources spec",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec: &myapigroupv1alpha1.MyAppResourceSpec{
				ReplicaCount: utils.Ptr[int32](3),
				Image: &myapigroupv1alpha1.Image{
					Repository: "ghcr.io/stefanprodan/podinfo",
					Tag:        "latest",
				},
				UI: &myapigroupv1alpha1.UI{
					Color:   "#34577c",
					Message: "some string",
				},
				Redis: &myapigroupv1alpha1.Redis{
					Enabled: true,
				},
			},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: utils.Ptr[int32](3),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app.kubernetes.io/name":      "testName-podinfo",
							"app.kubernetes.io/namespace": "testNamespace",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "testName-podinfo",
							Namespace: "testNamespace",
							Labels: map[string]string{
								"app.kubernetes.io/name":      "testName-podinfo",
								"app.kubernetes.io/namespace": "testNamespace",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "testName-podinfo",
									Image: "ghcr.io/stefanprodan/podinfo:latest",
									Ports: []corev1.ContainerPort{
										{
											Name:          "http",
											ContainerPort: 9898,
											Protocol:      "TCP",
										},
									},
									Env: []corev1.EnvVar{
										{
											Name:  "PODINFO_UI_COLOR",
											Value: "#34577c",
										},
										{
											Name:  "PODINFO_UI_MESSAGE",
											Value: "some string",
										},
										{
											Name:  "PODINFO_CACHE_SERVER",
											Value: "redis.server.com:6379",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:         "empty MyAppResourceSpec",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec:      &myapigroupv1alpha1.MyAppResourceSpec{},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
			},
		},
		{
			name:         "empty MyAppResourceSpec",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec:      &myapigroupv1alpha1.MyAppResourceSpec{ReplicaCount: utils.Ptr[int32](3)},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
			},
		},
		{
			name:         "MyAppResourceSpec with empty image",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec:      &myapigroupv1alpha1.MyAppResourceSpec{},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
			},
		},
		{
			name:         "MyAppResourceSpec with empty image repo",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec: &myapigroupv1alpha1.MyAppResourceSpec{
				Image: &myapigroupv1alpha1.Image{
					Tag: "testTag",
				},
			},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
			},
		},
		{
			name:         "MyAppResourceSpec with empty image tag",
			argNamespace: "testNamespace",
			argName:      "testName",
			argSpec: &myapigroupv1alpha1.MyAppResourceSpec{
				Image: &myapigroupv1alpha1.Image{
					Repository: "foo.com",
					Tag:        "",
				},
			},
			expected: &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testName-podinfo",
					Namespace: "testNamespace",
					Labels: map[string]string{
						"app.kubernetes.io/name":      "testName-podinfo",
						"app.kubernetes.io/namespace": "testNamespace",
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			podinfoDeployment := GetDeployment(tc.argName, tc.argNamespace, "redis.server.com:6379", tc.argSpec)

			if diff := cmp.Diff(tc.expected, podinfoDeployment); diff != "" {
				t.Errorf("GetDeployment: mismatch (-want +got):\n%s", diff)
			}

		})
	}
}
