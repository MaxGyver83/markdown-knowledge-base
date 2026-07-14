# Kubernetes

The project uses Kubernetes for orchestration with two environments:

| Environment | Tool | Namespace |
|-------------|------|-----------|
| Local development | kind | markdownkb |
| AWS production | k3s | markdownkb-aws |

## kind (Kubernetes in Docker)

kind runs a single-node Kubernetes cluster inside a Docker container. It is configured via `kubernetes/kind-config.yaml`:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
    extraPortMappings:
      - containerPort: 30080
        hostPort: 30080
```

## Manifest Structure

All Kubernetes resources are organized under the `kubernetes/` directory:

```
kubernetes/
├── namespace.yaml           # markdownkb namespace
├── configmap.yaml           # Backend environment variables
├── secret.yaml              # Database credentials (gitignored)
├── kustomization.yaml       # Groups all resources
├── postgres/
│   ├── pvc.yaml             # Persistent volume claim
│   ├── deployment.yaml      # Single replica with health checks
│   └── service.yaml         # ClusterIP on port 5432
├── backend/
│   ├── pvc.yaml             # Markdown file storage
│   ├── deployment.yaml      # REST API server
│   └── service.yaml         # ClusterIP on port 8080
└── frontend/
    ├── configmap.yaml       # nginx reverse proxy config
    ├── deployment.yaml      # nginx serving static files
    └── service.yaml         # NodePort on port 30080
```

## Key kubectl Commands

```sh
# Deploy to kind
kubectl apply -n markdownkb -k kubernetes/

# Deploy to k3s (AWS)
KUBECONFIG=terraform/kubeconfig.yaml kubectl apply -n markdownkb-aws -k kubernetes/

# Check status
kubectl get pods -A
kubectl get svc -A

# View logs
kubectl logs -n markdownkb -l app=frontend

# Restart a deployment
kubectl rollout restart -n markdownkb deployment/frontend

# Follow pod events
kubectl get events -n markdownkb --watch
```

## Architecture Notes

- **Backend** and **Frontend** use `Deployment` resources for stateless operation
- **PostgreSQL** uses a `Deployment` with a single replica and a `PersistentVolumeClaim`
- **Frontend** uses `NodePort` service type on port 30080 for direct external access
- The `kustomization.yaml` sets the default namespace and lists all resources
- **Namespaces** separate local (markdownkb) from AWS (markdownkb-aws) deployments
- Kind uses the `kind-markdownkb` kubectl context; k3s uses the `default` context
