namespace TelemetryService.Models;

public class TelemetryData
{
    public int DeviceId { get; set; }
    public DateTime? Timestamp { get; set; } = DateTime.UtcNow;
    public List<Metric> Metrics { get; set; } = new();
}

public class Metric
{
    public string Name { get; set; } = string.Empty;
    public double Value { get; set; }
    public string? Unit { get; set; }
}

public class TelemetryPoint
{
    public DateTime Timestamp { get; set; }
    public double Value { get; set; }
    public string? Unit { get; set; }
}

public class AggregatedTelemetry
{
    public int HouseId { get; set; }
    public string Period { get; set; } = string.Empty;
    public double TotalEnergyConsumption { get; set; }
    public double AverageTemperature { get; set; }
    public int DeviceCount { get; set; }
    public Dictionary<string, object> Metrics { get; set; } = new();
}
