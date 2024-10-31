package validators

import (
	"testing"
)

// TestClassifyWord verifica la correcta clasificación de palabras.
func TestClassifyWord(t *testing.T) {
	testCases := []struct {
		name     string
		word     string
		expected TiposPalabra
	}{
		{"Subject - I", "i", TipoSujeto},
		{"Subject - You", "you", TipoSujeto},
		{"Subject - He", "he", TipoSujeto},
		{"Verb Auxiliary - Am", "am", TipoVerboAuxiliar},
		{"Verb Auxiliary - Is", "is", TipoVerboAuxiliar},
		{"Article - A", "a", TipoArticulo},
		{"Article - The", "the", TipoArticulo},
		{"Preposition - In", "in", TipoPreposicion},
		{"Preposition - To", "to", TipoPreposicion},
		{"Simple Verb - Play", "play", TipoVerboSimple},
		{"Simple Verb - Eats", "eats", TipoVerboSimple},
		{"Progressive Verb - Playing", "playing", TipoVerboProgresivo},
		{"Adjective - Big", "big", TipoAdjetivo},
		{"Adjective - Happy", "happy", TipoAdjetivo},
		{"Complement - Book", "book", TipoComplemento},
		{"Complement - Games", "games", TipoComplemento},
		{"Unknown - Foobar", "foobar", TipoDesconocido},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := clasificarPalabra(tc.word)
			if result != tc.expected {
				t.Errorf("clasificarPalabra(%s) = %v, expected %v", tc.word, result, tc.expected)
			}
		})
	}
}

// TestTransition verifica las transiciones del autómata.
func TestTransition(t *testing.T) {
	testCases := []struct {
		name     string
		words    []string
		expected int
	}{
		{"Subject -> Verb Auxiliary", []string{"i", "am"}, EstadoVerboAuxiliar},
		{"Subject -> Simple Verb", []string{"i", "play"}, EstadoVerboSimple},
		{"Subject -> Verb Auxiliary -> Progressive Verb", []string{"i", "am", "playing"}, EstadoVerboProgresivo},
		{"Subject -> Verb Simple -> Article -> Adjective -> Complement", []string{"i", "play", "a", "big", "book"}, EstadoComplemento},
		{"Subject -> Verb Simple -> Preposition -> Article -> Complement", []string{"i", "play", "in", "the", "house"}, EstadoFinal},
		{"Incomplete Sentence", []string{"i", "play", "a"}, EstadoArticulo},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewAutomata()
			for _, word := range tc.words {
				a.Transicionar(word)
			}
			if a.estado != tc.expected {
				t.Errorf("Transition failed. Expected state: %d, Actual state: %d", tc.expected, a.estado)
			}
		})
	}
}

// TestValidateSentence verifica la validación de oraciones.
func TestValidateSentence(t *testing.T) {
	testCases := []struct {
		name     string
		sentence string
		expected string
	}{
		{"Valid Sentence", "i play a big book", "Válida"},
		{"Valid Sentence with Preposition", "i play in the house", "Válida"},
		{"Invalid Sentence - Missing Verb", "i a big book", "No válida"},
		{"Invalid Sentence - Unknown Word", "i play a big foobar", "No válida"},
		{"Incomplete Sentence", "i play a", "No válida"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status, _ := ValidarOracion(tc.sentence)
			if status != tc.expected {
				t.Errorf("ValidarOracion(%s) = %s, expected %s", tc.sentence, status, tc.expected)
			}
		})
	}
}
