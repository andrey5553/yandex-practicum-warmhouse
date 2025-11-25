using Microsoft.AspNetCore.Mvc;
using TemperatureApi.Models;

namespace TemperatureApi.Controllers;

[ApiController]
[Route("/temperature")]
public class SmartHomeTemperatureController : ControllerBase
{
    private readonly ILogger<SmartHomeTemperatureController> _logger;

    public SmartHomeTemperatureController(ILogger<SmartHomeTemperatureController> logger)
    {
        _logger = logger;
    }

    [HttpGet("{id}")]
    public TemperatureResponse Get(
        [FromRoute(Name = "id")] string sensorId,
        [FromQuery] string? location)
    {
        _logger.LogInformation("Get new temperature for the sensor: {SensorId} with location: {Location} started", 
            sensorId, 
            location);

        Random random = new Random();
        return new TemperatureResponse
        {
            Value = random.Next(1, 100),
            Status = "Ok",
            Location = location ?? "Unknown",
            SensorId = sensorId,
            Timestamp = DateTime.Now
        };



    }
}
