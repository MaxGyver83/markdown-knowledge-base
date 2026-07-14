# Terraform

Terraform manages the AWS infrastructure for the production deployment of the Markdown Knowledge Base.

## Infrastructure as Code

All infrastructure is defined in `terraform/main.tf` and provisioned by `terraform/deploy.sh`.

## AWS Resources

```hcl
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = { Name = "markdownkb-vpc" }
}

resource "aws_subnet" "main" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  tags = { Name = "markdownkb-subnet" }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
  tags   = { Name = "markdownkb-igw" }
}

resource "aws_security_group" "main" {
  name_prefix = "markdownkb-"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "k3s API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Frontend NodePort"
    from_port   = 30080
    to_port     = 30080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "k3s" {
  ami                    = "ami-0e1bed4f06a3b4636"
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.main.id
  associate_public_ip_address = true
  key_name               = aws_key_pair.main.key_name
  tags                   = { Name = "markdownkb-k3s" }
}
```

## Deployment Script

The `deploy.sh` script automates the full provisioning process:

```sh
#!/bin/bash
set -euo pipefail

terraform init
terraform apply -auto-approve

INSTANCE_IP=$(terraform output -raw public_ip)

# Install k3s on the EC2 instance
ssh -i key.pem "ubuntu@$INSTANCE_IP" << 'EOF'
  sudo apt update && sudo apt install -y curl
  curl -sfL https://get.k3s.io | sh -s - --tls-san "$INSTANCE_IP"
  sudo cat /etc/rancher/k3s/k3s.yaml
EOF

# Copy kubeconfig locally
ssh -i key.pem "ubuntu@$INSTANCE_IP" "sudo cat /etc/rancher/k3s/k3s.yaml" > kubeconfig.yaml
```

## Key Terraform Commands

```sh
# Initialize Terraform
terraform init

# Preview changes
terraform plan

# Apply infrastructure
terraform apply -auto-approve

# Destroy everything
terraform destroy -auto-approve
```

## CI/CD Integration

The GitHub Actions workflow builds Docker images, pushes them to ghcr.io, then SSHes into the EC2 instance and runs `kubectl set image` to update the running deployment on k3s.
