package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var (
	templates *template.Template
	tracer    = otel.Tracer("vote-app")
	results   = map[string]int{
		"aws":   0,
		"gcp":   0,
		"azure": 0,
	}
	mu sync.Mutex
)

func initTracer() (*sdktrace.TracerProvider, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://otel-collector.otel:4318"
	}

	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // Disable TLS for local/cluster use
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("vote-app"),
			semconv.DeploymentEnvironment("production"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return tp, nil
}

func main() {
	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	templates = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "show-vote-page")
		defer span.End()
		renderTemplate(ctx, w, "index.html", nil)
	})

	http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		_, span := tracer.Start(r.Context(), "process-vote")
		defer span.End()

		cloud := r.FormValue("cloud")
		span.SetAttributes(attribute.String("vote.choice", cloud))

		mu.Lock()
		results[cloud]++
		mu.Unlock()

		http.Redirect(w, r, "/results", http.StatusSeeOther)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "show-results")
		defer span.End()
		renderTemplate(ctx, w, "results.html", results)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func renderTemplate(ctx context.Context, w http.ResponseWriter, tmpl string, data interface{}) {
	_, span := tracer.Start(ctx, "render-"+tmpl)
	defer span.End()

	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		span.RecordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
