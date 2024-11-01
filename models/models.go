package models

import "time"

// NewConfig crea una nueva configuración con valores por defecto
func NewConfig() *Config {
	return &Config{
		Port:            "8080",
		StaticDir:       "static",
		TemplatesDir:    "templates",
		MaxRequestSize:  1 << 20, // 1 MB
		ReadTimeout:     5,       // 5 segundos
		WriteTimeout:    10,      // 10 segundos
		RequestTimeout:  30,      // 30 segundos
		MaxOraciones:    5,       // Número máximo de oraciones
		EnableCORS:      true,    // Habilitar CORS por defecto
		EnableRateLimit: true,    // Habilitar límite de tasa por defecto
	}
}

// ValidadorConfig contiene la configuración del validador
type ValidadorConfig struct {
	MinPalabras    int  // Número mínimo de palabras
	MaxPalabras    int  // Número máximo de palabras
	MaxOraciones   int  // Número máximo de oraciones
	LimpiarEntrada bool // Si se debe limpiar la entrada
}

// ResultadoOracion representa el resultado de la validación de una oración
type ResultadoOracion struct {
	Oracion     string // La oración validada
	EsValida    bool   // Indica si es válida o no
	Mensaje     string // Mensaje de validación
	Explicacion string // Explicación de la validación
}

// Estadisticas contiene estadísticas sobre las validaciones realizadas
type Estadisticas struct {
	PorcentajeExito float64        // Porcentaje de oraciones válidas
	ErroresComunes  map[string]int // Contador de errores comunes
	TiposValidos    map[string]int // Conteo de tipos válidos
}

// PageVariables contiene las variables para renderizar la plantilla
type PageVariables struct {
	Oraciones        []ResultadoOracion // Resultados de la validación
	TotalOraciones   int                // Total de oraciones procesadas
	OracionesValidas int                // Total de oraciones válidas
	ShowResults      bool               // Si se deben mostrar resultados
	ErrorMessage     string             // Mensaje de error si aplica
	Estadisticas     Estadisticas       // Estadísticas de las validaciones
}

// NewValidadorConfig crea una nueva instancia de ValidadorConfig con valores por defecto
func NewValidadorConfig() ValidadorConfig {
	return ValidadorConfig{
		MinPalabras:    1,    // Puedes ajustar este valor según tus requisitos
		MaxPalabras:    50,   // Puedes ajustar este valor según tus requisitos
		MaxOraciones:   5,    // Límite por defecto de oraciones
		LimpiarEntrada: true, // Por defecto, limpiar la entrada
	}
}

// Config contiene la configuración de la aplicación
type Config struct {
	Port            string
	StaticDir       string
	TemplatesDir    string
	MaxRequestSize  int64
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	RequestTimeout  time.Duration
	MaxOraciones    int
	EnableCORS      bool
	EnableRateLimit bool
}
