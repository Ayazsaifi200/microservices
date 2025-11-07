# Grafana Cloud Integration Script
# Run this after updating grafana-cloud-secret.yaml with your credentials

Write-Host " Grafana Cloud Integration Setup" -ForegroundColor Cyan
Write-Host "====================================`n" -ForegroundColor Cyan

# Step 1: Apply secret
Write-Host " Step 1: Creating Grafana Cloud credentials secret..." -ForegroundColor Yellow
kubectl apply -f monitoring/grafana-cloud-secret.yaml

if ($LASTEXITCODE -eq 0) {
    Write-Host " Secret created successfully`n" -ForegroundColor Green
} else {
    Write-Host " Failed to create secret. Please check the file.`n" -ForegroundColor Red
    exit 1
}

# Step 2: Upgrade Prometheus with Grafana Cloud config
Write-Host " Step 2: Upgrading Prometheus with Grafana Cloud remote write..." -ForegroundColor Yellow
helm upgrade prometheus prometheus-community/kube-prometheus-stack `
    -n monitoring `
    -f monitoring/prometheus-grafana-cloud-values.yaml

if ($LASTEXITCODE -eq 0) {
    Write-Host " Prometheus upgraded successfully`n" -ForegroundColor Green
} else {
    Write-Host " Failed to upgrade Prometheus`n" -ForegroundColor Red
    exit 1
}

# Step 3: Upgrade Promtail for log shipping
Write-Host " Step 3: Upgrading Promtail for Grafana Cloud Loki..." -ForegroundColor Yellow
helm upgrade loki grafana/loki-stack `
    -n monitoring `
    -f monitoring/promtail-grafana-cloud-values.yaml `
    --set promtail.enabled=true `
    --set loki.persistence.enabled=false

if ($LASTEXITCODE -eq 0) {
    Write-Host " Promtail upgraded successfully`n" -ForegroundColor Green
} else {
    Write-Host "Failed to upgrade Promtail`n" -ForegroundColor Red
    exit 1
}

# Step 4: Wait for pods to restart
Write-Host " Step 4: Waiting for monitoring pods to restart..." -ForegroundColor Yellow
Start-Sleep -Seconds 10
kubectl rollout status statefulset/prometheus-prometheus-kube-prometheus-prometheus -n monitoring --timeout=300s

# Step 5: Verify setup
Write-Host "`nSetup Complete! Now verify:" -ForegroundColor Green
Write-Host "================================`n" -ForegroundColor Green

Write-Host "1. Check monitoring pods:" -ForegroundColor Cyan
Write-Host "   kubectl get pods -n monitoring`n"

Write-Host "2. Check Prometheus logs:" -ForegroundColor Cyan
Write-Host "   kubectl logs -n monitoring prometheus-prometheus-kube-prometheus-prometheus-0 -f`n"

Write-Host "3. Port-forward Prometheus UI:" -ForegroundColor Cyan
Write-Host "   kubectl port-forward -n monitoring svc/prometheus-kube-prometheus-prometheus 9090:9090`n"
Write-Host "   Then visit: http://localhost:9090/targets`n"

Write-Host "4. Check Grafana Cloud:" -ForegroundColor Cyan
Write-Host "   https://ayazzssaifi.grafana.net/explore`n"

Write-Host "5. View application logs:" -ForegroundColor Cyan
Write-Host "   kubectl logs -f deployment/loadgenerator`n"

Write-Host "Next: Import dashboards from monitoring/SETUP_GUIDE.md" -ForegroundColor Magenta
