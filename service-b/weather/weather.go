package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func FetchTemperature(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("WEATHER_API_KEY not set")
	}

	escapedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, escapedCity)
	fmt.Println("[FetchTemperature] Requesting:", requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("[FetchTemperature] Error on request:", err)
		return 0, fmt.Errorf("error calling weather api: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[FetchTemperature] Non-200 response: %d\n", resp.StatusCode)
		fmt.Printf("[FetchTemperature] Body: %s\n", string(body))
		return 0, fmt.Errorf("weather api returned status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("[FetchTemperature] Raw response body:", string(body))

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("[FetchTemperature] Error decoding response:", err)
		return 0, fmt.Errorf("error decoding weather api response: %w", err)
	}

	fmt.Println("[FetchTemperature] Parsed temperature:", data.Current.TempC)
	return data.Current.TempC, nil
}
