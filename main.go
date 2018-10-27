package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: "metal-tile-dev1",
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", &ochttp.Handler{}))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.StartSpan(ctx, "/backendhellotime")
	defer span.End()

	fmt.Fprintf(w, "Hello Backend %s", time.Now())
}
