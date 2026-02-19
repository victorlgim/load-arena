APP := load-arena
IMG := load-arena:dev
NS  := load-arena
KIND_NAME := load-arena

.PHONY: all tidy test build docker kind-load k8s-apply k8s-delete pf logs top hpa

all: docker kind-load k8s-apply

tidy:
	go mod tidy

test:
	go test ./...

build:
	CGO_ENABLED=0 go build -o bin/$(APP) ./cmd/api

docker:
	docker build -t $(IMG) -f deploy/docker/Dockerfile .

kind-load:
	kind load docker-image $(IMG) --name $(KIND_NAME)

k8s-apply:
	kubectl apply -f deploy/k8s/namespace.yaml
	kubectl apply -f deploy/k8s/configmap.yaml
	kubectl apply -f deploy/k8s/deployment.yaml
	kubectl apply -f deploy/k8s/service.yaml
	@echo "Optional:"
	@echo "  kubectl apply -f deploy/k8s/servicemonitor.yaml"
	@echo "  kubectl apply -f deploy/k8s/hpa.yaml"
	@echo "  kubectl apply -f deploy/k8s/ingress.yaml"
	@echo "  kubectl apply -f deploy/k8s/pdb.yaml"

k8s-delete:
	kubectl delete -f deploy/k8s/ -n $(NS) --ignore-not-found || true
	kubectl delete ns $(NS) --ignore-not-found || true

pf:
	kubectl port-forward -n $(NS) svc/$(APP) 8080:80

logs:
	kubectl logs -n $(NS) deploy/$(APP) -f

top:
	kubectl top pods -n $(NS)

hpa:
	kubectl apply -f deploy/k8s/hpa.yaml
	kubectl get hpa -n $(NS) -w
