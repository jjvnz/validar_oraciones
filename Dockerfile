# ---- Etapa de construcción ----
FROM golang:1.23.2 AS builder

WORKDIR /app

# Instalar dependencias para Node.js (que son necesarias para Tailwind CSS)
RUN apk add --no-cache \
    nodejs \
    npm

# Instalar las dependencias de Node.js
COPY package.json package-lock.json ./
RUN npm install

# Copiar el código Go y los archivos estáticos
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compilar Go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o validar_oraciones .

# Ejecutar el build de Tailwind CSS
RUN npm run build:css

# ---- Etapa de producción ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN adduser -D appuser

WORKDIR /home/appuser/

# Copiar el binario compilado, los archivos estáticos (incluyendo el CSS de Tailwind), y otros recursos
COPY --from=builder /app/validar_oraciones . 
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/words.json ./words.json

# Cambiar permisos del binario
RUN chmod 700 validar_oraciones

USER appuser

CMD ["./validar_oraciones"]