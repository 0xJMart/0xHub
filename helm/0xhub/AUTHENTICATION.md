# Setting Up Authentication for Private GHCR Images

Since your GitHub Container Registry packages are private, you need to set up authentication in your Kubernetes cluster to pull the images.

## Option 1: Create Kubernetes Secret (Recommended)

### Step 1: Create a GitHub Personal Access Token (PAT)

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate a new token with the `read:packages` scope
3. Copy the token (you won't see it again!)

### Step 2: Create the Kubernetes Secret

```bash
# Create the secret in the target namespace
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_GITHUB_TOKEN> \
  --namespace=0xhub

# Or if namespace doesn't exist yet, create it first:
kubectl create namespace 0xhub
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_GITHUB_TOKEN> \
  --namespace=0xhub
```

### Step 3: Install/Upgrade Helm Chart with the Secret

```bash
helm install 0xhub oci://ghcr.io/0xjmart/0xhub/0xhub \
  --namespace 0xhub \
  --create-namespace \
  --set imagePullSecrets[0].name=ghcr-secret

# Or upgrade existing installation:
helm upgrade 0xhub oci://ghcr.io/0xjmart/0xhub/0xhub \
  --namespace 0xhub \
  --set imagePullSecrets[0].name=ghcr-secret
```

### Step 4: Verify the Secret is Used

```bash
# Check that pods are using the secret
kubectl get pods -n 0xhub -o jsonpath='{.items[*].spec.imagePullSecrets[*].name}'

# Check pod status to ensure images are being pulled
kubectl get pods -n 0xhub
```

## Option 2: Use Service Account with Image Pull Secret

If you want to use a service account (recommended for production):

```bash
# Create the secret
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_GITHUB_TOKEN> \
  --namespace=0xhub

# Patch the default service account to use the secret
kubectl patch serviceaccount default -n 0xhub \
  -p '{"imagePullSecrets": [{"name": "ghcr-secret"}]}'
```

Then install the Helm chart normally - it will automatically use the service account's imagePullSecrets.

## Option 3: Using Flux (GitOps)

If you're using Flux, you can create a `Secret` resource:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: ghcr-secret
  namespace: 0xhub
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: <base64-encoded-docker-config>
```

To generate the base64-encoded docker config:

```bash
# Create docker config JSON
echo -n '{"auths":{"ghcr.io":{"username":"<YOUR_GITHUB_USERNAME>","password":"<YOUR_GITHUB_TOKEN>","auth":"'$(echo -n '<YOUR_GITHUB_USERNAME>:<YOUR_GITHUB_TOKEN>' | base64)'"}}}' | base64 -w 0
```

Or use `kubectl create secret` and then extract it:

```bash
kubectl create secret docker-registry ghcr-secret \
  --docker-server=ghcr.io \
  --docker-username=<YOUR_GITHUB_USERNAME> \
  --docker-password=<YOUR_GITHUB_TOKEN> \
  --namespace=0xhub \
  --dry-run=client -o yaml | \
  kubectl apply -f -
```

Then in your HelmRelease, reference the secret:

```yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: 0xhub
  namespace: 0xhub
spec:
  chart:
    spec:
      chart: 0xhub
      sourceRef:
        kind: HelmRepository
        name: 0xhub-repo
  values:
    imagePullSecrets:
      - name: ghcr-secret
```

## Troubleshooting

### Check if secret exists:
```bash
kubectl get secret ghcr-secret -n 0xhub
```

### View secret details:
```bash
kubectl describe secret ghcr-secret -n 0xhub
```

### Check pod events for pull errors:
```bash
kubectl describe pod <pod-name> -n 0xhub
```

### Test image pull manually:
```bash
# Create a test pod to verify authentication
kubectl run test-pull --image=ghcr.io/0xjmart/0xhub/backend:main \
  --image-pull-policy=Always \
  --overrides='{"spec":{"imagePullSecrets":[{"name":"ghcr-secret"}]}}' \
  --rm -it --restart=Never -- /bin/sh
```

## Security Best Practices

1. **Use Fine-Grained PATs**: Create tokens with minimal required scopes (`read:packages` only)
2. **Rotate Tokens Regularly**: Update secrets periodically
3. **Use Namespace-Scoped Secrets**: Don't use cluster-wide secrets unless necessary
4. **Consider Using External Secrets Operator**: For better secret management in production
5. **Use Service Accounts**: Instead of adding secrets to each deployment

