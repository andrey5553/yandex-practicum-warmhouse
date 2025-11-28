FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Копируем весь проект
COPY . .

# Восстанавливаем и собираем
RUN dotnet restore "src/TemperatureApi/TemperatureApi.csproj"
RUN dotnet publish "src/TemperatureApi/TemperatureApi.csproj" -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
COPY --from=build /app/publish .
EXPOSE 8081
ENTRYPOINT ["dotnet", "TemperatureApi.dll"]
