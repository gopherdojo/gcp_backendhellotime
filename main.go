package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	projectID, err := GetProjectID()
	if err != nil {
		panic(err)
	}
	// Create and register a OpenCensus Stackdriver Trace exporter.
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()}) // defaultでは10,000回に1回のサンプリングになっているが、リクエストが少ないと出てこないので、とりあえず全部出す

	server := &http.Server{
		Addr: ":8080",
		Handler: &ochttp.Handler{
			Handler:        http.DefaultServeMux,
			Propagation:    &propagation.HTTPFormat{},
			FormatSpanName: formatSpanName,
		},
	}

	http.Handle("/", ochttp.WithRouteTag(func() http.Handler { return http.HandlerFunc(handler) }(), "/backendhellotime/"))
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Backend %s", time.Now())
}

func formatSpanName(r *http.Request) string {
	return fmt.Sprintf("/backendhellotime%s", r.URL.Path)
}
