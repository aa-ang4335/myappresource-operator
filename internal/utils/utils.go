package utils

// Ptr returns a pointer to the placeholder object
func Ptr[C any](in C) *C {
	return &in
}

// MergeLabels merges all the label maps into a single new label map.
func MergeLabels(allLabels ...map[string]string) map[string]string {
	out := map[string]string{}

	for _, labels := range allLabels {
		for k, v := range labels {
			out[k] = v
		}
	}
	return out
}

// GenerateDefaultLabels generates a set of commong labels across k8s resources managed by the operator.
func GenerateDefaultLabels(name, namespace string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      name,
		"app.kubernetes.io/namespace": namespace,
	}
}
