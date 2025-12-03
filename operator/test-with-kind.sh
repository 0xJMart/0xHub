#!/bin/bash
# Test script for the operator with kind

set -e

export PATH="$HOME/.local/bin:$PATH"

CLUSTER_NAME="0xhub-test"
OPERATOR_IMAGE="0xhub-operator:latest"
# For kind, we need to use host.docker.internal or the host IP
# On Linux, we can use the host network or a service
# For now, we'll use host.docker.internal which works on Docker Desktop
# On Linux, you may need to use the host IP: $(docker network inspect kind | jq -r '.[0].IPAM.Config[0].Gateway')
BACKEND_URL="${BACKEND_URL:-http://host.docker.internal:8080}"

echo "=== Testing 0xHub Operator with kind ==="
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
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
    echo "Using existing kind cluster: ${CLUSTER_NAME}"
    kubectl config use-context kind-${CLUSTER_NAME}
else
    echo "Creating kind cluster: ${CLUSTER_NAME}"
    kind create cluster --name ${CLUSTER_NAME} --wait 60s
fi

echo ""
echo "=== Installing Project CRD ==="
kubectl apply -f ../crd/project.yaml
kubectl wait --for condition=established --timeout=60s crd/projects.hub.0xhub.io

echo ""
echo "=== Building Operator ==="
cd "$(dirname "$0")"
make build

echo ""
echo "=== Building Docker Image ==="
make docker-build

echo ""
echo "=== Loading Image into kind ==="
kind load docker-image ${OPERATOR_IMAGE} --name ${CLUSTER_NAME}

echo ""
echo "=== Setting up RBAC ==="
kubectl create namespace system --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f config/rbac/

echo ""
echo "=== Deploying Operator ==="
# Create deployment with backend URL
cat > /tmp/operator-deployment.yaml <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
  labels:
    app: 0xhub-operator
spec:
  selector:
    matchLabels:
      app: 0xhub-operator
  replicas: 1
  template:
    metadata:
      labels:
        app: 0xhub-operator
    spec:
      serviceAccountName: controller-manager
      containers:
      - command:
        - /manager
        args:
        - --backend-url=${BACKEND_URL}
        image: ${OPERATOR_IMAGE}
        name: manager
        resources:
          limits:
            cpu: 500m
            memory: 512Mi
          requests:
            cpu: 100m
            memory: 128Mi
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
EOF
kubectl apply -f /tmp/operator-deployment.yaml
rm /tmp/operator-deployment.yaml

echo ""
echo "=== Waiting for Operator to be Ready ==="
kubectl wait --for=condition=available --timeout=120s deployment/controller-manager -n system

echo ""
echo "=== Checking Operator Logs ==="
kubectl logs -n system deployment/controller-manager --tail=20

echo ""
echo "=== Creating Test Project ==="
kubectl apply -f - <<EOF
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: test-project
  namespace: default
spec:
  name: Test Project
  description: A test project created by the operator test script
  url: https://test.example.com
  category: testing
  status: active
EOF

echo ""
echo "=== Waiting for Project to be Synced ==="
sleep 5

echo ""
echo "=== Checking Project Status ==="
kubectl get project test-project -o yaml

echo ""
echo "=== Checking Operator Logs for Sync ==="
kubectl logs -n system deployment/controller-manager --tail=10

echo ""
echo "=== Test Complete ==="
echo ""
echo "To check the operator logs:"
echo "  kubectl logs -n system deployment/controller-manager -f"
echo ""
echo "To check project status:"
echo "  kubectl get projects"
echo "  kubectl get project test-project -o yaml"
echo ""
echo "To clean up:"
echo "  kind delete cluster --name ${CLUSTER_NAME}"

