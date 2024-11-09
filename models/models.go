package models

import (
	"strings"
	"time"
)

// TipoPalabra representa el tipo de palabra con más categorías
type TipoPalabra uint8

const (
	TipoDesconocido TipoPalabra = iota
	TipoSujeto
	TipoVerboSimple
	TipoVerboEstado
	TipoVerboAuxiliar
	TipoVerboModalPasado // Nuevos verbos modales en pasado
	TipoComplemento
	TipoTiempo
	TipoPreposicion
	TipoArticulo
	TipoAdjetivo
	TipoAdverbio
	TipoConjuncion
	TipoPronombre
	TipoPuntuacion
	TipoNegativo       // Nuevas construcciones negativas
	TipoCausaEfecto    // Nuevas frases de causa y efecto
	TipoRespuestaCorta // Respuestas cortas
)

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

// NewConfig crea una nueva configuración con valores por defecto
func NewConfig() *Config {
	return &Config{
		Port:            "8080",
		StaticDir:       "static",
		TemplatesDir:    "templates",
		MaxRequestSize:  1 << 20,          // 1 MB
		ReadTimeout:     5 * time.Second,  // 5 segundos
		WriteTimeout:    10 * time.Second, // 10 segundos
		RequestTimeout:  30 * time.Second, // 30 segundos
		MaxOraciones:    5,                // Número máximo de oraciones
		EnableCORS:      true,             // Habilitar CORS por defecto
		EnableRateLimit: true,             // Habilitar límite de tasa por defecto
	}
}

// ValidadorConfig contiene la configuración del validador
type ValidadorConfig struct {
	MinPalabras    int  // Número mínimo de palabras
	MaxPalabras    int  // Número máximo de palabras
	MaxOraciones   int  // Número máximo de oraciones
	LimpiarEntrada bool // Si se debe limpiar la entrada
}

// NewValidadorConfig crea una nueva instancia de ValidadorConfig con valores por defecto
func NewValidadorConfig() ValidadorConfig {
	return ValidadorConfig{
		MinPalabras:    1,
		MaxPalabras:    50,
		MaxOraciones:   5,
		LimpiarEntrada: true,
	}
}

// Metadata almacena información adicional sobre la palabra
type Metadata struct {
	EsNombrePropio bool
	EsAbreviatura  bool
	EsContraccion  bool
	SubTipo        string
	EsVerboEstado  bool
}

// Palabra representa una palabra con su tipo y metadata adicional
type Palabra struct {
	Tipo     TipoPalabra
	Texto    string
	Original string
	Posicion int
	Metadata Metadata
}

// Token representa un token de entrada con metadata
type Token struct {
	Tipo     TipoPalabra
	Texto    string
	Original string
	Posicion int
	Metadata Metadata
}

// ElementoOracion representa el estado de un elemento dentro de una oración
type ElementoOracion struct {
	Encontrado bool // Cambia a mayúscula para exportar
	Posicion   int  // Cambia a mayúscula para exportar
	Cantidad   int  // Cambia a mayúscula para exportar
}

// ResultadoOracion representa el resultado de la validación de una oración
type ResultadoOracion struct {
	Oracion     string
	EsValida    bool
	Mensaje     string
	Explicacion string
}

// Estadisticas contiene estadísticas sobre las validaciones realizadas
type Estadisticas struct {
	PorcentajeExito float64
	ErroresComunes  map[string]int
	TiposValidos    map[string]int
}

// PageVariables contiene las variables para renderizar la plantilla
type PageVariables struct {
	Oraciones        []ResultadoOracion
	TotalOraciones   int
	OracionesValidas int
	ShowResults      bool
	ErrorMessage     string
	Estadisticas     Estadisticas
}

// Contexto almacena información sobre el contexto de análisis
type Contexto struct {
	PalabraAnterior   string
	PalabraSiguiente  string
	TipoAnterior      TipoPalabra
	TipoSiguiente     TipoPalabra
	PosicionEnOracion int
}

// ErrorAnalisis representa un error durante el análisis
type ErrorAnalisis struct {
	Mensaje  string
	Posicion int
	Contexto string
}

func (e *ErrorAnalisis) Error() string {
	return e.Mensaje
}

// IdentificarTipoPalabra identifica el tipo de palabra en función del texto
func IdentificarTipoPalabra(texto string) TipoPalabra {
	// Convertimos el texto a minúsculas para hacer la comparación insensible a mayúsculas
	texto = strings.ToLower(texto)

	switch texto {
	case "could", "might", "would":
		return TipoVerboModalPasado
	case "did not", "didn't":
		return TipoNegativo
	case "yes", "no":
		return TipoRespuestaCorta
	case "because", "therefore":
		return TipoCausaEfecto
	}
	return TipoDesconocido
}
