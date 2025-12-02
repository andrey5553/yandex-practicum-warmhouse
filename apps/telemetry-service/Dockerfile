FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Копируем весь проект
COPY . .

# Восстанавливаем и собираем
RUN dotnet restore "src/telemetry_service/TelemetryService.csproj"
RUN dotnet publish "src/telemetry_service/TelemetryService.csproj" -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
COPY --from=build /app/publish .
EXPOSE 8084
ENTRYPOINT ["dotnet", "TelemetryService.dll"]
