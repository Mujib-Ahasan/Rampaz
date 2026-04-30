# Rampaz Kubernetes Deployment

This directory contains Kubernetes manifests to deploy **rampaz-server** (gRPC backend) and **rampaz-ai-server** (HTTP/AI interface).

---

## 🚀 Deploy

```bash
kubectl apply -f deploy/kubernetes/
```

verify: 
```bash
kubectl get pods -n rampaz
kubectl get svc -n rampaz
```

## 🔗 Communication
`rampaz-ai-server` → `rampaz-server` via:
     rampaz-server:50052
Configured using RAMPAZ_GRPC_ADDR environment variable.

## 🌐 Access UI
```bash
kubectl port-forward -n rampaz svc/rampaz-ai-server 8080:8080
```
Open in browser: http://localhost:8080

## ⚠️ Notes
- Uses Docker images
- Both services run in rampaz namespace