package validators

import (
	"reflect"
	"testing"
	"validar_oraciones/models"
)

func TestPreprocesarTexto(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"texto básico", "Hello World", "Hello World"},
		{"múltiples espacios", "Hello   World   Test", "Hello World Test"},
		{"nombres propios", "John visited London yesterday", "John visited London yesterday"},
		{"texto vacío", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := preprocesarTexto(tt.input)
			if resultado != tt.expected {
				t.Errorf("preprocesarTexto() = %v, expected %v", resultado, tt.expected)
			}
		})
	}
}

func TestEsPosibleNombrePropio(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"nombre propio", "John", true},
		{"palabra minúscula", "cat", false},
		{"vacío", "", false},
		{"número inicial", "123test", false},
		{"camelCase", "iPhone", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := esPosibleNombrePropio(tt.input)
			if resultado != tt.expected {
				t.Errorf("esPosibleNombrePropio(%s) = %v, expected %v", tt.input, resultado, tt.expected)
			}
		})
	}
}

func TestClasificarPalabra(t *testing.T) {
	inicializarDiccionario()

	tests := []struct {
		name     string
		palabra  string
		ctx      models.Contexto
		expected models.Palabra
	}{
		{"verbo simple conocido", "played", models.Contexto{PosicionEnOracion: 1},
			models.Palabra{Tipo: models.TipoVerboSimple, Texto: "played", Original: "played", Posicion: 1}},
		{"nombre propio", "John", models.Contexto{PosicionEnOracion: 0},
			models.Palabra{Tipo: models.TipoSujeto, Texto: "john", Original: "John", Posicion: 0, Metadata: models.Metadata{EsNombrePropio: true}}},
		{"palabra después de artículo", "house", models.Contexto{PosicionEnOracion: 1, TipoAnterior: models.TipoArticulo},
			models.Palabra{Tipo: models.TipoComplemento, Texto: "house", Original: "house", Posicion: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := ClasificarPalabra(tt.palabra, tt.ctx)
			if !reflect.DeepEqual(resultado, tt.expected) {
				t.Errorf("ClasificarPalabra() = %+v, expected %+v", resultado, tt.expected)
			}
		})
	}
}

func TestAnalizarLexico(t *testing.T) {
	tests := []struct {
		name         string
		oracion      string
		expectedLen  int
		expectError  bool
		errorMessage string
	}{
		{"oración válida", "John played football", 3, false, ""},
		{"oración vacía", "", 0, true, "la oración está vacía"},
		{"oración con múltiples espacios", "I    played   football    yesterday", 4, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := AnalizarLexico(tt.oracion)

			if tt.expectError {
				if err == nil {
					t.Errorf("AnalizarLexico() expected error but got nil")
				} else if err.Error() != tt.errorMessage {
					t.Errorf("AnalizarLexico() error = %v, expected %v", err.Error(), tt.errorMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("AnalizarLexico() unexpected error = %v", err)
				return
			}

			if len(tokens) != tt.expectedLen {
				t.Errorf("AnalizarLexico() returned %d tokens, expected %d", len(tokens), tt.expectedLen)
			}
		})
	}
}

func TestValidarTokens(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []models.Token
		expectedStatus string
		expectedMsg    string
	}{
		{
			"oración válida",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Válida",
			"La oración tiene una estructura válida en pasado simple afirmativo.",
		},
		{
			"sin tokens",
			[]models.Token{},
			"Inválida",
			"No se encontraron tokens.",
		},
		{
			"falta sujeto",
			[]models.Token{{Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Inválida",
			"Falta el sujeto en la oración.",
		},
		{
			"falta verbo",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Inválida",
			"Falta un verbo en pasado en la oración.",
		},
		{
			"verbo auxiliar presente",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoVerboAuxiliar, Texto: "did"}, {Tipo: models.TipoVerboSimple, Texto: "played"}},
			"Inválida",
			"La oración no debe contener verbos auxiliares.",
		},
		{
			"orden incorrecto sujeto-verbo",
			[]models.Token{{Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoSujeto, Texto: "I"}},
			"Inválida",
			"El verbo debe seguir al sujeto.",
		},
		{
			"complemento antes del verbo",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoComplemento, Texto: "football"}, {Tipo: models.TipoVerboSimple, Texto: "played"}},
			"Inválida",
			"El complemento debe ir después del verbo.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := ValidarTokens(tt.tokens)
			if status != tt.expectedStatus || msg != tt.expectedMsg {
				t.Errorf("ValidarTokens() = (%v, %v), expected (%v, %v)",
					status, msg, tt.expectedStatus, tt.expectedMsg)
			}
		})
	}
}

func TestValidarOracion(t *testing.T) {
	tests := []struct {
		name           string
		oracion        string
		expectedStatus string
		expectedMsg    string
	}{
		{
			"oración válida simple",
			"I played football",
			"Válida",
			"La oración tiene una estructura válida en pasado simple afirmativo.",
		},
		{
			"oración con nombre propio",
			"John visited London yesterday",
			"Válida",
			"La oración tiene una estructura válida en pasado simple afirmativo.",
		},
		{
			"oración vacía",
			"",
			"Inválida",
			"Error en análisis léxico: la oración está vacía",
		},
		{
			"orden incorrecto",
			"football played I",
			"Inválida",
			"El verbo debe seguir al sujeto.",
		},
		{
			"con verbo auxiliar",
			"I did play football",
			"Inválida",
			"La oración no debe contener verbos auxiliares.",
		},
		{
			"verbo auxiliar al inicio",
			"did I play football",
			"Inválida",
			"La oración no debe contener verbos auxiliares.",
		},
		{
			"múltiples verbos auxiliares",
			"I did not play football",
			"Inválida",
			"La oración no debe contener verbos auxiliares.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := ValidarOracion(tt.oracion)
			if status != tt.expectedStatus || msg != tt.expectedMsg {
				t.Errorf("ValidarOracion(%q) = (%v, %v), expected (%v, %v)",
					tt.oracion, status, msg, tt.expectedStatus, tt.expectedMsg)
			}
		})
	}
}

func TestErrorAnalisis_Error(t *testing.T) {
	err := models.ErrorAnalisis{
		Mensaje:  "error de prueba",
		Posicion: 5,
		Contexto: "contexto de prueba",
	}
	expected := "error de prueba"
	if got := err.Error(); got != expected {
		t.Errorf("ErrorAnalisis.Error() = %v, expected %v", got, expected)
	}
}
