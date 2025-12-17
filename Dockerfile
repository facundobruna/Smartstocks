# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Copiar go.mod y go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Build de la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario
COPY --from=builder /app/main .

# Exponer puerto
EXPOSE 8080

# Comando para ejecutar
CMD ["./main"]