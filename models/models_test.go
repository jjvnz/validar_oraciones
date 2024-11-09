package models

import (
	"testing"
)

// TestTipoPalabra verifica que los nuevos tipos de palabra sean correctos
func TestTipoPalabra(t *testing.T) {
	tests := []struct {
		name     string
		tipo     TipoPalabra
		expected TipoPalabra
	}{
		{"TipoDesconocido", TipoDesconocido, TipoDesconocido},
		{"TipoSujeto", TipoSujeto, TipoSujeto},
		{"TipoVerboSimple", TipoVerboSimple, TipoVerboSimple},
		{"TipoVerboModalPasado", TipoVerboModalPasado, TipoVerboModalPasado},
		{"TipoNegativo", TipoNegativo, TipoNegativo},
		{"TipoCausaEfecto", TipoCausaEfecto, TipoCausaEfecto},
		{"TipoRespuestaCorta", TipoRespuestaCorta, TipoRespuestaCorta},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tipo; got != tt.expected {
				t.Errorf("TipoPalabra = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestValidarPalabra verifica que los tipos de palabra se identifiquen correctamente
func TestValidarPalabra(t *testing.T) {
	tests := []struct {
		name     string
		texto    string
		expected TipoPalabra
	}{
		{"Verbo Modal Pasado - Could", "could", TipoVerboModalPasado},
		{"Verbo Modal Pasado - Might", "might", TipoVerboModalPasado},
		{"Negación - Did Not", "did not", TipoNegativo},
		{"Respuesta Corta - Yes", "yes", TipoRespuestaCorta},
		{"Causa y Efecto - Because", "because", TipoCausaEfecto},
		{"Palabra Desconocida - Dog", "dog", TipoDesconocido},
		{"Palabra Desconocida - Running", "running", TipoDesconocido},
		{"Palabra Desconocida - Quickly", "quickly", TipoDesconocido},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tipo := IdentificarTipoPalabra(tt.texto)
			if tipo != tt.expected {
				t.Errorf("IdentificarTipoPalabra(%v) = %v, want %v", tt.texto, tipo, tt.expected)
			}
		})
	}
}
