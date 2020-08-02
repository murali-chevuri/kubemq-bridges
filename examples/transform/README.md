# KubeMQ Bridges - Transform Example

In this example we demonstrate how to transform and aggregating of events in one cluster and replicate them as qeuue messages to other clusters.

![transform-example](../../.github/assets/transform-example.jpeg)

## Run

Run the following deployment

```bash
kubectl apply -f ./deploy.yaml
```
Where deploy.yaml:

```yaml
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-a
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30501
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-b
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30502
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-c
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30503
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-d
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30504
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-e
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30505
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-f
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30506
---
apiVersion: core.k8s.kubemq.io/v1alpha1
kind: KubemqCluster
metadata:
  name: kubemq-cluster-g
  namespace: kubemq
spec:
  replicas: 3
  grpc:
    expose: NodePort
    nodePort: 30507
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubemq-bridges-deployment
  namespace: kubemq
  labels:
    app: kubemq-bridges
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kubemq-bridges
  template:
    metadata:
      labels:
        app: kubemq-bridges
    spec:
      containers:
        - name: kubemq-bridges
          image: kubemq/kubemq-bridges:latest
          ports:
            - containerPort: 8080
          volumeMounts:
            - mountPath: /kubemq-bridges/config.yaml
              name: config-file
              subPath: config.yaml
      volumes:
        - name: config-file
          configMap:
            name: kubemq-bridges-config
            items:
              - key: config.yaml
                path: config.yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: kubemq-bridges-config
  namespace: kubemq
data:
  config.yaml: |-
    apiPort: 8081
    bindings:
      - name: cluster-a-d-events-aggregate
        properties:
          log_level: "debug"
        sources:
          kind: source.events
          name: cluster-sources
          connections:
            - address: "kubemq-cluster-a-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-a-events-source"
              auth_token: ""
              channel: "events.a"
              group:   ""
              concurrency: "1"
              auto_reconnect: "true"
              reconnect_interval_seconds: "1"
              max_reconnects: "0"
            - address: "kubemq-cluster-b-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-b-events-source"
              auth_token: ""
              channel: "events.b"
              group:   ""
              concurrency: "1"
              auto_reconnect: "true"
              reconnect_interval_seconds: "1"
              max_reconnects: "0"
            - address: "kubemq-cluster-c-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-c-d-events-source"
              auth_token: ""
              channel: "events.c"
              group:   ""
              concurrency: "1"
              auto_reconnect: "true"
              reconnect_interval_seconds: "1"
              max_reconnects: "0"
        targets:
          kind: target.events
          name: cluster-targets
          connections:
            - address: "kubemq-cluster-d-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-a-d-events-target"
              auth_token: ""
              channels: "events.aggregate"
      - name: cluster-transform-queue
        properties:
          log_level: "debug"
        sources:
          kind: source.queue
          name: cluster-sources
          connections:
            - address: "kubemq-cluster-d-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-d-queue-source"
              auth_token: ""
              channel: "queue.e"
        targets:
          kind: target.queue
          name: cluster-targets
          connections:
            - address: "kubemq-cluster-d-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-g-queue-target"
              auth_token: ""
              channels: "queue"
            - address: "kubemq-cluster-g-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-f-queue-target"
              auth_token: ""
              channels: "queue"
            - address: "kubemq-cluster-g-grpc.kubemq.svc.cluster.local:50000"
              client_id: "cluster-g-queue-target"
              auth_token: ""
              channels: "queue"

```
