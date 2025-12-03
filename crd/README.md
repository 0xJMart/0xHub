# Project Custom Resource Definition (CRD)

This directory contains the Kubernetes Custom Resource Definition for Projects in the 0xHub system.

## Overview

The Project CRD allows you to define projects as Kubernetes resources, which are automatically synced to the backend API by the operator.

## Installation

To install the CRD in your Kubernetes cluster:

```bash
kubectl apply -f crd/project.yaml
```

Verify the installation:

```bash
kubectl get crd projects.hub.0xhub.io
```

## Testing with kind

You can test the CRD locally using [kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker):

### Quick Test

Run the automated test script:

```bash
cd crd
./test-with-kind.sh
```

This script will:
1. Install kind and kubectl if needed
2. Create a local kind cluster
3. Install the CRD
4. Create example projects
5. Test validation rules
6. Verify everything works

### Manual Testing

1. **Create a kind cluster:**
   ```bash
   kind create cluster --name 0xhub-test
   ```

2. **Install the CRD:**
   ```bash
   kubectl apply -f crd/project.yaml
   kubectl wait --for condition=established crd/projects.hub.0xhub.io
   ```

3. **Create a test project:**
   ```bash
   kubectl apply -f crd/examples/example-project.yaml
   kubectl get projects
   ```

4. **Clean up:**
   ```bash
   kind delete cluster --name 0xhub-test
   ```

## Schema

### Required Fields

- `spec.name` (string): The name of the project (1-100 characters)
- `spec.description` (string): A description of the project (1-1000 characters)
- `spec.url` (string): The URL of the project (must be a valid URI)

### Optional Fields

- `spec.icon` (string): URL or path to the project icon (must be a valid URI)
- `spec.category` (string): The category or group this project belongs to (1-50 characters)
- `spec.status` (string): The status of the project. Valid values:
  - `active` (default)
  - `inactive`
  - `archived`
  - `maintenance`

### Status Fields (Managed by Operator)

- `status.synced` (boolean): Whether the project has been synced to the backend
- `status.lastSyncedAt` (string): Timestamp of the last successful sync
- `status.error` (string): Error message if sync failed

## Examples

### Basic Project

```yaml
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: my-project
  namespace: default
spec:
  name: My Awesome Project
  description: This is a description of my project
  url: https://myproject.com
```

### Full Project with All Fields

```yaml
apiVersion: hub.0xhub.io/v1
kind: Project
metadata:
  name: full-project
  namespace: default
spec:
  name: Full Featured Project
  description: A project with all optional fields filled in
  url: https://fullproject.com
  icon: https://fullproject.com/icon.png
  category: infrastructure
  status: active
```

See the `examples/` directory for more examples.

## Usage

### Create a Project

```bash
kubectl apply -f examples/example-project.yaml
```

### List Projects

```bash
kubectl get projects
# or using short name
kubectl get proj
```

### Get Project Details

```bash
kubectl get project example-web-app -o yaml
```

### Update a Project

Edit the resource and apply:

```bash
kubectl edit project example-web-app
# or
kubectl apply -f examples/example-project.yaml
```

### Delete a Project

```bash
kubectl delete project example-web-app
```

## Validation

The CRD includes OpenAPI schema validation that enforces:

- Required fields must be present
- String length constraints
- URL format validation for `url` and `icon` fields
- Enum validation for `status` field
- Type checking for all fields

Invalid resources will be rejected by the Kubernetes API server.

## Operator Integration

Once the operator is running (Phase 3), it will:

1. Watch for Project CRD create/update/delete events
2. Sync the project data to the backend API
3. Update the `status` subresource with sync information

Check the project status to see if it has been synced:

```bash
kubectl get project example-web-app -o jsonpath='{.status}'
```

## Troubleshooting

### CRD Not Found

If you get an error that the CRD doesn't exist, make sure you've installed it:

```bash
kubectl apply -f crd/project.yaml
```

### Validation Errors

If your resource is rejected, check the validation errors:

```bash
kubectl describe project <project-name>
```

Common issues:
- Missing required fields (`name`, `description`, `url`)
- Invalid URL format
- Status value not in allowed enum values
- String length exceeds maximum

