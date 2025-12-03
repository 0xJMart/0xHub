# 0xHub Helm Chart

This Helm chart deploys the 0xHub project hub with all its components:
- Backend API server
- Frontend React application
- Kubernetes Operator
- Project CRD

## Prerequisites

- Kubernetes 1.24+
- Helm 3.0+
- kubectl configured to access your cluster

## Installation

### Install from local chart

```bash
helm install 0xhub ./helm/0xhub
```

### Install with custom values

```bash
helm install 0xhub ./helm/0xhub -f my-values.yaml
```

## Configuration

The following table lists the configurable parameters and their default values:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.image.repository` | Backend image repository | `0xhub/backend` |
| `backend.image.tag` | Backend image tag | `latest` |
| `backend.replicaCount` | Number of backend replicas | `1` |
| `frontend.image.repository` | Frontend image repository | `0xhub/frontend` |
| `frontend.image.tag` | Frontend image tag | `latest` |
| `frontend.replicaCount` | Number of frontend replicas | `1` |
| `operator.image.repository` | Operator image repository | `0xhub/operator` |
| `operator.image.tag` | Operator image tag | `latest` |
| `operator.backendURL` | Backend URL for operator | `http://0xhub-backend:8080` |
| `crd.install` | Whether to install Project CRD | `true` |
| `rbac.create` | Whether to create RBAC resources | `true` |
| `namespace.create` | Whether to create namespace | `true` |
| `namespace.name` | Namespace name | `0xhub` |

## Usage

After installation, you can create Project resources:

```yaml
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: my-project
spec:
  name: My Project
  description: A sample project
  url: https://example.com
  category: web
  status: active
```

## Uninstallation

```bash
helm uninstall 0xhub
```

## Development

This is a basic Helm chart structure. For production use, you should:
- Add proper resource limits and requests
- Configure ingress for frontend
- Set up persistent storage if needed
- Add monitoring and logging
- Configure secrets management
- Add health checks and probes

