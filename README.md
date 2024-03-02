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

## Design

The list below covers various components of the operator code that's responsible for  managing Kubernetes resources.

#### Operator Code

- api -> Definition of the myappresource CRD.
- cmd -> Contains the starting point of the application.
- internal/controller/cleaner -> Contains functions to clean up Kubernetes resources.
- internal/controller/myappresource_controller -> The main controller logic manages requests from Kubernetes, creating, updating, or deleting pieces as necessary.
- internal/service/podinfo -> The logic to generate podinfo Kubernetes resources from the values defined in the CRD.
- internal/service/redis -> The logic to generate Redis Kubernetes resources from the values defined in the CRD.
- vendor -> Vendored packages used by the application.

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
   git clone https://github.com/aa-ang4335/myappresource-operator.git
   ```

4. Build and deploy the provider using the steps below:

   ```sh
   img="example.com/myappresource-operator:v0.0.1"

   # For Darwin/amd64 users, ensure you are using a darwin docker builder profile. 
   make docker-build "IMG=${img}"

   kind load docker-image "${img}"
   make install
   make deploy "IMG=${img}"
   ```

## Validation

- Ensure that existing tests pass successfully:

```
# Run unit tests
make test

# Run end-to-end tests
make test-e2e
```

- Apply the example CRD without encountering any errors:

```
kubectl apply -f example/whatever.yaml
kubectl wait myappresource.my.api.group/whatever --for=jsonpath='{.status.valid}'=true --timeout=5m
```

- Perform port forwarding to the podinfo deployment:

```
kubectl port-forward deploy/whatever-podinfo 9898:9898
```

- Verify that the endpoint returns the expected environment variables:

```
curl -s localhost:9898/env | grep -i "podinfo"
```

- Confirm that the endpoint can store data in the cache server:

```
# Store 'bar' under the key 'foo'
curl -s -X POST -d "bar" localhost:9898/cache/foo

# Retrieve the stored value for 'foo'
curl -s localhost:9898/cache/foo
```

## Cleanup

To remove the operator from your Kubernetes cluster, you must first delete the Custom Resource Definition, followed by the oeprator resources. Deleting the CRD automatically removes all podinfo resources and Redis stateful sets. However, for disaster recovery purposes, the persistent volumes associated with Redis remain untouched. You can manually clean them using kubectl."

```sh
# Delete the CRD
kubectl delete myappresources.my.api.group --all

# Delete the operator
kubectl delete deploy -l control-plane=controller-manager --all-namespaces

# Delete the persistent volumes
kubectl delete pvc -l app.kubernetes.io/name=whatever-redis --all-namespaces
kubectl delete pv -l app.kubernetes.io/name=whatever-redis --all-namespaces
```

## Feature Wishlist

The list below covers a set of features we'd like to see before productizing the project.

### Config Management

- Ensure the naming of the CRD object reflects the service. Replace 'myappresources' with 'podinfo' or something similar.
- Publish the operator binary via Github Action.
- Use Helm to package the operator service.

### Testing

- Ensure 100% unit test coverage.
- Add end-to-end testing coverage.

### Telemetry

- Add a Prometheus metric service for the application.
- Add telemetry for the operator/reconciler/syncer/cleaner health.
- Add ServiceMonitor in the Helm chart templates.

### Core

- Split the CRD into two different components: backend/frontend.
- Ensure that the majority of deployment specs for podinfo can be set via the CRD.
- Remove support for Redis from the operator in favor of using an independently maintained, production-ready Helm chart to support the backend.
- Ensure the reconciler runs periodically to address any configuration drifts.
- Ensure the Kubernetes resource cleanup process runs asynchronously and control the frequency of cleanups.
