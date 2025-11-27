package handlers

import (
	"net/http"
	"strconv"
	"time"

	"smarthome/services"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler handles telemetry-related requests
type TelemetryHandler struct {
	TelemetryService *services.TelemetryService
}

// NewTelemetryHandler creates a new TelemetryHandler
func NewTelemetryHandler(telemetryService *services.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{
		TelemetryService: telemetryService,
	}
}

// RegisterRoutes registers the telemetry routes
func (h *TelemetryHandler) RegisterRoutes(router *gin.RouterGroup) {
	telemetry := router.Group("/telemetry")
	{
		telemetry.POST("", h.SendTelemetry)
		telemetry.GET("/devices/:id", h.GetDeviceTelemetry)
		telemetry.GET("/houses/:id/aggregated", h.GetAggregatedHouseTelemetry)
	}
}

// SendTelemetry handles POST /api/v1/telemetry
func (h *TelemetryHandler) SendTelemetry(c *gin.Context) {
	var telemetryData services.TelemetryData
	if err := c.ShouldBindJSON(&telemetryData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set timestamp if not provided
	if telemetryData.Timestamp.IsZero() {
		telemetryData.Timestamp = time.Now()
	}

	err := h.TelemetryService.SendTelemetry(telemetryData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Telemetry data accepted"})
}

// GetDeviceTelemetry handles GET /api/v1/telemetry/devices/:id
func (h *TelemetryHandler) GetDeviceTelemetry(c *gin.Context) {
	deviceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	metric := c.Query("metric")
	aggregation := c.DefaultQuery("aggregation", "raw")

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &parsed
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &parsed
		}
	}

	telemetry, err := h.TelemetryService.GetDeviceTelemetry(deviceID, metric, from, to, aggregation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_id": deviceID,
		"metrics":   telemetry,
	})
}

// GetAggregatedHouseTelemetry handles GET /api/v1/telemetry/houses/:id/aggregated
func (h *TelemetryHandler) GetAggregatedHouseTelemetry(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
		return
	}

	period := c.DefaultQuery("period", "today")
	validPeriods := map[string]bool{
		"today": true, "yesterday": true, "week": true, "month": true,
	}
	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid period. Must be one of: today, yesterday, week, month",
		})
		return
	}

	aggregated, err := h.TelemetryService.GetAggregatedHouseTelemetry(houseID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aggregated)
}
