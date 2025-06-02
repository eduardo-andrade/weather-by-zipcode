package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("service-b")

func FetchTemperature(ctx context.Context, city string) (float64, error) {
	ctx, span := tracer.Start(ctx, "FetchTemperature")
	defer span.End()

	span.SetAttributes(attribute.String("city", city))

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		err := fmt.Errorf("WEATHER_API_KEY not set")
		span.RecordError(err)
		return 0, err
	}

	escapedCity := url.QueryEscape(city)
	requestURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, escapedCity)
	fmt.Println("[FetchTemperature] Requesting:", requestURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		span.RecordError(err)
		fmt.Println("[FetchTemperature] Error creating request:", err)
		return 0, fmt.Errorf("error creating weather api request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		span.RecordError(err)
		fmt.Println("[FetchTemperature] Error on request:", err)
		return 0, fmt.Errorf("error calling weather api: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("weather api returned status %d: %s", resp.StatusCode, string(body))
		span.RecordError(err)
		fmt.Printf("[FetchTemperature] Non-200 response: %d\n", resp.StatusCode)
		fmt.Printf("[FetchTemperature] Body: %s\n", string(body))
		return 0, err
	}

	fmt.Println("[FetchTemperature] Raw response body:", string(body))

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		span.RecordError(err)
		fmt.Println("[FetchTemperature] Error decoding response:", err)
		return 0, fmt.Errorf("error decoding weather api response: %w", err)
	}

	span.SetAttributes(attribute.Float64("temperature_celsius", data.Current.TempC))
	fmt.Println("[FetchTemperature] Parsed temperature:", data.Current.TempC)
	return data.Current.TempC, nil
}
