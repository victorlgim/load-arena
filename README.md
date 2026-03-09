# Kube Load Lab

Mini aplicação para testar **concorrência HTTP**, **autoscaling** e **observabilidade** com Kubernetes.

Projeto pessoal para estudar **CKA + Kubernetes em produção**.

---

## ✨ O que essa aplicação faz

API simples com endpoints que simulam carga:

* `/cpu` → trabalho pesado de CPU
* `/io` → latência simulada
* `/mem` → alocação de memória
* `/chaos` → erros aleatórios
* `/metrics` → métricas Prometheus

Serve para testar:

* requests simultâneas
* autoscaling (HPA)
* limits / OOM
* observabilidade (Prometheus + Grafana)

---

## 🧱 Arquitetura

```
User → Service → Pods (API)
                  ↓
              Prometheus
                  ↓
                Grafana
```

---

## 🛠️ Requisitos

* Docker
* kubectl
* kind **ou** minikube
* Helm
* opcional: `hey` ou `k6`

---

## ▶️ Rodando local

### 1️⃣ Criar cluster

```bash
kind create cluster --name load-arena
```

---

### 2️⃣ Build da imagem

```bash
docker build -t load-arena:dev -f deploy/docker/Dockerfile .
kind load docker-image load-arena:dev --name load-arena
```

---

### 3️⃣ Deploy no Kubernetes

```bash
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/configmap.yaml
kubectl apply -f deploy/k8s/deployment.yaml
kubectl apply -f deploy/k8s/service.yaml
```

---

### 4️⃣ Acessar API

```bash
kubectl port-forward -n load-arena svc/load-arena 8080:80
```

Teste:

```bash
curl "http://localhost:8080/cpu?n=50000"
curl "http://localhost:8080/io?delay=200"
curl "http://localhost:8080/mem?mb=50"
```

---

## 📊 Observabilidade

### Instalar Prometheus + Grafana

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm install monitoring prometheus-community/kube-prometheus-stack \
  -n monitoring --create-namespace
```

Acessar Grafana:

```bash
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:80
```

Usuário: `admin`
Senha:

```bash
kubectl --namespace monitoring get secrets monitoring-grafana \
  -o jsonpath="{.data.admin-password}" | base64 -d ; echo
```

---

### Conectar métricas da app

```bash
kubectl apply -f deploy/k8s/servicemonitor.yaml
```

Depois no Grafana use PromQL:

```
sum(rate(http_requests_total[1m]))
```

---

## 🔥 Testar carga

### hey

```bash
hey -n 10000 -c 200 "http://localhost:8080/cpu?n=60000"
```

### k6

```bash
k6 run scripts/load/k6.js
```

Observe:

```bash
kubectl top pods -n load-arena
kubectl get pods -n load-arena -w
```

---

## 📈 Autoscaling (HPA)

```bash
kubectl apply -f deploy/k8s/hpa.yaml
kubectl get hpa -n load-arena
```

Faça carga no `/cpu` e veja os Pods crescerem.

---

## 🧪 Experimentos legais

* Delete um Pod → Kubernetes recria
* Limite memória → teste `/mem`
* Canary deploy com versão nova
* Chaos endpoint para testar erros
  
---

Se quiser contribuir ou sugerir melhorias, fique à vontade 👍

