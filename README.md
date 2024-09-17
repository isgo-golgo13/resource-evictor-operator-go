# Resource Evictor Operator in Go
Go 1.23 and OperatorSDK for Kubernetes to provide a Resource Evictor Operator for Deployments and StatefulSets not providing resources.requests and limits



## Run the Operator

For Kubernetes `Deployments` without `resources.requests and resources.limits' set.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: no-resource-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      app: no-resource-app
  template:
    metadata:
      labels:
        app: no-resource-app
    spec:
      containers:
      - name: app
        image: <app-image>:latest
        resources: {}  # Missing limits
```
To apply and trigger the resource evictor.

```shell
kubectl apply -f no-resource-deployment.yaml
```


For Kubernetes `StatefulSets` without `resources.requests and resources.limits' set.

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: no-resource-statefulset
spec:
  serviceName: "no-resource-service"
  replicas: 1
  selector:
    matchLabels:
      app: no-resource-app
  template:
    metadata:
      labels:
        app: no-resource-app
    spec:
      containers:
        - name: app
          image: <app-image>:latest
          resources: {}  # Missing resource limits
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```

To apply and trigger the resource evictor.

```shell
kubectl apply -f no-resource-statefulset.yaml
```
