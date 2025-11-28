using System.Text.Json.Serialization;

namespace TemperatureApi.Models;

public class TemperatureResponse
{
    [JsonPropertyName("value")]
    public float Value { get; set; }

    [JsonPropertyName("status")]
    public required string Status { get; set; }

    [JsonPropertyName("location")]
    public required string Location { get; set; }

    [JsonPropertyName("sensor_id")]
    public required string SensorId { get; set; }

    [JsonPropertyName("timestamp")]
    public DateTime Timestamp { get; set; }

}
