# OpenTelemetry + Jaeger Playground

Learn tracing by running it on your own machine.

This repo has a small **Todo app**. Every time you use the app, it records a **trace** (what happened and how long each step took). You can see those traces in **Jaeger**.

```
Todo app  ──>  OpenTelemetry Collector  ──>  Jaeger (see traces here)
```

---

## Run it locally

**You need:** [Docker](https://docs.docker.com/get-docker/)

```bash
git clone https://github.com/sajedul5/opentelemetry-jaeger-kubernetes.git
cd opentelemetry-jaeger-kubernetes/app
docker compose up -d --build
```

Open in your browser:

- Todo app: http://localhost:8000
- Jaeger: http://localhost:16686

## See your first trace

1. In the **Todo app**, add a few todos, then tick or delete them.
2. In **Jaeger**, choose service **`todo-html-app`** and click **Find Traces**.
3. Click any trace to see each step of that request.

## Stop it

```bash
docker compose down
```

---

## How it works

| Part | File | What it does |
| ---- | ---- | ------------ |
| Todo app | `app/main.py` | FastAPI app. It creates a trace for every request. |
| Collector | `app/otel-collector-config.yaml` | Receives traces from the app and sends them to Jaeger. |
| Jaeger | `app/docker-compose.yml` | Stores traces and shows them in a web UI. |

The tracing code in `app/main.py`:

```python
trace.set_tracer_provider(TracerProvider(resource=Resource.create({SERVICE_NAME: "todo-html-app"})))
trace.get_tracer_provider().add_span_processor(
    BatchSpanProcessor(OTLPSpanExporter(endpoint="http://otel-collector:4318/v1/traces"))
)
FastAPIInstrumentor.instrument_app(app)  # traces every request automatically
```

## Try next

- Change the service name `todo-html-app` in `main.py`, restart, and find it in Jaeger.
- Add your own span:
  ```python
  with trace.get_tracer(__name__).start_as_current_span("my-step"):
      ...
  ```
- Watch the raw trace data: `docker compose logs -f otel-collector`

---

## Bonus: run it on Kubernetes

Use a local cluster such as kind, minikube, or Docker Desktop. Run these from the repo root:

```bash
kubectl apply -f k8s/
kubectl get pods -n otel        # wait until all pods are Running

kubectl port-forward -n otel svc/fastapi-todo 8000:80 &
kubectl port-forward -n otel svc/jaeger-collector 16686:16686 &
```

Open the same URLs as above. To clean up, run `kubectl delete namespace otel`.

> If `kubectl apply -f k8s/` complains that the `otel` namespace doesn't exist, run `kubectl apply -f k8s/namespace.yaml` first, then run it again.

---

## Learn more

- [OpenTelemetry](https://opentelemetry.io/)
- [Jaeger](https://www.jaegertracing.io/)

MIT License
