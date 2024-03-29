SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

LOCAL_IMG ?= localhost:5000/grpc-productinfo-server-test:latest

docker-build-local:
	docker build . -t ${LOCAL_IMG} --target run_test

docker-push-local:
	docker push ${LOCAL_IMG}

REG_NAME="kind-registry"

cluster:
	./kind-with-registry.sh $(REG_NAME)

cluster-off:
	kind delete cluster
	docker stop $(REG_NAME)
