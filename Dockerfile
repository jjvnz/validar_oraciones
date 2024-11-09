# ---- Etapa de construcción ----
FROM golang:1.23.2 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o validar_oraciones .

# ---- Etapa de producción ----
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN adduser -D appuser

WORKDIR /home/appuser/
COPY --from=builder /app/validar_oraciones .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/words.json ./words.json

# Cambiar permisos del binario
RUN chmod 700 validar_oraciones

USER appuser

CMD ["./validar_oraciones"]