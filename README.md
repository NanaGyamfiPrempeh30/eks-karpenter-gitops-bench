# EKS Karpenter GitOps Bench

[![CI - Tenki Cloud](https://github.com/NanaGyamfiPrempeh30/eks-karpenter-gitops-bench/actions/workflows/ci.yml/badge.svg)](https://github.com/NanaGyamfiPrempeh30/eks-karpenter-gitops-bench/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-grade Kubernetes platform demonstrating EKS with Karpenter autoscaling, ArgoCD GitOps deployments, and CI/CD pipeline benchmarking between Tenki Cloud and GitHub Actions.

![Architecture Diagram](docs/architecture.png)

## Project Overview

This project showcases a complete DevOps workflow:

- **Infrastructure as Code** - 86 AWS resources managed by Terraform
- **GitOps Deployments** - ArgoCD managing dev and prod environments
- **Intelligent Autoscaling** - Karpenter for node provisioning + HPA for pod scaling
- **Multi-Architecture Builds** - Container images supporting AMD64 and ARM64 (Graviton)
- **CI/CD Benchmarking** - Performance comparison between Tenki Cloud and GitHub Actions

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         GitHub Repository                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │
│  │ app/        │  │ infra/      │  │ gitops/     │                 │
│  │ (Go API)    │  │ (Terraform) │  │ (ArgoCD)    │                 │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘                 │
└─────────┼────────────────┼────────────────┼─────────────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Tenki Cloud Runners                              │
│         tenki-standard-medium-4c-8g (4 vCPU, 8GB RAM)              │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │ Build & Test ──► Docker Multi-Arch ──► Push to GHCR         │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
          │
          ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        AWS EKS Cluster                              │
│                    eks-karpenter-bench (eu-west-1)                  │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │                      Control Plane                            │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐              │  │
│  │  │ Karpenter  │  │  ArgoCD    │  │  Metrics   │              │  │
│  │  │            │  │  v2.9.3    │  │  Server    │              │  │
│  │  └────────────┘  └────────────┘  └────────────┘              │  │
│  └──────────────────────────────────────────────────────────────┘  │
│  ┌─────────────────────┐  ┌─────────────────────┐                  │
│  │   Dev Environment   │  │  Prod Environment   │                  │
│  │   1 replica         │  │  2-10 replicas (HPA)│                  │
│  │   ClusterIP         │  │  LoadBalancer       │                  │
│  └─────────────────────┘  └─────────────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
eks-karpenter-gitops-bench/
├── app/                          # Go REST API application
│   ├── main.go                   # Application source code
│   ├── main_test.go              # Unit tests
│   ├── go.mod                    # Go module definition
│   └── Dockerfile                # Multi-arch Dockerfile
├── infra/
│   └── terraform/                # Infrastructure as Code
│       ├── main.tf               # Main Terraform configuration
│       ├── variables.tf          # Input variables
│       ├── outputs.tf            # Output values
│       ├── karpenter-config.tf   # Karpenter NodePool & EC2NodeClass
│       └── argocd.tf             # ArgoCD Helm installation
├── gitops/
│   ├── apps/
│   │   ├── dev/                  # Dev environment manifests
│   │   │   ├── namespace.yaml
│   │   │   ├── deployment.yaml
│   │   │   └── service.yaml
│   │   └── prod/                 # Prod environment manifests
│   │       ├── namespace.yaml
│   │       ├── deployment.yaml
│   │       ├── service.yaml
│   │       └── hpa.yaml
│   └── argocd/
│       ├── dev-application.yaml  # ArgoCD Application for dev
│       └── prod-application.yaml # ArgoCD Application for prod
├── .github/
│   └── workflows/
│       └── ci.yml                # CI pipeline (Tenki Cloud)
└── docs/
    └── architecture.png          # Architecture diagram
```

## Quick Start

### Prerequisites

- AWS CLI configured with appropriate credentials
- Terraform >= 1.5.0
- kubectl
- Git

### 1. Clone the Repository

```bash
git clone https://github.com/NanaGyamfiPrempeh30/eks-karpenter-gitops-bench.git
cd eks-karpenter-gitops-bench
```

### 2. Deploy Infrastructure

```bash
cd infra/terraform

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Deploy (takes ~15-20 minutes)
terraform apply
```

### 3. Configure kubectl

```bash
aws eks update-kubeconfig --region eu-west-1 --name eks-karpenter-bench
```

### 4. Verify Deployment

```bash
# Check nodes
kubectl get nodes

# Check Karpenter
kubectl get pods -n karpenter

# Check ArgoCD
kubectl get pods -n argocd
```

### 5. Deploy Applications via GitOps

```bash
# Apply ArgoCD Applications
kubectl apply -f gitops/argocd/dev-application.yaml
kubectl apply -f gitops/argocd/prod-application.yaml

# Watch ArgoCD sync
kubectl get applications -n argocd
```

### 6. Access the API

```bash
# Get the LoadBalancer URL
kubectl get svc -n gitops-bench-prod

# Test the API
curl http://<EXTERNAL-IP>/health
```

## CI/CD Benchmark Results

Performance comparison between Tenki Cloud and GitHub Actions:

| Metric | GitHub Actions | Tenki Cloud |
|--------|---------------|-------------|
| Runner Specs | 2 vCPU, 7GB | 4 vCPU, 8GB |
| Build & Test | 25s | 24s |
| Docker Multi-Arch | 42s | 1m 43s |
| Cost per minute | $0.008 | $0.006 |

**Key Findings:**
- **Migration Experience:** Seamless - only requires changing `runs-on` value
- **Cost Savings:** 25% cheaper per minute
- **Best For:** Larger workloads where extra CPU provides real gains

## Configuration

### Terraform Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `aws_region` | AWS region | `eu-west-1` |
| `cluster_name` | EKS cluster name | `eks-karpenter-bench` |
| `cluster_version` | Kubernetes version | `1.29` |
| `enable_spot_instances` | Enable Spot for cost savings | `true` |
| `enable_graviton` | Enable ARM64 instances | `true` |

### Karpenter NodePool

The default NodePool is configured for:
- Instance categories: c, m, r, t families
- Capacity types: Spot and On-Demand
- Architectures: AMD64 and ARM64
- Instance sizes: medium, large, xlarge, 2xlarge

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | Root endpoint with API info |
| `GET /health` | Health check |
| `GET /live` | Liveness probe |
| `GET /ready` | Readiness probe |
| `GET /metrics` | Application metrics |

### Sample Response

```json
{
  "architecture": "amd64",
  "environment": "prod",
  "hostname": "gitops-bench-api-86598666cc-w7d9b",
  "message": "EKS Karpenter GitOps Bench API",
  "timestamp": "2026-01-21T16:17:00Z",
  "version": "1.0.0"
}
```

## Cleanup

**Important:** Delete Kubernetes-created resources before Terraform destroy.

```bash
# Delete services (releases LoadBalancers)
kubectl delete svc gitops-bench-api -n gitops-bench-prod
kubectl delete svc argocd-server -n argocd

# Wait for AWS to release resources
sleep 60

# Destroy infrastructure
cd infra/terraform
terraform destroy
```

If Terraform destroy fails with dependency errors:

```bash
# Find and delete orphaned Load Balancers
aws elb describe-load-balancers --region eu-west-1
aws elb delete-load-balancer --load-balancer-name <name> --region eu-west-1

# Find and delete orphaned Security Groups
aws ec2 describe-security-groups --region eu-west-1 --filters "Name=vpc-id,Values=<vpc-id>"
aws ec2 delete-security-group --group-id <sg-id> --region eu-west-1

# Retry destroy
terraform destroy
```

## Blog & Documentation

- [Medium Article: EKS + Karpenter + ArgoCD: From Zero to Production GitOps](https://medium.com/@yawgyamfiprempeh27/eks-karpenter-argocd-from-zero-to-production-gitops-with-25-cheaper-ci-cd-pipelines-heres-8bd81d4d657b)
- [LinkedIn Post with Tenki Cloud Review](https://linkedin.com/in/yaw-nana-gyamfi-prempeh)

## Acknowledgments

- [Marina Rivosecchi](https://www.linkedin.com/in/marina-rivosecchi/) at Tenki Cloud for the CI/CD runner evaluation opportunity
- AWS EKS and Karpenter teams for excellent documentation
- ArgoCD community for GitOps best practices

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Yaw Nana Gyamfi Prempeh**
- GitHub: [@NanaGyamfiPrempeh30](https://github.com/NanaGyamfiPrempeh30)
- LinkedIn: [Yaw Nana Gyamfi Prempeh](www.linkedin.com/in/yaw-gyamfi-prempeh-4042a0129)
- Medium: [@yawgyamfiprempeh27](https://medium.com/@yawgyamfiprempeh27)
