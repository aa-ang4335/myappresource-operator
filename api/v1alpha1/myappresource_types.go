/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MyAppResourceSpec defines the desired state of MyAppResource.
type MyAppResourceSpec struct {

	// ReplicaCount specifies the number of frontend replicas.
	ReplicaCount *int32 `json:"replicaCount,omitempty"`

	// Resources specifies system resources for the frontend pods.
	Resources *Resources `json:"resources,omitempty"`

	// Image specifies the image information for the frontend pods.
	Image *Image `json:"image,omitempty"`

	// UI specifies the UI configuration for the frontend pods.
	UI *UI `json:"ui,omitempty"`

	// Redis specifies the Redis configuration for the frontend pods.
	Redis *Redis `json:"redis,omitempty"`
}

// Resources defines the resource requirements for the frontend pods.
type Resources struct {
	// MemoryLimit specifies the maximum memory limit for the frontend pods.
	MemoryLimit string `json:"memoryLimit,omitempty"`

	// CPURequest specifies the CPU request for the frontend pods.
	CPURequest string `json:"cpuRequest,omitempty"`
}

// Image specifies the details of the container image.
type Image struct {
	// Repository specifies the repository of the container image.
	Repository string `json:"repository,omitempty"`
	// Tag specifies the tag of the container image.
	Tag string `json:"tag,omitempty"`
}

// UI specifies the configuration for the user interface.
type UI struct {
	// Color specifies the color scheme for the user interface.
	Color string `json:"color,omitempty"`
	// Message specifies a message for the user interface.
	Message string `json:"message,omitempty"`
}

// Redis specifies the configuration for Redis.
type Redis struct {
	// Enabled indicates whether Redis is enabled or not.
	Enabled bool `json:"enabled,omitempty"`
}

// MyAppResourceStatus defines the observed state of MyAppResource
type MyAppResourceStatus struct {
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// MyAppResource is the Schema for the myappresources API
type MyAppResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyAppResourceSpec   `json:"spec,omitempty"`
	Status MyAppResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyAppResourceList contains a list of MyAppResource
type MyAppResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyAppResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyAppResource{}, &MyAppResourceList{})
}
