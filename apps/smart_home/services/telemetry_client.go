package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TelemetryData represents the data structure for telemetry
type TelemetryData struct {
	DeviceID  int       `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Metrics   []Metric  `json:"metrics"`
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit,omitempty"`
}

type TelemetryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit,omitempty"`
}

type AggregatedTelemetry struct {
	HouseID               int                    `json:"house_id"`
	Period                string                 `json:"period"`
	TotalEnergyConsumption float64                `json:"total_energy_consumption"`
	AverageTemperature    float64                `json:"average_temperature"`
	DeviceCount           int                    `json:"device_count"`
	Metrics               map[string]interface{} `json:"metrics"`
}

// TelemetryService handles communication with the telemetry service
type TelemetryService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTelemetryService creates a new telemetry service client
func NewTelemetryService(baseURL string) *TelemetryService {
	return &TelemetryService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendTelemetry sends telemetry data to the telemetry service
func (s *TelemetryService) SendTelemetry(data TelemetryData) error {
	url := fmt.Sprintf("%s/telemetry", s.BaseURL)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling telemetry data: %w", err)
	}

	resp, err := s.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending telemetry data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetDeviceTelemetry retrieves telemetry data for a specific device
func (s *TelemetryService) GetDeviceTelemetry(deviceID int, metric string, from, to *time.Time, aggregation string) ([]TelemetryPoint, error) {
	url := fmt.Sprintf("%s/telemetry/devices/%d", s.BaseURL, deviceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	if metric != "" {
		q.Add("metric", metric)
	}
	if from != nil {
		q.Add("from", from.Format(time.RFC3339))
	}
	if to != nil {
		q.Add("to", to.Format(time.RFC3339))
	}
	if aggregation != "" {
		q.Add("aggregation", aggregation)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching device telemetry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		DeviceID int             `json:"device_id"`
		Metrics  []TelemetryPoint `json:"metrics"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding telemetry response: %w", err)
	}

	return response.Metrics, nil
}

// GetAggregatedHouseTelemetry retrieves aggregated telemetry for a house
func (s *TelemetryService) GetAggregatedHouseTelemetry(houseID int, period string) (*AggregatedTelemetry, error) {
	url := fmt.Sprintf("%s/telemetry/houses/%d/aggregated", s.BaseURL, houseID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Add("period", period)
	req.URL.RawQuery = q.Encode()

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching aggregated telemetry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var aggregated AggregatedTelemetry
	if err := json.NewDecoder(resp.Body).Decode(&aggregated); err != nil {
		return nil, fmt.Errorf("error decoding aggregated telemetry response: %w", err)
	}

	return &aggregated, nil
}
