#!/bin/bash
# Test script to install and verify the Project CRD

set -e

echo "Installing Project CRD..."
kubectl apply -f project.yaml

echo "Waiting for CRD to be established..."
kubectl wait --for condition=established --timeout=60s crd/projects.hub.0xhub.io

echo "Verifying CRD installation..."
kubectl get crd projects.hub.0xhub.io

echo "CRD installed successfully!"
echo ""
echo "You can now create Project resources. Example:"
echo "  kubectl apply -f examples/example-project.yaml"

