SHELL := /bin/bash

# ==============================================================================
# Testing running system
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
dev.setup.mac:
	brew update
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
# ==============================================================================
# Building containers

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: images-api metrics

images-api:
	docker build \
		-f conf/docker/dockerfile.images-api \
		-t images-api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics:
	docker build \
		-f conf/docker/dockerfile.metrics \
		-t metrics-amd64:1.0 \
		--build-arg BUILD_REF=1.0\
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := semi-cluster

# Upgrade to latest Kind: brew upgrade kind
# For full Kind v0.14 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.14.0
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e \
		--name $(KIND_CLUSTER) \
		--config conf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=images-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-load:
	cd conf/k8s/kind/images-pod; kustomize edit set image images-api-image=images-api-amd64:$(VERSION)
	cd conf/k8s/kind/images-pod; kustomize edit set image metrics-image=metrics-amd64:$(VERSION)
	kind load docker-image images-api-amd64:$(VERSION) --name $(KIND_CLUSTER)
	kind load docker-image metrics-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build conf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
#	kustomize build conf/k8s/kind/zipkin-pod | kubectl apply -f -
#	kubectl wait --namespace=zipkin-system --timeout=120s --for=condition=Available deployment/zipkin-pod
	kustomize build conf/k8s/kind/images-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=images --all-containers=true -f --tail=100 --namespace=images-system | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment images-pod --namespace=images-system

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-logs-metrics:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go -service=METRICS

kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100

kind-logs-zipkin:
	kubectl logs -l app=zipkin --namespace=zipkin-system --all-containers=true -f --tail=100

kind-status-sales:
	kubectl get pods -o wide --watch --namespace=sales-system

kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-status-zipkin:
	kubectl get pods -o wide --watch --namespace=zipkin-system

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=sales

kind-describe-deployment:
	kubectl describe deployment sales-pod

kind-describe-replicaset:
	kubectl get rs
	kubectl describe rs -l app=sales

kind-events:
	kubectl get ev --sort-by metadata.creationTimestamp

kind-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

kind-context-sales:
	kubectl config set-context --current --namespace=sales-system

kind-shell:
	kubectl exec -it $(shell kubectl get pods | grep sales | cut -c1-26) --container sales-api -- /bin/sh

kind-database:
 	# ./admin --db-disable-tls=1 migrate
 	# ./admin --db-disable-tls=1 seed

 # ==============================================================================
 # Administration

migrate:
	go run app/tooling/admin/main.go migrate

seed: migrate
	go run app/tooling/admin/main.go seed

 # ==============================================================================
 # Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

 # ==============================================================================
 # Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
 	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

 # ==============================================================================
 # Docker support

docker-down:
	docker rm -f $(shell docker ps -aq)

docker-clean:
	docker system prune -f

docker-kind-logs:
	docker logs -f $(KIND_CLUSTER)-control-plane
