# Primera etapa: generar la documentación Swagger
FROM golang:1.22-alpine AS swagger

# Instalar git y dependencias necesarias
RUN apk add --no-cache git

# Configurar el entorno para evitar problemas de certificados
ENV GONOSUMDB=* \
    GOSUMDB=off \
    GOPROXY=direct

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Segunda etapa: construir la aplicación
FROM golang:1.22-alpine

WORKDIR /app

# Copiar los archivos necesarios
COPY . .

# Copiar swag desde la primera etapa
COPY --from=swagger /go/bin/swag /go/bin/swag

# Generar documentación Swagger
RUN swag init

# Instalar dependencias y compilar
RUN go mod download
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]