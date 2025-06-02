package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

type CepInput struct {
	Cep string `json:"cep"`
}

// HandleCEP valida o CEP e envia para o Serviço B.
func HandleCEP(w http.ResponseWriter, r *http.Request) {
	var input CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || !isValidCep(input.Cep) {
		http.Error(w, `{"error":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	// Serializa novamente pois já lemos o body e não podemos reutilizar r.Body diretamente
	jsonBody, err := json.Marshal(input)
	if err != nil {
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(os.Getenv("SERVICE_B_URL")+"/weather", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("error calling service B: %v", err)
		http.Error(w, `{"error":"service B unavailable"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response from service B: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func isValidCep(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
