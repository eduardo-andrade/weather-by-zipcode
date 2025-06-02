package main

import (
	"context"
	"log"
	"net/http"

	"github.com/eduardo-andrade/weather-by-zipcode/shared/tracer"
)

func main() {
	shutdown := tracer.InitTracer(context.Background(), "service-a")
	defer shutdown()

	mux := http.NewServeMux()
	mux.HandleFunc("/cep", HandleCEP)

	log.Println("Service A running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
