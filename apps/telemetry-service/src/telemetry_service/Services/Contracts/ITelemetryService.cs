using TelemetryService.Models;

namespace TelemetryService.Services.Contracts;

public interface ITelemetryService
{
    Task StoreTelemetryAsync(TelemetryData data);
    Task<List<TelemetryPoint>> GetDeviceTelemetryAsync(int deviceId, string? metric, DateTime? from, DateTime? to, string? aggregation);
    Task<AggregatedTelemetry> GetAggregatedHouseTelemetryAsync(int houseId, string period);
}
