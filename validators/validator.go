package validators

import (
	"errors"
	"strings"
	"sync"
)

// TipoPalabra representa el tipo de palabra usando un tipo byte para optimizar memoria
type TipoPalabra uint8

const (
	TipoDesconocido TipoPalabra = iota
	TipoSujeto
	TipoVerboSimple
	TipoVerboAuxiliar
	TipoComplemento
	TipoTiempo // Para palabras que indican tiempo como "yesterday", "last year", etc.
)

// Palabra representa una palabra con su tipo
type Palabra struct {
	Tipo  TipoPalabra
	Texto string
}

// Token representa un token de entrada
type Token struct {
	Tipo  TipoPalabra
	Texto string
}

// diccionario es un singleton thread-safe para el mapa de palabras
var (
	diccionario map[string]TipoPalabra
	once        sync.Once
)

// inicializarDiccionario crea el mapa de palabras una sola vez
func inicializarDiccionario() {
	once.Do(func() {
		diccionario = make(map[string]TipoPalabra, 200) // Aumentada la capacidad inicial

		// Agregar palabras por categorías
		agregarPalabras([]string{
			"i", "you", "he", "she", "it", "we", "they",
			"john", "mary", "peter", "julia", "mike", "ann",
			"the dog", "the cat", "my friend", "the teacher", "the students",
		}, TipoSujeto)

		// Verbos auxiliares y to be en pasado
		agregarPalabras([]string{
			"was", "were", "had", "did", "could", "would", "should",
			"might", "must", "shall", "will", "can", "may",
		}, TipoVerboAuxiliar)

		// Verbos en pasado simple (mantenemos los anteriores y añadimos más comunes)
		agregarPalabras([]string{
			"played", "visited", "walked", "talked", "worked", "studied",
			"watched", "listened", "ate", "went", "saw", "bought", "made",
			"read", "cleaned", "called", "finished", "liked", "traveled",
			"wrote", "spoke", "ran", "swam", "drank", "gave", "took", "flew",
			"thought", "came", "found", "felt", "broke", "chose", "held",
			"left", "taught", "built", "sent", "met", "lost", "said", "slept",
			"understood", "wore", "kept", "grew", "threw", "gained", "began",
			"ended", "arrived", "departed", "founded", "proved", "remained",
			"attended", "celebrated", "enjoyed", "helped", "created", "improved",
			"discussed", "explained", "described", "answered", "continued", "researched",
			"been", "gone", "done", "had", "made", "gotten", "become",
		}, TipoVerboSimple)

		// Complementos (lugares, objetos, etc.)
		agregarPalabras([]string{
			"football", "music", "movie", "book", "school", "home", "park",
			"store", "homework", "food", "game", "tv", "party", "meeting",
			"friend", "family", "house", "garden", "city", "beach", "restaurant",
			"concert", "trip", "vacation", "project", "presentation", "exercise",
			"lesson", "activity", "event", "test", "competition", "adventure",
			"challenge", "celebration", "gathering", "ceremony", "discussion",
			"session", "assignment", "work", "research", "field", "tour",
			"exploration", "training", "seminar", "japan", "headache",
		}, TipoComplemento)

		// Expresiones de tiempo
		agregarPalabras([]string{
			"yesterday", "today", "tomorrow", "last night", "last week",
			"last month", "last year", "ago", "before", "after",
			"in the morning", "in the afternoon", "in the evening",
		}, TipoTiempo)
	})
}

// agregarPalabras es una función auxiliar para agregar palabras al diccionario
func agregarPalabras(palabras []string, tipo TipoPalabra) {
	for _, palabra := range palabras {
		diccionario[palabra] = tipo
	}
}

// ClasificarPalabra determina el tipo de una palabra
func ClasificarPalabra(palabra string) Palabra {
	inicializarDiccionario()
	palabra = strings.ToLower(strings.TrimSpace(palabra))

	if tipo, existe := diccionario[palabra]; existe {
		return Palabra{tipo, palabra}
	}
	return Palabra{TipoDesconocido, palabra}
}

// AnalizarLexico recibe una oración y devuelve una lista de tokens
func AnalizarLexico(oracion string) ([]Token, error) {
	if oracion == "" {
		return nil, errors.New("la oración está vacía")
	}

	palabras := strings.Fields(oracion)
	tokens := make([]Token, 0, len(palabras))

	for _, palabra := range palabras {
		p := ClasificarPalabra(palabra)
		tokens = append(tokens, Token{p.Tipo, p.Texto})
	}

	return tokens, nil
}

// ValidarTokens valida los tokens generados a partir de la oración
func ValidarTokens(tokens []Token) (string, string) {
	if len(tokens) == 0 {
		return "Inválida", "No se encontraron tokens."
	}

	var (
		tieneSujeto bool
		tieneVerbo  bool // Ahora considera tanto verbos simples como auxiliares
	)

	for _, token := range tokens {
		switch token.Tipo {
		case TipoSujeto:
			tieneSujeto = true
		case TipoVerboSimple, TipoVerboAuxiliar: // Considera ambos tipos de verbos
			tieneVerbo = true
		}

		if tieneSujeto && tieneVerbo {
			break
		}
	}

	if tieneSujeto && tieneVerbo {
		return "Válida", "La oración tiene una estructura válida."
	}

	return "Inválida", "La oración debe contener al menos un sujeto y un verbo."
}
