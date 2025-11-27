using TelemetryService.Models;
using TelemetryService.Services.Contracts;


namespace TelemetryService.Services;

public class SmartHomeTelemetryService : ITelemetryService
{
    // Временное хранилище (в реальном проекте используйте БД)
    private static readonly List<TelemetryData> _telemetryStorage = new();

    public Task StoreTelemetryAsync(TelemetryData data)
    {
        _telemetryStorage.Add(data);
        Console.WriteLine($"Stored telemetry for device {data.DeviceId} with {data.Metrics.Count} metrics");
        return Task.CompletedTask;
    }

    public Task<List<TelemetryPoint>> GetDeviceTelemetryAsync(int deviceId, string? metric, DateTime? from, DateTime? to, string? aggregation)
    {
        var deviceData = _telemetryStorage
            .Where(t => t.DeviceId == deviceId)
            .SelectMany(t => t.Metrics
                .Where(m => metric == null || m.Name.Equals(metric, StringComparison.OrdinalIgnoreCase))
                .Select(m => new TelemetryPoint
                {
                    Timestamp = t.Timestamp ?? DateTime.UtcNow,
                    Value = m.Value,
                    Unit = m.Unit
                }))
            .Where(t => from == null || t.Timestamp >= from)
            .Where(t => to == null || t.Timestamp <= to)
            .ToList();

        // Простейшая агрегация
        var result = aggregation switch
        {
            "hour" => AggregateByHour(deviceData),
            "day" => AggregateByDay(deviceData),
            "week" => AggregateByWeek(deviceData),
            _ => deviceData
        };

        return Task.FromResult(result);
    }

    public Task<AggregatedTelemetry> GetAggregatedHouseTelemetryAsync(int houseId, string period)
    {
        // Простейшая логика агрегации по дому
        var houseDevices = _telemetryStorage
            .Where(t => t.DeviceId % 10 == houseId % 10) // Упрощенная логика принадлежности к дому
            .ToList();

        var aggregated = new AggregatedTelemetry
        {
            HouseId = houseId,
            Period = period,
            TotalEnergyConsumption = houseDevices
                .SelectMany(t => t.Metrics)
                .Where(m => m.Name.Equals("energy", StringComparison.OrdinalIgnoreCase))
                .Sum(m => m.Value),
            AverageTemperature = houseDevices
                .SelectMany(t => t.Metrics)
                .Where(m => m.Name.Equals("temperature", StringComparison.OrdinalIgnoreCase))
                .Average(m => m.Value),
            DeviceCount = houseDevices.Select(t => t.DeviceId).Distinct().Count(),
            Metrics = new Dictionary<string, object>
            {
                ["active_devices"] = houseDevices.Count,
                ["data_points"] = houseDevices.Sum(t => t.Metrics.Count)
            }
        };

        return Task.FromResult(aggregated);
    }

    private static List<TelemetryPoint> AggregateByHour(List<TelemetryPoint> data)
    {
        return data
            .GroupBy(t => new DateTime(t.Timestamp.Year, t.Timestamp.Month, t.Timestamp.Day, t.Timestamp.Hour, 0, 0))
            .Select(g => new TelemetryPoint
            {
                Timestamp = g.Key,
                Value = g.Average(t => t.Value),
                Unit = g.First().Unit
            })
            .ToList();
    }

    private static List<TelemetryPoint> AggregateByDay(List<TelemetryPoint> data)
    {
        return data
            .GroupBy(t => t.Timestamp.Date)
            .Select(g => new TelemetryPoint
            {
                Timestamp = g.Key,
                Value = g.Average(t => t.Value),
                Unit = g.First().Unit
            })
            .ToList();
    }

    private static List<TelemetryPoint> AggregateByWeek(List<TelemetryPoint> data)
    {
        return data
            .GroupBy(t => new DateTime(t.Timestamp.Year, 1, 1).AddDays((t.Timestamp.DayOfYear / 7) * 7))
            .Select(g => new TelemetryPoint
            {
                Timestamp = g.Key,
                Value = g.Average(t => t.Value),
                Unit = g.First().Unit
            })
            .ToList();
    }
}