# ---- Etapa de construcción ----
FROM golang:1.23.2 AS builder

# Establece el directorio de trabajo
WORKDIR /app

# Copia go.mod y go.sum primero
COPY go.mod go.sum ./

# Descarga las dependencias de Go
RUN go mod download

# Copia el resto de los archivos de la aplicación al contenedor
COPY . .

# Instala Node.js y TailwindCSS
RUN apt-get update && apt-get install -y nodejs npm && \
    npm install -g tailwindcss

# Compila los estilos de Tailwind CSS
RUN tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

# Compila la aplicación Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o validar_oraciones .

# ---- Etapa de producción ----
FROM alpine:latest

# Establece el directorio de trabajo
WORKDIR /root/

# Copia el binario de la aplicación y los archivos necesarios desde la etapa de construcción
COPY --from=builder /app/validar_oraciones .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Expone el puerto en el que la aplicación escuchará
EXPOSE 8080

# Comando por defecto para ejecutar la aplicación
CMD ["./validar_oraciones"]
