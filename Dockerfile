# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Instalar dependencias necesarias
RUN apk add --no-cache git postgresql-client

# Copiar los archivos de dependencias
COPY go.mod ./
# Asumiendo que tienes un go.sum, si no lo tienes, comenta esta línea
# COPY go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o main .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]