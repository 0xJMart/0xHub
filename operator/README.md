# 0xHub Operator

The Kubernetes operator that watches for Project CRDs and syncs them to the backend API.

## Overview

The operator continuously watches for Project Custom Resource Definitions (CRDs) in the Kubernetes cluster. When a Project resource is created, updated, or deleted, the operator automatically syncs the changes to the backend API.

## Architecture

- **Controller**: Watches Project CRDs using controller-runtime
- **Reconcile Loop**: Processes create/update/delete events
- **Backend Client**: HTTP client that communicates with the backend API
- **Status Updates**: Updates Project CRD status with sync information

## Prerequisites

- Go 1.24+
- Kubernetes cluster (or kind for local testing)
- kubectl configured to access the cluster
- Backend API running (default: http://localhost:8080)

## Building

### Build the operator binary:

```bash
make build
```

### Build the Docker image:

```bash
make docker-build
```

## Running Locally

### Option 1: Run directly (for development)

1. Make sure you have a kubeconfig configured:
   ```bash
   kubectl cluster-info
   ```

2. Install the CRD:
   ```bash
   make install-crd
   ```

3. Run the operator:
   ```bash
   make run
   # Or with custom backend URL:
   BACKEND_URL=http://localhost:8080 make run
   ```

### Option 2: Run in Kubernetes (kind)

See the [Testing with kind](#testing-with-kind) section below.

## Testing with kind

The easiest way to test the operator is using kind (Kubernetes in Docker):

```bash
./test-with-kind.sh
```

This script will:
1. Install kind and kubectl if needed
2. Create a local kind cluster
3. Install the Project CRD
4. Build and deploy the operator
5. Create a test project
6. Verify the sync operation

### Manual Testing Steps

1. **Create a kind cluster:**
   ```bash
   kind create cluster --name 0xhub-test
   ```

2. **Install the CRD:**
   ```bash
   kubectl apply -f ../crd/project.yaml
   ```

3. **Build and load the operator image:**
   ```bash
   make docker-build
   kind load docker-image 0xhub-operator:latest --name 0xhub-test
   ```

4. **Deploy the operator:**
   ```bash
   kubectl create namespace system --dry-run=client -o yaml | kubectl apply -f -
   kubectl apply -f config/rbac/
   kubectl apply -f config/manager/deployment.yaml
   ```

5. **Create a test project:**
   ```bash
   kubectl apply -f ../crd/examples/example-project.yaml
   ```

6. **Check the sync status:**
   ```bash
   kubectl get project example-web-app -o yaml
   kubectl logs -n system deployment/controller-manager
   ```

## Configuration

The operator can be configured using command-line flags or environment variables:

- `--backend-url`: URL of the backend API (default: http://localhost:8080)
- `--metrics-bind-address`: Address for metrics endpoint (default: :8080)
- `--health-probe-bind-address`: Address for health probe (default: :8081)
- `--leader-elect`: Enable leader election (default: false)

## Project Sync Flow

1. **Create**: When a Project CRD is created, the operator:
   - Converts the CRD spec to backend Project format
   - Calls `POST /api/projects` to create the project
   - Updates the CRD status with sync information

2. **Update**: When a Project CRD is updated, the operator:
   - Converts the updated spec to backend Project format
   - Calls `PUT /api/projects/{id}` to update the project
   - Updates the CRD status

3. **Delete**: When a Project CRD is deleted, the operator:
   - Calls `DELETE /api/projects/{id}` to remove the project from backend
   - The CRD is then removed from Kubernetes

## Status Fields

The operator updates the following status fields on Project CRDs:

- `status.synced`: Boolean indicating if the project was successfully synced
- `status.lastSyncedAt`: Timestamp of the last successful sync
- `status.error`: Error message if sync failed

## Troubleshooting

### Operator not syncing projects

1. Check operator logs:
   ```bash
   kubectl logs -n system deployment/controller-manager
   ```

2. Verify backend is accessible:
   ```bash
   kubectl exec -n system deployment/controller-manager -- wget -qO- http://backend-service:8080/api/health
   ```

3. Check Project CRD status:
   ```bash
   kubectl get project <project-name> -o yaml
   ```

### Backend connection errors

- Ensure the backend URL is correct in the deployment
- Check network connectivity from the operator pod to the backend
- Verify the backend is running and healthy

### RBAC issues

If you see permission errors, ensure the RBAC resources are applied:
```bash
kubectl apply -f config/rbac/
```

## Development

### Project Structure

```
operator/
├── api/v1/              # CRD type definitions
├── controllers/         # Controller logic
├── internal/backend/    # Backend API client
├── cmd/manager/         # Main entry point
├── config/              # Kubernetes manifests
│   ├── rbac/           # RBAC resources
│   └── manager/        # Deployment manifests
└── Makefile            # Build commands
```

### Adding Features

1. Update CRD types in `api/v1/project_types.go`
2. Update controller logic in `controllers/project_controller.go`
3. Update backend client if needed in `internal/backend/client.go`
4. Rebuild and redeploy

## License

MIT

