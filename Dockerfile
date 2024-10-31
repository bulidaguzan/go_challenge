FROM golang:1.22-alpine

WORKDIR /app

# Instalar swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copiar archivos del proyecto
COPY . .

# Generar documentación Swagger
RUN swag init --parseDepth 1 --parseDependency --parseInternal

# Descargar dependencias y compilar
RUN go mod download
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]