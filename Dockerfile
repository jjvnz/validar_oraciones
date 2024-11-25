# syntax=docker/dockerfile:1.4

# ---- Etapa de construcción (Builder) ----
    FROM golang:1.23-alpine3.20 AS builder

    # Establecer ARGs para versiones
    ARG NODE_VERSION=18
    ARG TAILWIND_VERSION=3.4.15
    ARG POSTCSS_VERSION=8.4.49
    ARG AUTOPREFIXER_VERSION=10.4.20
    
    # Establecer etiquetas
    LABEL maintainer="Validar Oraciones"
    LABEL version="1.0.0"
    LABEL description="Servidor web en Go para validación de oraciones en inglés (pasado simple afirmativo)"
    
    # Variables de entorno para la compilación
    ENV CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64 \
        GO111MODULE=on
    
    WORKDIR /build
    
    # Instalar dependencias del sistema
    RUN apk add --no-cache \
        nodejs \
        npm \
        && rm -rf /var/cache/apk/*
    
    # Copiar archivos de dependencias
    COPY package*.json ./
    COPY go.mod go.sum ./
    
    # Instalar dependencias específicas de Node.js
    RUN npm install -D \
        tailwindcss@${TAILWIND_VERSION} \
        postcss@${POSTCSS_VERSION} \
        autoprefixer@${AUTOPREFIXER_VERSION}
    
    # Instalar dependencias de Go
    RUN go mod download && \
        go mod verify
    
    # Copiar el resto del código fuente
    COPY . .
    
    # Construir CSS con Tailwind
    RUN npm run build:css
    
    # Construir la aplicación Go
    RUN go build -ldflags="-w -s" -o validar_oraciones .
    
    # ---- Etapa de producción ----
    FROM alpine:3.20
    
    # Copiar certificados SSL
    RUN apk --no-cache add ca-certificates tzdata
    
    # Crear usuario no privilegiado
    RUN adduser -D -u 10001 appuser
    
    # Establecer directorio de trabajo
    WORKDIR /app
    
    # Copiar archivos necesarios
    COPY --from=builder --chown=appuser:appuser /build/validar_oraciones .
    COPY --from=builder --chown=appuser:appuser /build/static ./static
    COPY --from=builder --chown=appuser:appuser /build/templates ./templates
    COPY --from=builder --chown=appuser:appuser /build/words.json ./words.json
    
    # Instalar OpenSSL (versión más reciente disponible)
    RUN apk update && \
        apk upgrade && \
        apk add --no-cache openssl && \
        rm -rf /var/cache/apk/*
    
    # Establecer permisos
    RUN chmod 500 validar_oraciones && \
        chmod -R 400 words.json && \
        chmod -R 500 static templates
    
    # Configurar usuario no privilegiado
    USER appuser
    
    # Puerto por defecto (ajústalo según tu aplicación)
    EXPOSE 8080
    
    # Healthcheck
    HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
        CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
    
    # Comando para ejecutar la aplicación
    CMD ["./validar_oraciones"]
    