package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Device represents a smart device from device-service
type Device struct {
	ID           int                    `json:"id"`
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	RoomID       string                 `json:"room_id"` // меняем на string
	SerialNumber string                 `json:"serial_number"`
	Status       string                 `json:"status"`
	Configuration map[string]interface{} `json:"configuration"`
	LastSeen     string                 `json:"last_seen"`  // оставляем как строку
	CreatedAt    string                 `json:"created_at"` // оставляем как строку
}

// DeviceCreate represents the data needed to create a new device
type DeviceCreate struct {
	Name         string                 `json:"name"`
	Type         string                 `json:"type"`
	RoomID       string                 `json:"room_id"` // меняем на string
	SerialNumber string                 `json:"serial_number"`
	Configuration map[string]interface{} `json:"configuration"`
}

// DeviceService handles communication with the device-service
type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewDeviceService creates a new device service client
func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetDevices fetches all devices from device-service
func (s *DeviceService) GetDevices() ([]Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices", s.BaseURL)

	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching devices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Временная структура для гибкого парсинга
	var response struct {
		Devices []struct {
			ID           interface{}            `json:"id"`           // interface{} для гибкости
			Name         string                 `json:"name"`
			Type         string                 `json:"type"`
			RoomID       interface{}            `json:"room_id"`      // interface{} для гибкости
			SerialNumber string                 `json:"serial_number"`
			Status       string                 `json:"status"`
			Configuration map[string]interface{} `json:"configuration"`
			LastSeen     string                 `json:"last_seen"`
			CreatedAt    string                 `json:"created_at"`
		} `json:"devices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding devices response: %w", err)
	}

	// Конвертируем в правильный формат
	devices := make([]Device, len(response.Devices))
	for i, d := range response.Devices {
		// Конвертируем ID
		var deviceID int
		switch v := d.ID.(type) {
		case float64:
			deviceID = int(v)
		case int:
			deviceID = v
		default:
			deviceID = 0
		}

		// Конвертируем RoomID в строку
		var roomID string
		switch v := d.RoomID.(type) {
		case float64:
			roomID = strconv.Itoa(int(v))
		case int:
			roomID = strconv.Itoa(v)
		case string:
			roomID = v
		default:
			roomID = ""
		}

		devices[i] = Device{
			ID:           deviceID,
			Name:         d.Name,
			Type:         d.Type,
			RoomID:       roomID,
			SerialNumber: d.SerialNumber,
			Status:       d.Status,
			Configuration: d.Configuration,
			LastSeen:     d.LastSeen,
			CreatedAt:    d.CreatedAt,
		}
	}

	return devices, nil
}

// CreateDevice creates a new device in device-service
func (s *DeviceService) CreateDevice(device DeviceCreate) (*Device, error) {
	url := fmt.Sprintf("%s/api/v1/devices", s.BaseURL)

	jsonData, err := json.Marshal(device)
	if err != nil {
		return nil, fmt.Errorf("error marshaling device data: %w", err)
	}

	resp, err := s.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Временная структура для гибкого парсинга
	var tempDevice struct {
		ID           interface{}            `json:"id"`
		Name         string                 `json:"name"`
		Type         string                 `json:"type"`
		RoomID       interface{}            `json:"room_id"`
		SerialNumber string                 `json:"serial_number"`
		Status       string                 `json:"status"`
		Configuration map[string]interface{} `json:"configuration"`
		LastSeen     string                 `json:"last_seen"`
		CreatedAt    string                 `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tempDevice); err != nil {
		return nil, fmt.Errorf("error decoding device response: %w", err)
	}

	// Конвертируем ID
	var deviceID int
	switch v := tempDevice.ID.(type) {
	case float64:
		deviceID = int(v)
	case int:
		deviceID = v
	default:
		deviceID = 0
	}

	// Конвертируем RoomID в строку
	var roomID string
	switch v := tempDevice.RoomID.(type) {
	case float64:
		roomID = strconv.Itoa(int(v))
	case int:
		roomID = strconv.Itoa(v)
	case string:
		roomID = v
	default:
		roomID = ""
	}

	createdDevice := &Device{
		ID:           deviceID,
		Name:         tempDevice.Name,
		Type:         tempDevice.Type,
		RoomID:       roomID,
		SerialNumber: tempDevice.SerialNumber,
		Status:       tempDevice.Status,
		Configuration: tempDevice.Configuration,
		LastSeen:     tempDevice.LastSeen,
		CreatedAt:    tempDevice.CreatedAt,
	}

	return createdDevice, nil
}

// SendCommand sends a command to a device
func (s *DeviceService) SendCommand(deviceID int, command string, parameters map[string]interface{}) error {
	url := fmt.Sprintf("%s/api/v1/devices/%d/commands", s.BaseURL, deviceID)

	commandData := map[string]interface{}{
		"command":    command,
		"parameters": parameters,
	}

	jsonData, err := json.Marshal(commandData)
	if err != nil {
		return fmt.Errorf("error marshaling command data: %w", err)
	}

	resp, err := s.HTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
