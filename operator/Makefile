.PHONY: build run test docker-build deploy clean

# Image URL to use all building/pushing image targets
IMG ?= 0xhub-operator:latest
# Backend URL for local development
BACKEND_URL ?= http://localhost:8080

# Build the operator binary
build:
	go build -o bin/manager cmd/manager/main.go

# Run the operator locally
run: build
	./bin/manager --backend-url=$(BACKEND_URL)

# Run tests
test:
	go test ./... -v

# Build the docker image
docker-build:
	docker build -t ${IMG} .

# Deploy to kind cluster
deploy-kind:
	kubectl apply -f config/rbac/
	kubectl apply -f config/manager/deployment.yaml

# Clean build artifacts
clean:
	rm -rf bin/

# Install CRD
install-crd:
	kubectl apply -f ../crd/project.yaml

# Uninstall CRD
uninstall-crd:
	kubectl delete -f ../crd/project.yaml

# Load image into kind
kind-load:
	kind load docker-image ${IMG} --name 0xhub-test

