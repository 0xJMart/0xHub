#!/bin/bash
# Test script to install and verify the Project CRD using kind

set -e

export PATH="$HOME/.local/bin:$PATH"

echo "=== Testing Project CRD with kind ==="
echo ""

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "Installing kind..."
    mkdir -p ~/.local/bin
    curl -Lo ~/.local/bin/kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
    chmod +x ~/.local/bin/kind
    export PATH="$HOME/.local/bin:$PATH"
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "Installing kubectl..."
    mkdir -p ~/.local/bin
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    mv kubectl ~/.local/bin/kubectl
    export PATH="$HOME/.local/bin:$PATH"
fi

# Check if cluster exists
if kind get clusters | grep -q "^0xhub-test$"; then
    echo "Using existing kind cluster: 0xhub-test"
    kubectl config use-context kind-0xhub-test
else
    echo "Creating kind cluster: 0xhub-test"
    kind create cluster --name 0xhub-test --wait 60s
fi

echo ""
echo "=== Installing Project CRD ==="
kubectl apply -f project.yaml

echo "Waiting for CRD to be established..."
kubectl wait --for condition=established --timeout=60s crd/projects.hub.0xhub.io

echo ""
echo "=== Verifying CRD Installation ==="
kubectl get crd projects.hub.0xhub.io

echo ""
echo "=== Testing Project Resource Creation ==="
echo "Creating example project..."
kubectl apply -f examples/example-project.yaml

echo "Creating minimal project..."
kubectl apply -f examples/example-project-minimal.yaml

echo ""
echo "=== Listing Projects ==="
kubectl get projects

echo ""
echo "=== Testing Short Names ==="
kubectl get proj

echo ""
echo "=== Testing Validation (should fail) ==="
echo "Testing invalid project (missing required fields)..."
if kubectl apply -f - <<EOF 2>&1 | grep -q "Required value"; then
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: test-invalid
spec:
  name: Test Invalid
EOF
    echo "✓ Validation correctly rejected missing required fields"
else
    echo "✗ Validation test failed"
    exit 1
fi

echo ""
echo "Testing invalid status enum..."
if kubectl apply -f - <<EOF 2>&1 | grep -q "Unsupported value"; then
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: test-invalid-status
spec:
  name: Test Invalid Status
  description: Test
  url: https://example.com
  status: invalid-value
EOF
    echo "✓ Validation correctly rejected invalid status enum"
else
    echo "✗ Enum validation test failed"
    exit 1
fi

echo ""
echo "=== Project Details ==="
kubectl get project example-web-app -o yaml | grep -A 10 "^spec:"

echo ""
echo "=== All Tests Passed! ==="
echo ""
echo "To clean up the cluster, run:"
echo "  kind delete cluster --name 0xhub-test"

