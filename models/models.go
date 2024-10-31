package models

import "time"

// ResultadoOracion representa el resultado de validar una oración
type ResultadoOracion struct {
	Oracion     string `json:"oracion"`
	EsValida    bool   `json:"esValida"`
	Mensaje     string `json:"mensaje"`
	Explicacion string `json:"explicacion"`
}

// PageVariables contiene todas las variables necesarias para la plantilla
type PageVariables struct {
	Oraciones        []ResultadoOracion
	ErrorMessage     string
	TotalOraciones   int
	OracionesValidas int
	ShowResults      bool
	Estadisticas     Estadisticas
}

// Estadisticas contiene información estadística sobre las oraciones analizadas
type Estadisticas struct {
	PorcentajeExito float64
	ErroresComunes  map[string]int
	TiposValidos    map[string]int
}

// ValidadorConfig contiene la configuración para el validador
type ValidadorConfig struct {
	MaxOraciones   int
	MinPalabras    int
	MaxPalabras    int
	LimpiarEntrada bool
	ModoEstricto   bool
}

// NewValidadorConfig crea una nueva configuración con valores por defecto
func NewValidadorConfig() ValidadorConfig {
	return ValidadorConfig{
		MaxOraciones:   5,
		MinPalabras:    2,
		MaxPalabras:    15,
		LimpiarEntrada: true,
		ModoEstricto:   false,
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
