package main

import (
	"log"
	"net/http"
	"os"
	"time"
	"validar_oraciones/handlers"
	"validar_oraciones/middleware"
	"validar_oraciones/models"
)

// NewConfig crea una nueva configuración con valores por defecto
func NewConfig() *models.Config {
	return &models.Config{
		Port:            "8080",
		StaticDir:       "static",
		TemplatesDir:    "templates",
		MaxRequestSize:  1 << 20, // 1 MB
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		RequestTimeout:  30 * time.Second,
		MaxOraciones:    5,
		EnableCORS:      true,
		EnableRateLimit: true,
	}
}

func main() {
	// Configurar logger
	logger := log.New(os.Stdout, "VALIDATOR: ", log.LstdFlags|log.Lshortfile)

	// Cargar configuración
	config := NewConfig()

	// Crear router
	mux := http.NewServeMux()

	// Configurar el handler de oraciones
	validadorConfig := models.NewValidadorConfig()
	validadorConfig.MaxOraciones = config.MaxOraciones

	oracionHandler, err := handlers.NewOracionHandler(validadorConfig, logger)
	if err != nil {
		logger.Fatal("Error al crear el handler de oraciones:", err)
	}

	// Configurar middleware
	var handler http.Handler = mux

	// Agregar middleware de logging
	handler = middleware.LogRequest(handler, logger)

	// Agregar middleware de recuperación de pánico
	handler = middleware.RecoverPanic(handler, logger)

	// Agregar middleware de timeout
	handler = middleware.Timeout(handler, config.RequestTimeout)

	// Agregar middleware de CORS si está habilitado
	if config.EnableCORS {
		handler = middleware.CORS(handler)
	}

	// Agregar middleware de rate limiting si está habilitado
	if config.EnableRateLimit {
		handler = middleware.RateLimit(handler, 100) // 100 requests per minute
	}

	// Configurar rutas estáticas
	fsHandler := http.FileServer(http.Dir(config.StaticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", cacheControl(fsHandler)))

	// Configurar rutas de la API
	mux.Handle("/", oracionHandler)
	mux.HandleFunc("/api/validar", oracionHandler.HandleAPIValidation)
	mux.HandleFunc("/api/health", handleHealth)

	// Configurar el servidor
	server := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: int(config.MaxRequestSize),
	}

	// Iniciar el servidor
	logger.Printf("Servidor iniciado en el puerto %s", config.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Error al iniciar el servidor:", err)
	}
}

// cacheControl agrega headers de caché para archivos estáticos
func cacheControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Agregar headers de caché para archivos estáticos
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		h.ServeHTTP(w, r)
	})
}

// handleHealth maneja el endpoint de health check
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"OK"}`))
}
