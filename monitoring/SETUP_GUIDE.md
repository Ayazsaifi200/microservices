# ==============================================================================
# GRAFANA CLOUD INTEGRATION SETUP GUIDE
# ==============================================================================
# Project: microservices-demo
# Date: November 7, 2025
# ==============================================================================

## STEP 1: GET GRAFANA CLOUD CREDENTIALS
## ======================================

1. Open Grafana Cloud Portal:
   https://ayazzssaifi.grafana.net/

2. Create Access Policy Token:
   - Go to: Administration > Access Policies
   - Or direct: https://ayazzssaifi.grafana.net/a/grafana-auth-app/access-policies
   - Click "Create access policy"
   - Name: "microservices-demo-monitoring"
   - Scopes: Select "metrics:write" and "logs:write"
   - Click "Create"
   - Click "Add token"
   - Copy the token (YOU WON'T SEE IT AGAIN!)

3. Get Prometheus Details:
   - Go to: Connections > Data sources
   - Click on "grafanacloud-ayazzssaifi-prom"
   - Note down:
     * URL: https://prometheus-prod-43-prod-ap-south-1.grafana.net/api/prom/push
     * Username: (shown in Basic Auth section - usually a number)

4. Get Loki Details:
   - Go to: Connections > Data sources
   - Click on "grafanacloud-ayazzssaifi-logs"  
   - Note down:
     * URL: Should be something like https://logs-prod-XXX.grafana.net/loki/api/v1/push
     * Username: (same as Prometheus or similar)


## STEP 2: UPDATE KUBERNETES SECRET
## ==================================

Edit the file: monitoring/grafana-cloud-secret.yaml

Replace these placeholders:
- YOUR_PROMETHEUS_USER_ID: The username from Prometheus datasource
- YOUR_ACCESS_TOKEN: The token you created in Step 1
- YOUR_LOKI_USER_ID: The username from Loki datasource  
- YOUR_LOKI_URL: The Loki push URL

Then apply:
kubectl apply -f monitoring/grafana-cloud-secret.yaml


## STEP 3: UPDATE PROMETHEUS CONFIG
## ==================================

File is ready at: monitoring/prometheus-grafana-cloud-values.yaml

Just verify the URLs match your endpoints, then:
helm upgrade prometheus prometheus-community/kube-prometheus-stack \
  -n monitoring \
  -f monitoring/prometheus-grafana-cloud-values.yaml


## STEP 4: UPDATE PROMTAIL CONFIG
## ================================

File is ready at: monitoring/promtail-grafana-cloud-values.yaml

Update with your Loki URL, then:
helm upgrade loki grafana/loki-stack \
  -n monitoring \
  -f monitoring/promtail-grafana-cloud-values.yaml


## STEP 5: VERIFY DATA IN GRAFANA CLOUD
## ======================================

1. Go to: https://ayazzssaifi.grafana.net/explore
2. Select Prometheus datasource
3. Query: up{job="kubernetes-pods"}
4. You should see metrics from your microservices!

5. Switch to Loki datasource
6. Query: {namespace="default"}
7. You should see logs!


## STEP 6: IMPORT DASHBOARDS
## ==========================

Pre-configured dashboards to import:

1. Kubernetes Cluster Monitoring:
   - Dashboard ID: 15757
   - https://grafana.com/grafana/dashboards/15757

2. Kubernetes Pod Monitoring:
   - Dashboard ID: 14205
   - https://grafana.com/grafana/dashboards/14205

3. Node Exporter Full:
   - Dashboard ID: 1860
   - https://grafana.com/grafana/dashboards/1860

How to import:
- Go to: Dashboards > Import
- Enter dashboard ID
- Select your Prometheus datasource
- Click Import


## QUICK VERIFICATION COMMANDS
## =============================

# Check if Prometheus is scraping
kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090

# Visit: http://localhost:9090/targets

# Check if Promtail is running
kubectl get pods -n monitoring -l app.kubernetes.io/name=promtail

# Check application pods
kubectl get pods

# Generate traffic
kubectl logs -f deployment/loadgenerator

# Check metrics endpoint of a service
kubectl port-forward deployment/frontend 8080:8080
# Visit: http://localhost:8080/metrics


## TROUBLESHOOTING
## ===============

If data not showing in Grafana Cloud:

1. Check secret is created:
   kubectl get secret -n monitoring grafana-cloud-credentials

2. Check Prometheus logs:
   kubectl logs -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0

3. Check remote write status:
   kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090
   # Visit: http://localhost:9090/targets
   # Look for "remote write" section

4. Check Promtail logs:
   kubectl logs -n monitoring -l app.kubernetes.io/name=promtail


## BONUS: DATABASE PERSISTENCE
## ============================

For order persistence, we'll add PostgreSQL.
See: docs/DATABASE_SETUP.md (will be created)

==============================================================================
