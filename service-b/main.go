package main

import (
	"context"
	"log"
	"net/http"

	"github.com/eduardo-andrade/weather-by-zipcode/shared/tracer"
)

func main() {
	shutdown := tracer.InitTracer(context.Background(), "service-b")
	defer shutdown()

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", HandleWeather)

	log.Println("Service B running on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
