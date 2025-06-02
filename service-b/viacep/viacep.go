package viacep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrCepNotFound = errors.New("CEP not found")

type viacepResponse struct {
	Localidade string `json:"localidade"`
}

func FetchCity(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrCepNotFound
	}

	var data viacepResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.Localidade == "" {
		return "", ErrCepNotFound
	}

	return data.Localidade, nil
}
