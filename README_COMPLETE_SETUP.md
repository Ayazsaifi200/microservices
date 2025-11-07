# 🚀 Microservices Demo - Complete Setup Guide
## Google Cloud Microservices with Grafana Cloud Monitoring & Database Persistence

---

## 📋 Project Overview

This project demonstrates:
✅ **12 Microservices** deployed on Kubernetes (Minikube)
✅ **Prometheus** + **Loki** for metrics and logs collection
✅ **Grafana Cloud** integration for visualization
✅ **PostgreSQL** database for order persistence (BONUS)
✅ **Load Generation** with realistic traffic simulation

---

## 🎯 Current Status

### ✅ Completed Tasks

1. **✅ Kubernetes Cluster**: Minikube with Docker driver
2. **✅ Application Deployment**: All 12 microservices running
3. **✅ Monitoring Stack**: Prometheus + Loki deployed
4. **✅ Database**: PostgreSQL deployed with schema
5. **✅ Code Changes**: Checkout service modified for persistence

### 🔄 Pending Tasks

1. **Update Grafana Cloud Credentials** (see below)
2. **Rebuild Checkout Service** with database code
3. **Import Dashboards** to Grafana Cloud
4. **Verify Data Flow** in Grafana

---

## 🚀 Quick Start Commands

### 1. Check Cluster Status
\`\`\`powershell
kubectl get nodes
kubectl get pods
kubectl get pods -n monitoring
\`\`\`

### 2. Access Application
\`\`\`powershell
# Start port forwarding
minikube service frontend-external

# Or manually:
kubectl port-forward svc/frontend-external 8080:80
# Visit: http://localhost:8080
\`\`\`

### 3. Check Database
\`\`\`powershell
kubectl get pod -l app=postgres
kubectl exec -it <postgres-pod-name> -- psql -U postgres -d ordersdb -c "SELECT * FROM orders;"
\`\`\`

---

## 🔐 Grafana Cloud Integration

### Step 1: Get Your Credentials

1. **Open Grafana Cloud**: https://ayazzssaifi.grafana.net/

2. **Create Access Policy Token**:
   - Go to: **Administration** → **Access Policies**
   - Or: https://ayazzssaifi.grafana.net/a/grafana-auth-app/access-policies
   - Click "**Create access policy**"
   - Name: `microservices-demo-monitoring`
   - Scopes: Select `metrics:write` and `logs:write`
   - Click "**Create**" → "**Add token**"
   - **COPY THE TOKEN** (you won't see it again!)

3. **Get Prometheus Username**:
   - Go to: **Connections** → **Data sources**
   - Click on "**grafanacloud-ayazzssaifi-prom**"
   - Find the **Username** (usually a number like `784995`)

4. **Get Loki URL**:
   - In Data sources, find "**grafanacloud-ayazzssaifi-logs**"
   - Copy the **URL** (e.g., `https://logs-prod-006.grafana.net`)

### Step 2: Update Credentials

Edit file: `monitoring/grafana-cloud-secret.yaml`

Replace these values:
\`\`\`yaml
prometheus-username: "YOUR_PROMETHEUS_USER_ID"     # e.g., "784995"
prometheus-password: "YOUR_ACCESS_TOKEN"          # Token from step 2
loki-username: "YOUR_LOKI_USER_ID"               # Usually same as Prometheus
loki-password: "YOUR_ACCESS_TOKEN"               # Same token
\`\`\`

### Step 3: Apply Configuration

\`\`\`powershell
# Run the automated setup script
.\monitoring\setup-grafana-cloud.ps1
\`\`\`

Or manually:
\`\`\`powershell
# 1. Create secret
kubectl apply -f monitoring/grafana-cloud-secret.yaml

# 2. Upgrade Prometheus
helm upgrade prometheus prometheus-community/kube-prometheus-stack `
  -n monitoring `
  -f monitoring/prometheus-grafana-cloud-values.yaml

# 3. Upgrade Promtail
helm upgrade loki grafana/loki-stack `
  -n monitoring `
  -f monitoring/promtail-grafana-cloud-values.yaml `
  --set promtail.enabled=true
\`\`\`

---

## 🗄️ Database Persistence Layer (BONUS)

### What's Implemented

- **PostgreSQL 15** deployed in Kubernetes
- **2 Tables**: `orders` and `order_items`
- **Automatic schema** initialization
- **Checkout service** saves every order

### Verify Database

\`\`\`powershell
# Get postgres pod name
kubectl get pods -l app=postgres

# Connect to database
kubectl exec -it postgres-xxxxxxxxx-xxxxx -- psql -U postgres -d ordersdb

# Run queries
\du
# List tables
\dt

# View orders
SELECT order_id, user_email, created_at, total_items FROM orders;

# View order details with items
SELECT o.order_id, o.user_email, oi.product_id, oi.quantity
FROM orders o
JOIN order_items oi ON o.order_id = oi.order_id
ORDER BY o.created_at DESC
LIMIT 10;

# Exit
\q
\`\`\`

### Current Status

⚠️ **IMPORTANT**: The code is ready but service needs to be rebuilt with new code.

To rebuild and deploy:
\`\`\`powershell
# Option 1: Using Skaffold (if installed)
skaffold build -f skaffold.yaml --tag=latest

# Option 2: Use pre-built image and apply DB config
kubectl apply -f database/checkoutservice-with-db.yaml
kubectl rollout restart deployment/checkoutservice
\`\`\`

---

## 📊 Grafana Dashboards

### Import These Dashboards

1. **Kubernetes Cluster Monitoring**
   - Dashboard ID: `15757`
   - https://grafana.com/grafana/dashboards/15757

2. **Kubernetes Pods Monitoring**
   - Dashboard ID: `14205`
   - https://grafana.com/grafana/dashboards/14205

3. **Node Exporter Full**
   - Dashboard ID: `1860`
   - https://grafana.com/grafana/dashboards/1860

4. **Loki Logs Dashboard**
   - Dashboard ID: `13639`
   - https://grafana.com/grafana/dashboards/13639

### How to Import

1. Go to: https://ayazzssaifi.grafana.net/dashboards
2. Click "**New**" → "**Import**"
3. Enter Dashboard ID
4. Select your **Prometheus** datasource
5. Click "**Import**"

---

## 🔍 Verification Checklist

### Application
- [ ] All pods running: `kubectl get pods`
- [ ] Frontend accessible: `minikube service frontend-external`
- [ ] Load generator working: `kubectl logs -f deployment/loadgenerator`

### Monitoring
- [ ] Prometheus pods running: `kubectl get pods -n monitoring`
- [ ] Metrics visible locally: `kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090`
- [ ] Visit: http://localhost:9090/targets

### Grafana Cloud
- [ ] Credentials configured: `kubectl get secret -n monitoring grafana-cloud-credentials`
- [ ] Metrics appearing: https://ayazzssaifi.grafana.net/explore (Prometheus datasource)
- [ ] Logs appearing: https://ayazzssaifi.grafana.net/explore (Loki datasource)
- [ ] Dashboards imported and showing data

### Database
- [ ] PostgreSQL running: `kubectl get pod -l app=postgres`
- [ ] Tables created: `kubectl exec -it <postgres-pod> -- psql -U postgres -d ordersdb -c "\dt"`
- [ ] Orders being saved: Check after placing orders via frontend

---

## 🧪 Generate Test Traffic

The load generator is already running and simulating user traffic!

\`\`\`powershell
# Watch load generator logs
kubectl logs -f deployment/loadgenerator

# Manually place an order via frontend
minikube service frontend-external
# Browse products → Add to cart → Checkout
\`\`\`

---

## 📈 Sample Grafana Queries

### Prometheus Queries

\`\`\`promql
# Application uptime
up{job="kubernetes-pods"}

# Request rate
rate(http_requests_total[5m])

# Pod memory usage
container_memory_usage_bytes{namespace="default"}

# Kubernetes pods by namespace
kube_pod_status_phase{namespace="default"}
\`\`\`

### Loki Queries

\`\`\`logql
# All logs from default namespace
{namespace="default"}

# Logs from specific service
{namespace="default", pod=~"frontend.*"}

# Error logs only
{namespace="default"} |= "error" or "ERROR" or "Error"

# Checkout service logs
{namespace="default", pod=~"checkoutservice.*"} |= "PlaceOrder"
\`\`\`

---

## 🐛 Troubleshooting

### Pods not starting
\`\`\`powershell
kubectl describe pod <pod-name>
kubectl logs <pod-name>
\`\`\`

### Prometheus not sending metrics
\`\`\`powershell
# Check Prometheus logs
kubectl logs -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0

# Check remote write status
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
# Visit: http://localhost:9090/targets
\`\`\`

### Database connection issues
\`\`\`powershell
# Check postgres logs
kubectl logs -l app=postgres

# Test connection
kubectl exec -it <postgres-pod> -- pg_isready -U postgres
\`\`\`

### Grafana Cloud not showing data
1. Verify secret is created: `kubectl get secret -n monitoring grafana-cloud-credentials`
2. Check token has correct scopes (`metrics:write`, `logs:write`)
3. Verify URLs match your Grafana Cloud instance
4. Wait 1-2 minutes for initial data sync

---

## 📝 Files Structure

\`\`\`
microservices-demo/
├── kubernetes-manifests/     # Original K8s manifests
├── monitoring/
│   ├── SETUP_GUIDE.md       # Detailed setup instructions
│   ├── grafana-cloud-secret.yaml          # ** UPDATE THIS **
│   ├── prometheus-grafana-cloud-values.yaml
│   ├── promtail-grafana-cloud-values.yaml
│   └── setup-grafana-cloud.ps1            # Automated setup
├── database/
│   ├── postgres-deployment.yaml           # PostgreSQL setup
│   └── checkoutservice-with-db.yaml       # Updated service
└── src/checkoutservice/
    ├── database.go                        # New database logic
    ├── main.go                           # Modified for DB
    └── go.mod                            # Added PostgreSQL driver
\`\`\`

---

## ✅ Submission Checklist

### Required (Guarantees Reply)
- [x] ✅ Microservices deployed on Kubernetes
- [x] ✅ Traffic generation (load generator running)
- [ ] 🔄 Dashboard with metrics visible
- [ ] 🔄 Application logs visible

### Bonus (Guarantees Interview)
- [x] ✅ Database persistence layer implemented
- [x] ✅ PostgreSQL with orders schema
- [x] ✅ Checkout service modified to save orders
- [ ] 🔄 Non-ngrok endpoint (Grafana Cloud URL)

### What to Share

1. **Grafana Dashboard Link**: https://ayazzssaifi.grafana.net/dashboards
2. **Login Credentials**: Already have access (ayazzssaifi account)
3. **GitHub Repo**: This fork with database implementation
4. **Email**: Send to siddarth@drdroid.io

---

## 🎉 Next Steps

1. **Update `monitoring/grafana-cloud-secret.yaml`** with your actual credentials
2. **Run** `.\monitoring\setup-grafana-cloud.ps1`
3. **Import dashboards** from Grafana marketplace
4. **Verify data** appearing in Grafana Cloud
5. **Place test orders** via frontend
6. **Check database** for saved orders
7. **Take screenshots** of dashboards
8. **Submit** to siddarth@drdroid.io

---

## 📧 Contact

For issues or questions:
- **Email**: siddarth@drdroid.io
- **Include**: "[DrDroid Assignment]" in subject line

---

**Good luck! 🚀**

