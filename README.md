
# Implementing OpenTelemetry with Jaeger on Kubernetes for Scalable Application Monitoring

This repository provides a complete setup for integrating OpenTelemetry (OTEL) into your application and visualizing traces in Jaeger. It is designed for easy deployment on Kubernetes, allowing developers to clone, deploy, and integrate observability into real-world applications.

---

## ✨ Why Use OpenTelemetry?

- **Unified Telemetry**: Collect traces, metrics, and logs using a single standard.
- **Vendor-Agnostic**: Compatible with various backends like Jaeger, Prometheus, Tempo, etc.
- **Improved Observability**: Gain deep insights into application performance and user behavior.
- **Open Source and CNCF Maintained**.

---

## 🧱 Architecture

```
[Your App] ──> [OpenTelemetry Collector] ──> [Jaeger] ──> Jaeger UI
```

---

## 📁 Project Structure

```
.
├── app/                            # Sample instrumented application
│   ├── main.py                     # Python app with OTEL tracing
│   └── Dockerfile                  # Dockerfile for building the app image
├── k8s/                            # Kubernetes manifests
│   ├── namespace.yaml              # Creates the 'otel' namespace
│   ├── jaeger.yaml                 # Jaeger all-in-one deployment
│   ├── otel-collector.yaml         # OTEL Collector config and deployment
│   └── app-deployment.yaml         # App deployment + service + ingress
├── README.md                       # Project documentation

```

---

## 🚀 Quick Start Guide

### 1. Clone the Repo

```bash
git clone https://github.com/sajedul5/OpenTelemetry.git
cd OpenTelemetry
```

### 2. Apply Kubernetes Manifests

```bash

```bash
    kubectl apply -f deploy/namespace.yaml
    kubectl apply -f deploy/jaeger.yaml
    kubectl apply -f deploy/otel-collector.yaml
    kubectl apply -f deploy/app-deployment.yaml
    kubectl apply -f deploy/ingress.yaml
```
```

> Make sure your cluster has an ingress controller installed (like NGINX).

---

## 🔧 Sample Application

The app is a basic Python service instrumented using `opentelemetry-sdk` and `opentelemetry-exporter-otlp`.

```python
# main.py
from flask import Flask
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

trace.set_tracer_provider(
    TracerProvider(resource=Resource.create({SERVICE_NAME: "sample-app"}))
)
span_processor = BatchSpanProcessor(OTLPSpanExporter(endpoint="http://otel-collector.otel.svc.cluster.local:4318"))
trace.get_tracer_provider().add_span_processor(span_processor)

app = Flask(__name__)

@app.route("/")
def index():
    with trace.get_tracer(__name__).start_as_current_span("index-span"):
        return "Hello from OTEL Instrumented App!"

app.run(host="0.0.0.0", port=5000)
```

---

## 🖥️ Access Jaeger UI

Edit your `/etc/hosts`:

```
127.0.0.1 jaeger.todo.com todo.com

```

Jaeger UI: http://jaeger.todo.com

App: http://todo.com

---

## 🛠 Best Practices

- Use different pipelines for logs, metrics, and traces.
- Add OTEL sidecar or auto-instrumentation for larger services.
- Secure OTLP endpoints in production.
- Visualize metrics using Prometheus + Grafana if needed.
- Prefer OTEL Collector as a central routing point.

---

## ✅ Recommended for Production

- Enable persistent storage for Jaeger.
- Use `StatefulSet` or Helm for production-grade deployments.
- Configure TLS and authentication.
- Use OpenTelemetry Operator for auto-instrumentation.

---

## 📚 Resources

- [OpenTelemetry](https://opentelemetry.io/)
- [Jaeger](https://www.jaegertracing.io/)
- [OpenTelemetry Python](https://opentelemetry-python.readthedocs.io/)

---

---

## 📄 License

MIT