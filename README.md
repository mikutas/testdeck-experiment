# testdeck-experiment

## Prepare local cluster and registry

[kind](kind.sigs.k8s.io/) required

```
make cluster
```

## Build

```
make docker-build-local
```

## Push

```
make docker-push-local
```

## Deploy

```
kubectl apply -f grpc-prodinfo-server.yaml
```
