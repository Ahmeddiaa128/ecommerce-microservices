# Kubernetes Deployment (Production-Oriented)

This folder contains Kubernetes manifests to run the full platform with gateway-only access and internal service authentication.

## What’s Included

- Namespace
- ConfigMaps and Secrets
- Deployments and Services for API Gateway and services
- StatefulSets for PostgreSQL and Redis (for non-managed environments)
- NetworkPolicies to restrict access
- Optional Ingress for the API Gateway

## Prerequisites

- A Kubernetes cluster (v1.24+)
- kubectl configured for your cluster
- A StorageClass for PVCs
- Container images for each service

## Build and Push Images

Replace the image references in [app-deployments.yaml](app-deployments.yaml) with your registry and tags.

Example images:

- ghcr.io/your-org/ecommerce-api-gateway:latest
- ghcr.io/your-org/ecommerce-user-service:latest
- ghcr.io/your-org/ecommerce-product-service:latest
- ghcr.io/your-org/ecommerce-cart-service:latest
- ghcr.io/your-org/ecommerce-order-service:latest

## Set Secrets

Edit [secrets.yaml](secrets.yaml) and replace the placeholder values:

- JWT_SECRET
- INTERNAL_AUTH_TOKEN
- POSTGRES_PASSWORD
- USER_DB_DSN
- PRODUCT_DB_DSN
- ORDER_DB_DSN

## Deploy

Apply manifests in order:

1. Namespace
2. Secrets + ConfigMaps
3. Data stores (Postgres/Redis) if you are not using managed services
4. Apps + Services
5. NetworkPolicies
6. Ingress (optional)

Recommended commands:

- kubectl apply -f k8s/namespace.yaml
- kubectl apply -f k8s/secrets.yaml
- kubectl apply -f k8s/configmaps.yaml
- kubectl apply -f k8s/data.yaml
- kubectl apply -f k8s/app-deployments.yaml
- kubectl apply -f k8s/app-services.yaml
- kubectl apply -f k8s/networkpolicies.yaml
- kubectl apply -f k8s/ingress.yaml

## Verify

- kubectl -n ecommerce get pods
- kubectl -n ecommerce get svc

## Production Notes

- Prefer managed PostgreSQL and Redis and remove [data.yaml](data.yaml) in that case.
- Set strong secrets and rotate them regularly.
- Use a service mesh (Istio/Linkerd) for mTLS between services.
- Add HPA and PodDisruptionBudgets for critical services.
- Configure observability (OpenTelemetry/Jaeger) and centralized logging.
