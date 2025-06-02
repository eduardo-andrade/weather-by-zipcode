package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracerOpel = otel.Tracer("service-a")

type CepInput struct {
	Cep string `json:"cep"`
}

func HandleCEP(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracerOpel.Start(r.Context(), "HandleCEP")
	defer span.End()

	var input CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || !isValidCep(input.Cep) {
		span.RecordError(err)
		http.Error(w, `{"error":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	span.SetAttributes(attribute.String("cep", input.Cep))

	jsonBody, err := json.Marshal(input)
	if err != nil {
		span.RecordError(err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://service-b:8081/weather", bytes.NewBuffer(jsonBody))
	if err != nil {
		span.RecordError(err)
		log.Printf("error creating request to service B: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		log.Printf("error calling service B: %v", err)
		http.Error(w, `{"error":"service B unavailable"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		log.Printf("error reading response from service B: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	span.SetAttributes(attribute.Int("status_code_service_b", resp.StatusCode))

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func isValidCep(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
