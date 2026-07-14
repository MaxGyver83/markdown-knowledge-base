# Markdown Knowledge Base

- Backend: REST API written in Go (using go-chi for routing)
- Frontend: HTML + Javascript + CSS
- Database: PostgreSQL
- CI/CD: GitHub Actions
- Container: Docker / kind (Kubernetes in Docker)
- Infrastructure / Deployment: Terraform (AWS) + k3s on EC2

Check this running instance of this project (at AWS): <http://63.178.0.65:30080/>

## Build and run with Docker

```sh
cd infra
docker compose up -d
```

### Stop

```sh
cd infra
docker compose down
```

## Rebuild

Backend for example:

```sh
cd infra
docker compose build backend
```

## Run with kind (Kubernetes in Docker)

### Prerequisites

- Docker installed
- kind installed
- kubectl installed

### Setup

```sh
# 1. Create kind cluster (if not already present)
#    kind-config.yaml maps NodePort 30080 to host port 30080
kind create cluster --name markdownkb --config kubernetes/kind-config.yaml

# 2a. Download Docker images
docker pull ghcr.io/maxgyver83/markdown-knowledge-base/backend:latest
docker pull ghcr.io/maxgyver83/markdown-knowledge-base/frontend:latest

# 2b. If 2a fails, build Docker images yourself
docker build -t ghcr.io/maxgyver83/markdown-knowledge-base/backend:latest backend/markdown-api
docker build -t ghcr.io/maxgyver83/markdown-knowledge-base/frontend:latest frontend

# 3. Load images into kind
kind load docker-image \
  ghcr.io/maxgyver83/markdown-knowledge-base/backend:latest \
  ghcr.io/maxgyver83/markdown-knowledge-base/frontend:latest \
  --name markdownkb

# 4. Generate Kubernetes secret (not tracked in git)
#    Override with: POSTGRES_PASSWORD=mysecret ./kubernetes/make-secret.sh
./kubernetes/make-secret.sh

# 5. Deploy everything
kubectl apply -k kubernetes/

# 6. Wait for Postgres to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n markdownkb --timeout=60s

# 7. Access the frontend
xdg-open http://localhost:30080
```

### Stop (data preserved)

Pause the kind node container without data loss:

```sh
docker stop kind-markdownkb
docker start kind-markdownkb
```

After `docker start`, Pods may need a moment until Postgres is ready.

### Apply config updates

```sh
kubectl apply -k kubernetes/
kubectl rollout restart -n markdownkb deploy/backend
kubectl rollout restart -n markdownkb deploy/frontend
```

ConfigMap changes take effect after the rollout.

### Teardown (data loss!)

```sh
# WARNING: PVC data is destroyed with the kind node container.
kind delete cluster --name markdownkb
```

Optionally backup all data (database + markdown files) before deleting:

```sh
# Backup
kubectl exec -n markdownkb deploy/postgres -- pg_dump -U postgres markdownkb > backup.sql
kubectl exec -n markdownkb deploy/backend -- tar czf - -C /app/data . > markdown-backup.tar.gz

# Restore (after `kind create cluster` + `kubectl apply -k kubernetes/`)
kubectl exec -i -n markdownkb deploy/postgres -- psql -U postgres markdownkb < backup.sql
kubectl exec -i -n markdownkb deploy/backend -- tar xzf - -C /app/data < markdown-backup.tar.gz
```

## Deploy to AWS with Terraform

Deploy a t2.micro EC2 instance (AWS Free Tier) with k3s using Terraform.

### Prerequisites

- AWS CLI installed and configured (`aws configure`)
- Terraform installed

### Setup

```sh
cd terraform

# Create the infrastructure
./deploy.sh

# Or step by step:
terraform init
terraform apply   # type "yes" when prompted
```

After `terraform apply`, the script:

1. Creates: VPC, subnet, internet gateway, security group (SSH + k3s + frontend), EC2 (Debian 12, 8 GB gp3)
2. Extracts the SSH key to `key.pem`
3. Installs k3s on the EC2 instance via SSH
4. Copies the kubeconfig locally as `kubeconfig.yaml`

### Access

```sh
# SSH
ssh -i terraform/key.pem admin@$(terraform -chdir=terraform output -raw public_ip)

# kubectl
export KUBECONFIG=$(pwd)/terraform/kubeconfig.yaml
kubectl get nodes

# Frontend
open http://$(terraform -chdir=terraform output -raw public_ip):30080
```

### Deploy the app

```sh
export KUBECONFIG=$(pwd)/terraform/kubeconfig.yaml
kubectl apply -k kubernetes/
kubectl wait --for=condition=ready pod -l app=postgres -n markdownkb --timeout=60s
```

### Teardown

```sh
cd terraform
terraform destroy   # removes everything: EC2, VPC, SG, key pair
```

### Data persistence notes

- **kind / minikube / k3d:** PVCs live inside the node container. `kind delete cluster` removes the container and all its data. `docker stop/start` preserves data.
- **Cloud (AKS, EKS, GKE):** PVCs use cloud disks (Azure Disk, EBS, Persistent Disk) that survive cluster deletion.
- **VPS with k3s / kubeadm:** Pod data lives directly on the host filesystem. Cluster restarts or pod deletions do not remove data.
