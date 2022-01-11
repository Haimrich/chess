FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build-env
WORKDIR /app

COPY *.csproj .
RUN dotnet restore

COPY . .
RUN dotnet publish -c Release -o output

#---------

FROM nginx:alpine
WORKDIR /app
COPY --from=build-env /app/output/wwwroot .
RUN chmod o+r -R *
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80