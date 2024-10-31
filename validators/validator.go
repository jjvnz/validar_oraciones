package validators

import (
	"fmt"
	"strings"
)

// Estados del autómata
const (
	EstadoInicio = iota
	EstadoSujeto
	EstadoVerboAuxiliar
	EstadoVerboProgresivo
	EstadoVerboSimple
	EstadoArticulo
	EstadoAdjetivo
	EstadoComplemento
	EstadoPreposicion
	EstadoFinal
)

// TiposPalabra define los posibles tipos de palabras
type TiposPalabra uint8

const (
	TipoDesconocido TiposPalabra = iota
	TipoSujeto
	TipoVerboAuxiliar
	TipoVerboProgresivo
	TipoVerboSimple
	TipoArticulo
	TipoAdjetivo
	TipoComplemento
	TipoPreposicion
)

// Automata representa el autómata finito
type Automata struct {
	estado       int
	tiempoVerbal string
}

// NewAutomata crea una nueva instancia del autómata
func NewAutomata() *Automata {
	return &Automata{
		estado:       EstadoInicio,
		tiempoVerbal: "presente",
	}
}

// clasificarPalabra determina el tipo de una palabra
func clasificarPalabra(palabra string) TiposPalabra {
	palabra = strings.ToLower(strings.TrimSpace(palabra))

	// Mapas de palabras válidas por tipo
	palabrasPorTipo := map[string]TiposPalabra{
		// Sujetos
		"i": TipoSujeto, "you": TipoSujeto, "he": TipoSujeto, "she": TipoSujeto,
		"it": TipoSujeto, "we": TipoSujeto, "they": TipoSujeto,

		// Verbos auxiliares
		"am": TipoVerboAuxiliar, "is": TipoVerboAuxiliar, "are": TipoVerboAuxiliar,
		"was": TipoVerboAuxiliar, "were": TipoVerboAuxiliar,
		"have": TipoVerboAuxiliar, "has": TipoVerboAuxiliar, "had": TipoVerboAuxiliar,
		"do": TipoVerboAuxiliar, "does": TipoVerboAuxiliar, "did": TipoVerboAuxiliar,

		// Artículos
		"a": TipoArticulo, "an": TipoArticulo, "the": TipoArticulo,

		// Preposiciones
		"in": TipoPreposicion, "on": TipoPreposicion, "at": TipoPreposicion,
		"to": TipoPreposicion, "for": TipoPreposicion, "with": TipoPreposicion,
	}

	// Verbos en forma simple
	verbosSimples := map[string]bool{
		"play": true, "plays": true, "played": true,
		"eat": true, "eats": true, "ate": true,
		"go": true, "goes": true, "went": true,
		"like": true, "likes": true, "liked": true,
		"see": true, "sees": true, "saw": true,
		"know": true, "knows": true, "knew": true,
		"visit": true, "visited": true, "visits": true,
		"buy": true, "bought": true,
		"walk": true, "walked": true,
	}

	// Verbos en forma progresiva
	verbosProgresivos := map[string]bool{
		"playing": true, "eating": true, "going": true,
		"liking": true, "seeing": true, "knowing": true,
		"reading": true, "watching": true,
	}

	// Adjetivos
	adjetivos := map[string]bool{
		"big": true, "small": true, "good": true, "bad": true,
		"happy": true, "sad": true, "new": true, "old": true,
	}

	// Complementos
	complementos := map[string]bool{
		"book": true, "books": true, "food": true, "game": true,
		"games": true, "music": true, "movie": true, "movies": true,
		"house": true, "car": true, "dog": true, "cat": true,
		"school": true, "grandparents": true, "store": true,
		"friends": true, "soccer": true,
	}

	// Verificar tipo de palabra
	if tipo, existe := palabrasPorTipo[palabra]; existe {
		return tipo
	}
	if verbosSimples[palabra] {
		return TipoVerboSimple
	}
	if verbosProgresivos[palabra] {
		return TipoVerboProgresivo
	}
	if adjetivos[palabra] {
		return TipoAdjetivo
	}
	if complementos[palabra] {
		return TipoComplemento
	}

	return TipoDesconocido
}

// Transicionar realiza la transición del autómata según la palabra de entrada
// Transicionar realiza la transición del autómata según la palabra de entrada
func (a *Automata) Transicionar(palabra string) bool {
	tipoPalabra := clasificarPalabra(palabra)

	switch a.estado {
	case EstadoInicio:
		if tipoPalabra == TipoSujeto {
			a.estado = EstadoSujeto
			return true
		}

	case EstadoSujeto:
		switch tipoPalabra {
		case TipoVerboAuxiliar:
			a.estado = EstadoVerboAuxiliar
			return true
		case TipoVerboSimple, TipoVerboProgresivo:
			a.estado = EstadoVerboSimple
			return true
		}

	case EstadoVerboAuxiliar:
		if tipoPalabra == TipoVerboProgresivo {
			a.estado = EstadoVerboProgresivo
			return true
		}

	case EstadoVerboProgresivo, EstadoVerboSimple:
		switch tipoPalabra {
		case TipoArticulo:
			a.estado = EstadoArticulo
			return true
		case TipoAdjetivo:
			a.estado = EstadoAdjetivo
			return true
		case TipoComplemento:
			a.estado = EstadoComplemento
			return true
		case TipoPreposicion:
			a.estado = EstadoPreposicion
			return true
		}

	case EstadoArticulo:
		switch tipoPalabra {
		case TipoAdjetivo:
			a.estado = EstadoAdjetivo
			return true
		case TipoComplemento:
			a.estado = EstadoComplemento
			return true
		}

	case EstadoAdjetivo:
		if tipoPalabra == TipoComplemento {
			a.estado = EstadoComplemento
			return true
		}

	case EstadoComplemento:
		// Ahora verificamos que, tras un complemento, siempre pasamos al estado final.
		a.estado = EstadoFinal
		return true // Aceptar que ha terminado la oración

	case EstadoPreposicion:
		switch tipoPalabra {
		case TipoArticulo:
			a.estado = EstadoArticulo
			return true
		case TipoComplemento:
			a.estado = EstadoComplemento
			return true
		}
	}

	return false
}

// ValidarOracion valida una oración completa
func ValidarOracion(oracion string) (string, string) {
	a := NewAutomata()
	palabras := strings.Fields(oracion)

	for _, palabra := range palabras {
		if !a.Transicionar(palabra) {
			return "No válida", fmt.Sprintf("Error en la palabra: %s", palabra)
		}
	}

	// Aceptamos la oración si el autómata termina en EstadoFinal o EstadoComplemento
	if a.estado == EstadoFinal || a.estado == EstadoComplemento {
		return "Válida", "La oración cumple con la estructura gramatical"
	}

	return "No válida", "La oración está incompleta"
}
