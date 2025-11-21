# Use the official ASP.NET Core runtime image
FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS base
WORKDIR /app
EXPOSE 80

# Use the SDK image to build the application
FROM mcr.microsoft.com/dotnet/sdk:7.0 AS build
WORKDIR /src
COPY ["YourMvcApp.csproj", "./"]
RUN dotnet restore "./YourMvcApp.csproj"
COPY . .
WORKDIR "/src/."
RUN dotnet build "YourMvcApp.csproj" -c Release -o /app/build

# Publish the application
FROM build AS publish
RUN dotnet publish "YourMvcApp.csproj" -c Release -o /app/publish

# Final stage: run the application
FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "YourMvcApp.dll"]