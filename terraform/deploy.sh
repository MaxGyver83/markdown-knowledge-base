#!/bin/sh
set -e

echo "=== Terraform init ==="
terraform init

echo "=== Terraform apply ==="
terraform apply -auto-approve

echo ""
echo "=== Extract SSH key ==="
terraform output -raw ssh_private_key > key.pem
chmod 600 key.pem

IP=$(terraform output -raw public_ip)
echo "Public IP: $IP"

echo ""
echo "=== Wait for SSH ==="
for i in $(seq 1 20); do
  if ssh -i key.pem -o StrictHostKeyChecking=no -o ConnectTimeout=5 admin@$IP whoami 2>/dev/null; then
    echo "SSH ready after ${i} attempts"
    break
  fi
  echo "Waiting for SSH... attempt $i"
  sleep 5
done

echo ""
echo "=== System setup (swap + curl) ==="
ssh -i key.pem -o StrictHostKeyChecking=no admin@$IP << 'ENDSSH'
  sudo apt update && sudo apt install -y curl

  # Swap for t2.micro (1 GB RAM)
  sudo fallocate -l 1G /swapfile
  sudo chmod 600 /swapfile
  sudo mkswap /swapfile
  sudo swapon /swapfile
  echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
ENDSSH

echo ""
echo "=== Install k3s with --tls-san ==="
ssh -i key.pem -o StrictHostKeyChecking=no admin@$IP \
  "curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC='server --tls-san $IP' sh - && sudo chmod 644 /etc/rancher/k3s/k3s.yaml"

echo ""
echo "=== Copy kubeconfig ==="
ssh -i key.pem -o StrictHostKeyChecking=no admin@$IP "cat /etc/rancher/k3s/k3s.yaml" | \
  sed "s/127.0.0.1/$IP/g" > kubeconfig.yaml

echo ""
echo "=== Done ==="
echo "SSH:         ssh -i key.pem admin@$IP"
echo "Kubeconfig:  export KUBECONFIG=$(pwd)/kubeconfig.yaml"
echo "Frontend:    http://$IP:30080"
