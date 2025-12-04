# 0xHub Helm Chart

This Helm chart deploys the 0xHub project hub with all its components:
- Backend API server (Go/Gin)
- Frontend React application (served via nginx)
- Kubernetes Operator (controller-runtime)
- Project CRD (Custom Resource Definition)

## Prerequisites

- Kubernetes 1.24+
- Helm 3.0+
- kubectl configured to access your cluster
- Docker images built and available in your registry (or use local images)

## Installation

### Quick Start

```bash
# Install with default values
helm install 0xhub ./helm/0xhub

# Install in a specific namespace
helm install 0xhub ./helm/0xhub --namespace 0xhub --create-namespace

# Install with custom values file
helm install 0xhub ./helm/0xhub -f my-values.yaml

# Install with inline value overrides
helm install 0xhub ./helm/0xhub \
  --set backend.image.tag=v1.0.0 \
  --set frontend.image.tag=v1.0.0 \
  --set operator.image.tag=v1.0.0
```

### Building Docker Images

Before installing, you need to build and push the Docker images:

```bash
# Build backend
cd backend
docker build -t 0xhub/backend:latest .

# Build frontend
cd frontend
docker build -t 0xhub/frontend:latest .

# Build operator
cd operator
docker build -t 0xhub/operator:latest .

# Push to your registry (example)
docker tag 0xhub/backend:latest your-registry/0xhub/backend:latest
docker push your-registry/0xhub/backend:latest
```

## Configuration

The following table lists the configurable parameters and their default values:

### Global Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `nameOverride` | String to partially override fullname | `""` |
| `fullnameOverride` | String to fully override fullname | `""` |
| `imagePullSecrets` | Image pull secrets | `[]` |

### Backend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `backend.image.repository` | Backend image repository | `0xhub/backend` |
| `backend.image.tag` | Backend image tag | `latest` |
| `backend.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `backend.replicaCount` | Number of backend replicas | `1` |
| `backend.service.type` | Service type | `ClusterIP` |
| `backend.service.port` | Service port | `8080` |
| `backend.resources.limits.cpu` | CPU limit | `500m` |
| `backend.resources.limits.memory` | Memory limit | `512Mi` |
| `backend.resources.requests.cpu` | CPU request | `100m` |
| `backend.resources.requests.memory` | Memory request | `128Mi` |
| `backend.env` | Additional environment variables | `[]` |
| `backend.podSecurityContext` | Pod security context | `{}` |
| `backend.securityContext` | Container security context | `{}` |
| `backend.podAnnotations` | Pod annotations | `{}` |
| `backend.podLabels` | Pod labels | `{}` |
| `backend.nodeSelector` | Node selector | `{}` |
| `backend.tolerations` | Tolerations | `[]` |
| `backend.affinity` | Affinity rules | `{}` |

### Frontend Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `frontend.image.repository` | Frontend image repository | `0xhub/frontend` |
| `frontend.image.tag` | Frontend image tag | `latest` |
| `frontend.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `frontend.replicaCount` | Number of frontend replicas | `1` |
| `frontend.service.type` | Service type | `ClusterIP` |
| `frontend.service.port` | Service port | `80` |
| `frontend.resources.limits.cpu` | CPU limit | `200m` |
| `frontend.resources.limits.memory` | Memory limit | `256Mi` |
| `frontend.resources.requests.cpu` | CPU request | `50m` |
| `frontend.resources.requests.memory` | Memory request | `64Mi` |
| `frontend.env` | Additional environment variables | `[]` |
| `frontend.podSecurityContext` | Pod security context | `{}` |
| `frontend.securityContext` | Container security context | `{}` |
| `frontend.podAnnotations` | Pod annotations | `{}` |
| `frontend.podLabels` | Pod labels | `{}` |
| `frontend.nodeSelector` | Node selector | `{}` |
| `frontend.tolerations` | Tolerations | `[]` |
| `frontend.affinity` | Affinity rules | `{}` |

### Operator Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `operator.image.repository` | Operator image repository | `0xhub/operator` |
| `operator.image.tag` | Operator image tag | `latest` |
| `operator.image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `operator.replicaCount` | Number of operator replicas | `1` |
| `operator.backendURL` | Backend URL for operator (auto-generated if empty) | `""` |
| `operator.leaderElection` | Enable leader election | `true` |
| `operator.args` | Additional command arguments | `[]` |
| `operator.env` | Additional environment variables | `[]` |
| `operator.resources.limits.cpu` | CPU limit | `500m` |
| `operator.resources.limits.memory` | Memory limit | `512Mi` |
| `operator.resources.requests.cpu` | CPU request | `100m` |
| `operator.resources.requests.memory` | Memory request | `128Mi` |
| `operator.podSecurityContext` | Pod security context | `{}` |
| `operator.securityContext` | Container security context | `{}` |
| `operator.podAnnotations` | Pod annotations | `{}` |
| `operator.podLabels` | Pod labels | `{}` |
| `operator.nodeSelector` | Node selector | `{}` |
| `operator.tolerations` | Tolerations | `[]` |
| `operator.affinity` | Affinity rules | `{}` |

### Other Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `crd.install` | Whether to install Project CRD | `true` |
| `rbac.create` | Whether to create RBAC resources | `true` |
| `serviceAccount.create` | Whether to create service account | `true` |
| `serviceAccount.name` | Service account name (auto-generated if empty) | `""` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `namespace.create` | Whether to create namespace | `true` |
| `namespace.name` | Namespace name | `0xhub` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Ingress hosts configuration | See values.yaml |
| `ingress.tls` | Ingress TLS configuration | `[]` |

## Usage Examples

### Basic Installation

```bash
helm install 0xhub ./helm/0xhub
```

### Installation with Custom Image Registry

```yaml
# custom-values.yaml
backend:
  image:
    repository: registry.example.com/0xhub/backend
    tag: v1.0.0

frontend:
  image:
    repository: registry.example.com/0xhub/frontend
    tag: v1.0.0

operator:
  image:
    repository: registry.example.com/0xhub/operator
    tag: v1.0.0

imagePullSecrets:
  - name: regcred
```

```bash
helm install 0xhub ./helm/0xhub -f custom-values.yaml
```

### Installation with Ingress

```yaml
# ingress-values.yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: 0xhub.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: 0xhub-tls
      hosts:
        - 0xhub.example.com
```

```bash
helm install 0xhub ./helm/0xhub -f ingress-values.yaml
```

### Installation with Resource Limits

```yaml
# resources-values.yaml
backend:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 200m
      memory: 256Mi

frontend:
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 128Mi

operator:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 200m
      memory: 256Mi
```

### Installation with Node Affinity

```yaml
# affinity-values.yaml
backend:
  nodeSelector:
    kubernetes.io/os: linux
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: node-type
            operator: In
            values:
            - compute
```

### Creating Project Resources

After installation, you can create Project CRD resources:

```yaml
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: my-project
  namespace: 0xhub
spec:
  name: My Project
  description: A sample project description
  url: https://example.com
  icon: https://example.com/icon.png
  category: web
  status: active
```

```bash
kubectl apply -f project.yaml
```

The operator will automatically sync this Project to the backend, and it will appear in the frontend.

## Upgrading

```bash
# Upgrade with new values
helm upgrade 0xhub ./helm/0xhub -f my-values.yaml

# Upgrade with new image tags
helm upgrade 0xhub ./helm/0xhub \
  --set backend.image.tag=v1.1.0 \
  --set frontend.image.tag=v1.1.0 \
  --set operator.image.tag=v1.1.0
```

## Uninstallation

```bash
# Uninstall the release
helm uninstall 0xhub

# Uninstall and delete namespace
helm uninstall 0xhub
kubectl delete namespace 0xhub
```

**Note:** The CRD will remain after uninstallation. To remove it:

```bash
kubectl delete crd projects.hub.0xhub.io
```

## Troubleshooting

### Pods Not Starting

Check pod status and logs:

```bash
# Check pod status
kubectl get pods -n 0xhub

# Check backend logs
kubectl logs -n 0xhub -l app.kubernetes.io/component=backend

# Check frontend logs
kubectl logs -n 0xhub -l app.kubernetes.io/component=frontend

# Check operator logs
kubectl logs -n 0xhub -l app.kubernetes.io/component=operator
```

### Operator Not Syncing Projects

1. Verify the operator can reach the backend:

```bash
# Check operator logs for connection errors
kubectl logs -n 0xhub -l app.kubernetes.io/component=operator

# Test backend connectivity from operator pod
kubectl exec -n 0xhub -it deployment/0xhub-operator -- wget -O- http://0xhub-backend:8080/api/health
```

2. Verify the backend URL is correct:

```bash
# Check operator environment variables
kubectl exec -n 0xhub -it deployment/0xhub-operator -- env | grep BACKEND_URL
```

3. Check Project CRD status:

```bash
# List projects
kubectl get projects -n 0xhub

# Check project status
kubectl describe project <project-name> -n 0xhub
```

### Frontend Not Loading

1. Check if frontend service is accessible:

```bash
# Port forward to test locally
kubectl port-forward -n 0xhub svc/0xhub-frontend 8080:80

# Access at http://localhost:8080
```

2. Verify API URL configuration:

```bash
# Check frontend environment variables
kubectl exec -n 0xhub -it deployment/0xhub-frontend -- env | grep VITE_API_URL
```

3. Check nginx configuration:

```bash
# Check nginx config in frontend pod
kubectl exec -n 0xhub -it deployment/0xhub-frontend -- cat /etc/nginx/conf.d/default.conf
```

### Image Pull Errors

If you see `ImagePullBackOff` errors:

1. Verify image names and tags in values.yaml
2. Check image pull secrets are configured correctly
3. Verify registry credentials:

```bash
# Create image pull secret
kubectl create secret docker-registry regcred \
  --docker-server=<registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  --docker-email=<email> \
  -n 0xhub
```

### CRD Installation Issues

If the CRD fails to install:

```bash
# Check CRD status
kubectl get crd projects.hub.0xhub.io

# Manually install CRD
kubectl apply -f crd/project.yaml

# Check CRD validation
kubectl describe crd projects.hub.0xhub.io
```

## Testing

### Dry Run Installation

```bash
# Test template rendering
helm template 0xhub ./helm/0xhub

# Test with dry-run
helm install 0xhub ./helm/0xhub --dry-run --debug
```

### Validate Configuration

```bash
# Lint the chart
helm lint ./helm/0xhub

# Validate against Kubernetes schema
helm template 0xhub ./helm/0xhub | kubeval
```

## Architecture

The chart deploys the following components:

1. **Backend Deployment & Service**: Go API server on port 8080
2. **Frontend Deployment & Service**: React app served via nginx on port 80
3. **Operator Deployment**: Kubernetes operator watching Project CRDs
4. **CRD**: Project CustomResourceDefinition
5. **RBAC**: ServiceAccount, ClusterRole, ClusterRoleBinding for operator
6. **Ingress** (optional): External access to frontend

## Security Considerations

- Use specific image tags instead of `latest` in production
- Configure resource limits and requests appropriately
- Enable security contexts for pods and containers
- Use image pull secrets for private registries
- Configure network policies to restrict traffic
- Use TLS for ingress in production
- Regularly update images for security patches

## Development

For local development and testing:

```bash
# Install in development mode with local images
helm install 0xhub ./helm/0xhub \
  --set backend.image.pullPolicy=Never \
  --set frontend.image.pullPolicy=Never \
  --set operator.image.pullPolicy=Never \
  --set namespace.create=true \
  --set namespace.name=0xhub-dev
```

## Support

For issues and questions:
- Check the troubleshooting section above
- Review pod logs and events
- Verify all configuration values
- Check Kubernetes cluster resources
