using Microsoft.AspNetCore.Mvc;
using TelemetryService.Models;
using TelemetryService.Services.Contracts;

namespace TelemetryService.Controllers;

[ApiController]
[Route("/telemetry")]
public class TelemetryController : ControllerBase
{
    private readonly ITelemetryService _telemetryService;

    public TelemetryController(ITelemetryService telemetryService)
    {
        _telemetryService = telemetryService;
    }

    [HttpPost]
    public async Task<IActionResult> PostTelemetry([FromBody] TelemetryData data)
    {
        await _telemetryService.StoreTelemetryAsync(data);
        return Accepted();
    }

    [HttpGet("devices/{deviceId}")]
    public async Task<IActionResult> GetDeviceTelemetry(
        [FromRoute] int deviceId,
        [FromQuery] string? metric,
        [FromQuery] DateTime? from,
        [FromQuery] DateTime? to,
        [FromQuery] string? aggregation = "raw")
    {
        var telemetry = await _telemetryService.GetDeviceTelemetryAsync(deviceId, metric, from, to, aggregation);

        return Ok(new
        {
            device_id = deviceId,
            metrics = telemetry
        });
    }

    [HttpGet("houses/{houseId}/aggregated")]
    public async Task<IActionResult> GetAggregatedHouseTelemetry(
        [FromRoute] int houseId,
        [FromQuery] string period)
    {
        var validPeriods = new[] { "today", "yesterday", "week", "month" };
        if (!validPeriods.Contains(period.ToLower()))
        {
            return BadRequest($"Invalid period. Must be one of: {string.Join(", ", validPeriods)}");
        }

        var aggregated = await _telemetryService.GetAggregatedHouseTelemetryAsync(houseId, period);
        return Ok(aggregated);
    }
}