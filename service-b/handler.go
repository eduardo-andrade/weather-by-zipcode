package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/eduardo-andrade/weather-by-zipcode/service-b/viacep"
	"github.com/eduardo-andrade/weather-by-zipcode/service-b/weather"
)

type CepInput struct {
	Cep string `json:"cep"`
}

type WeatherResponse struct {
	City string  `json:"city"`
	C    float64 `json:"temp_C"`
	F    float64 `json:"temp_F"`
	K    float64 `json:"temp_K"`
}

// HandleWeather orquestra a resposta de cidade e temperatura a partir do CEP.
func HandleWeather(w http.ResponseWriter, r *http.Request) {
	var input CepInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if len(input.Cep) != 8 {
		http.Error(w, `{"error":"invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}

	city, err := viacep.FetchCity(input.Cep)
	if err != nil {
		if errors.Is(err, viacep.ErrCepNotFound) {
			http.Error(w, `{"error":"can not find zipcode"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"viacep service error"}`, http.StatusInternalServerError)
		}
		return
	}

	celsius, err := weather.FetchTemperature(r.Context(), city)
	if err != nil {
		http.Error(w, `{"error":"weather service error"}`, http.StatusInternalServerError)
		return
	}

	resp := WeatherResponse{
		City: city,
		C:    celsius,
		F:    celsius*1.8 + 32,
		K:    celsius + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
