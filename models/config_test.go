package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Llamamos a NewConfig para obtener una instancia de la configuración
	config := NewConfig()

	// Comparamos los valores esperados con los valores reales
	assert.Equal(t, 5*time.Second, config.ReadTimeout, "ReadTimeout debe ser 5 segundos")
	assert.Equal(t, 10*time.Second, config.WriteTimeout, "WriteTimeout debe ser 10 segundos")
	assert.Equal(t, 30*time.Second, config.RequestTimeout, "RequestTimeout debe ser 30 segundos")

	// Comparación de otros valores
	assert.Equal(t, "8080", config.Port, "El puerto debe ser 8080")
	assert.Equal(t, "static", config.StaticDir, "El directorio estático debe ser 'static'")
	assert.Equal(t, "templates", config.TemplatesDir, "El directorio de plantillas debe ser 'templates'")
	assert.Equal(t, int64(1<<20), config.MaxRequestSize, "El tamaño máximo de la solicitud debe ser 1 MB")
	assert.Equal(t, true, config.EnableCORS, "CORS debe estar habilitado")
	assert.Equal(t, true, config.EnableRateLimit, "El límite de tasa debe estar habilitado")
}

func TestNewValidadorConfig(t *testing.T) {
	// Crear la configuración del validador utilizando la función NewValidadorConfig
	validadorConfig := NewValidadorConfig()

	// Comprobar que los valores predeterminados sean correctos
	if validadorConfig.MinPalabras != 1 {
		t.Errorf("Expected MinPalabras 1, but got %d", validadorConfig.MinPalabras)
	}
	if validadorConfig.MaxPalabras != 50 {
		t.Errorf("Expected MaxPalabras 50, but got %d", validadorConfig.MaxPalabras)
	}
	if validadorConfig.MaxOraciones != 5 {
		t.Errorf("Expected MaxOraciones 5, but got %d", validadorConfig.MaxOraciones)
	}
	if !validadorConfig.LimpiarEntrada {
		t.Errorf("Expected LimpiarEntrada true, but got %v", validadorConfig.LimpiarEntrada)
	}
}

func TestErrorAnalisis_Error(t *testing.T) {
	// Crear una instancia de ErrorAnalisis
	errorAnalisis := &ErrorAnalisis{
		Mensaje:  "An error occurred",
		Posicion: 1,
		Contexto: "context",
	}

	// Verificar si el método Error devuelve el mensaje correcto
	expectedErrorMessage := "An error occurred"
	if errorAnalisis.Error() != expectedErrorMessage {
		t.Errorf("Expected Error() to return '%s', but got '%s'", expectedErrorMessage, errorAnalisis.Error())
	}
}
