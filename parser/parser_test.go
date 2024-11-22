package validators

import (
	"reflect"
	"testing"
	"validar_oraciones/models"
)

// Constants for repeated error messages
const (
	ErrNoTokensFound             = "No tokens found."
	ErrMissingSubject            = "The subject is missing in the sentence."
	ErrMissingPastVerb           = "A past tense verb is missing in the sentence."
	ErrNoAuxiliaryVerbs          = "The sentence should not contain auxiliary verbs."
	ErrVerbFollowsSubject        = "The verb must immediately follow the subject."
	ErrComplementAfterVerb       = "The complement must come after the verb."
	ErrEmptySentence             = "The sentence is empty."
	ErrIncorrectOrderSubjectVerb = "The verb must follow the subject."
	ErrLexicalAnalysisEmpty      = "Error in lexical analysis: the sentence is empty"
)

// TestPreprocesarTexto tests preprocessing of text for various cases
func TestPreprocesarTexto(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic text", "Hello World", "Hello World"},
		{"multiple spaces", "Hello   World   Test", "Hello World Test"},
		{"proper nouns", "John visited London yesterday", "John visited London yesterday"},
		{"empty text", "", ""},
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

// TestEsPosibleNombrePropio tests possible proper name detection
func TestEsPosibleNombrePropio(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"proper name", "John", true},
		{"lowercase word", "cat", false},
		{"empty", "", false},
		{"number starting", "123test", false},
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

// TestClasificarPalabra tests word classification for various contexts
func TestClasificarPalabra(t *testing.T) {
	inicializarDiccionario()

	tests := []struct {
		name     string
		palabra  string
		ctx      models.Contexto
		expected models.Palabra
	}{
		{"known simple verb", "played", models.Contexto{PosicionEnOracion: 1},
			models.Palabra{Tipo: models.TipoVerboSimple, Texto: "played", Original: "played", Posicion: 1}},
		{"proper noun", "John", models.Contexto{PosicionEnOracion: 0},
			models.Palabra{Tipo: models.TipoSujeto, Texto: "john", Original: "John", Posicion: 0, Metadata: models.Metadata{EsNombrePropio: true}}},
		{"word after article", "house", models.Contexto{PosicionEnOracion: 1, TipoAnterior: models.TipoArticulo},
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

// TestAnalizarLexico tests lexical analysis
func TestAnalizarLexico(t *testing.T) {
	tests := []struct {
		name         string
		oracion      string
		expectedLen  int
		expectError  bool
		errorMessage string
	}{
		{"valid sentence", "John played football", 3, false, ""},
		{"empty sentence", "", 0, true, ErrLexicalAnalysisEmpty},
		{"sentence with multiple spaces", "I    played   football    yesterday", 4, false, ""},
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

// TestValidarTokens tests sentence validation logic for various cases
func TestValidarTokens(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []models.Token
		expectedStatus string
		expectedMsg    string
	}{
		{
			"valid sentence",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Valid",
			"The sentence has a valid structure in the simple past affirmative.",
		},
		{
			"no tokens",
			[]models.Token{},
			"Invalid",
			ErrNoTokensFound,
		},
		{
			"missing subject",
			[]models.Token{{Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Invalid",
			ErrMissingSubject,
		},
		{
			"missing verb",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoComplemento, Texto: "football"}},
			"Invalid",
			ErrMissingPastVerb,
		},
		{
			"present auxiliary verb",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoVerboAuxiliar, Texto: "did"}, {Tipo: models.TipoVerboSimple, Texto: "played"}},
			"Invalid",
			ErrNoAuxiliaryVerbs,
		},
		{
			"incorrect subject-verb order",
			[]models.Token{{Tipo: models.TipoVerboSimple, Texto: "played"}, {Tipo: models.TipoSujeto, Texto: "I"}},
			"Invalid",
			ErrIncorrectOrderSubjectVerb,
		},
		{
			"complement before verb",
			[]models.Token{{Tipo: models.TipoSujeto, Texto: "I"}, {Tipo: models.TipoComplemento, Texto: "football"}, {Tipo: models.TipoVerboSimple, Texto: "played"}},
			"Invalid",
			ErrComplementAfterVerb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := ValidarTokens(tt.tokens)

			if status != tt.expectedStatus || msg != tt.expectedMsg {
				t.Errorf("ValidarTokens() = %v, %v, expected %v, %v", status, msg, tt.expectedStatus, tt.expectedMsg)
			}
		})
	}
}

// TestValidarOracion tests sentence validation logic for the entire sentence
func TestValidarOracion(t *testing.T) {
	tests := []struct {
		name           string
		oracion        string
		expectedStatus string
		expectedMsg    string
	}{
		{
			"valid sentence",
			"I played football",
			"Valid",
			"The sentence has a valid structure in the simple past affirmative.",
		},
		{
			"invalid sentence (missing subject)",
			"played football",
			"Invalid",
			ErrMissingSubject,
		},
		{
			"invalid sentence (missing verb)",
			"I football",
			"Invalid",
			ErrMissingPastVerb,
		},
		{
			"invalid sentence (incorrect word order)",
			"football I played",
			"Invalid",
			ErrIncorrectOrderSubjectVerb,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := ValidarOracion(tt.oracion)

			if status != tt.expectedStatus || msg != tt.expectedMsg {
				t.Errorf("ValidarOracion() = %v, %v, expected %v, %v", status, msg, tt.expectedStatus, tt.expectedMsg)
			}
		})
	}
}
