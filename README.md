# MyAppResource-operator

The MyAppResource operator handles creating, configuring, and managing podinfo on Kubernetes.
For more details about podinfo, please refer to the project link:
<https://github.com/stefanprodan/podinfo>

## Requirements

- Kubernetes version: 1.21 or higher
- kind version 0.22 or higher
- Go version 1.21
- kustomize
- Makefile

The MyAppResource operator is being tested against kind Kubernetes 1.29. All dependencies have been vendored to prevent any dependency downloads during building.

## Deploying the operator

Before deploying the operator, ensure that the cluster is created via kind, and the correct Kubernetes context is set. Follow these steps to verify and set the appropriate context:

1. create kind cluster.

   ```sh
   kind create cluster
   ```

1. verify the current Kubernetes context.

   ```sh
   kubectl config current-context
   ```

2. Set the correct context if the value doesn't match the desired cluster managed by kind.

   ```sh
   kubectl config use-context <context-name>
   ```

3. Clone the operator repository from GitHub:

   ```sh
   git clone <tbd>
   ```

4. Build and deploy the provider using the steps below:

   ```sh
   img="example.com/int03012024:v0.0.1"

   # For Darwin/amd64 users, ensure you are using a darwin docker builder profile. 
   make docker-build "IMG=${img}"

   kind load docker-image "${img}"
   make install
   make deploy "IMG=${img}"
   ```

## Validation

* Ensure that existing tests pass successfully:
```
# Run unit tests
make test

# Run end-to-end tests
make test-e2e
```

* Apply the example CRD without encountering any errors:
```
kubectl apply -f example/whatever.yaml
kubectl wait myappresource.my.api.group/whatever --for=jsonpath='{.status.valid}'=true --timeout=5m
```

* Perform port forwarding to the podinfo deployment:
```
kubectl port-forward deploy/whatever-podinfo 9898:9898
```

* Verify that the endpoint returns the expected environment variables:
```
curl -s localhost:9898/env | grep -i "podinfo"
```

* Confirm that the endpoint can store data in the cache server:
```
# Store 'bar' under the key 'foo'
curl -s -X POST -d "bar" localhost:9898/cache/foo

# Retrieve the stored value for 'foo'
curl -s localhost:9898/cache/foo
```