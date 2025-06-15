from fastapi import FastAPI, Request, Form
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.templating import Jinja2Templates

import json, os
from uuid import uuid4

from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry import trace
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter

app = FastAPI()
templates = Jinja2Templates(directory="templates")

# OpenTelemetry
trace.set_tracer_provider(TracerProvider(resource=Resource.create({SERVICE_NAME: "todo-html-app"})))
trace.get_tracer_provider().add_span_processor(
    BatchSpanProcessor(OTLPSpanExporter(endpoint="http://otel-collector:4318/v1/traces"))
)
FastAPIInstrumentor.instrument_app(app)

# DB setup
DB_PATH = "./db/todos.json"
os.makedirs("db", exist_ok=True)
if not os.path.exists(DB_PATH):
    with open(DB_PATH, "w") as f: json.dump([], f)

def read_db(): return json.load(open(DB_PATH))
def write_db(data): json.dump(data, open(DB_PATH, "w"), indent=2)

#app

@app.get("/", response_class=HTMLResponse)
def home(request: Request):
    todos = read_db()
    return templates.TemplateResponse("index.html", {"request": request, "todos": todos})

@app.post("/add")
def add(title: str = Form(...)):
    todos = read_db()
    todos.append({"id": str(uuid4()), "title": title, "completed": False})
    write_db(todos)
    return RedirectResponse("/", status_code=302)

@app.post("/toggle")
def toggle(todo_id: str = Form(...)):
    todos = read_db()
    for todo in todos:
        if todo["id"] == todo_id:
            todo["completed"] = not todo["completed"]
            break
    write_db(todos)
    return RedirectResponse("/", status_code=302)

@app.post("/delete")
def delete(todo_id: str = Form(...)):
    todos = read_db()
    todos = [todo for todo in todos if todo["id"] != todo_id]
    write_db(todos)
    return RedirectResponse("/", status_code=302)
