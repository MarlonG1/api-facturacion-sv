FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar herramientas necesarias
RUN apk add --no-cache git

# Copiar proyecto
COPY . .

WORKDIR /app
RUN go build -o dte-app cmd/main.go

# Verificar dónde quedó el binario
RUN echo "=== Buscando binario compilado ===" && \
    find / -name "dte-app" -type f 2>/dev/null

# Imagen final
FROM alpine:latest
LABEL authors="Marlon"

WORKDIR /app

# Copiar binario compilado
COPY --from=builder /app/dte-app* /app/ || true
COPY --from=builder /app/cmd/dte-app* /app/ || true
COPY --from=builder /dte-app* /app/ || true

RUN ls -la /app/

RUN apk add --no-cache ca-certificates

# Configuración
COPY .env /app/

RUN chmod +x /app/dte-app || echo "No se encontró el binario para asignar permisos"

EXPOSE 7319

CMD ["sh", "-c", "ls -la /app && exec /app/dte-app"]