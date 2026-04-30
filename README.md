# Rampaz

Rampaz is an AI-powered assistant for Kubernetes clusters, designed to provide insights and simplify debugging.
It exposes cluster data through gRPC APIs and an intelligent chat interface, enabling users to query cluster state, workloads, and metrics using natural language.
## 🚀 Features
- 📡 gRPC-based Kubernetes data APIs
- 🤖 AI-powered chat interface for cluster insights
- 📊 Cluster overview and namespace summaries
- 📦 Workload inspection (Pods, Deployments, Jobs, etc.)
- 📈 Metrics integration (via Metrics Server)
- 🔐 Safe, minimal data exposure for AI usage

## 🏗️ Architecture
```
User (Browser/UI)
        ↓
rampaz-ai-server (HTTP + AI)
        ↓ gRPC
rampaz-server (Kubernetes API client)
        ↓
Kubernetes Cluster
```
- rampaz-server → Fetches data from Kubernetes
- rampaz-ai-server → Handles chat + AI + API layer

## 🧩 Components
- rampaz-server → gRPC backend (Kubernetes data layer)
- rampaz-ai-server → AI + HTTP server (chat interface)
- proto/ → gRPC service definitions
- deploy/kubernetes/ → Kubernetes manifests

## ☸️ Kubernetes Deployment
```
kubectl apply -f deploy/kubernetes/
```
Access UI:
```
kubectl port-forward -n rampaz svc/rampaz-ai-server 8080:8080
```
Open: http://localhost:8080


## 🔌 gRPC Testing
- kubectl port-forward -n rampaz svc/rampaz-server 50052:50052
- grpcurl -plaintext localhost:50052 list


## 🔐 Security Approach
- Avoid exposing sensitive data (e.g., Secrets)
- Minimal, aggregated responses for AI
- RBAC with read-only permissions


## 📌 Roadmap
- AI tool-calling improvements
- Cluster health scoring
- Workload insights & anomaly detection
- Multi-session chat support
- Observability dashboards

## 🤝 Contributing
Contributions are welcome!

## 📄 License
MIT License
