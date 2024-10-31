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

		// Verbos simples
		"play": TipoVerboSimple, "plays": TipoVerboSimple, "played": TipoVerboSimple,
		"eat": TipoVerboSimple, "eats": TipoVerboSimple, "ate": TipoVerboSimple,
		"go": TipoVerboSimple, "goes": TipoVerboSimple, "went": TipoVerboSimple,
		"like": TipoVerboSimple, "likes": TipoVerboSimple, "liked": TipoVerboSimple,
		"see": TipoVerboSimple, "sees": TipoVerboSimple, "saw": TipoVerboSimple,
		"know": TipoVerboSimple, "knows": TipoVerboSimple, "knew": TipoVerboSimple,
		"visit": TipoVerboSimple, "visited": TipoVerboSimple, "visits": TipoVerboSimple,
		"buy": TipoVerboSimple, "bought": TipoVerboSimple,
		"walk": TipoVerboSimple, "walked": TipoVerboSimple,

		// Verbos progresivos
		"playing": TipoVerboProgresivo, "eating": TipoVerboProgresivo,
		"going": TipoVerboProgresivo, "liking": TipoVerboProgresivo,
		"seeing": TipoVerboProgresivo, "knowing": TipoVerboProgresivo,
		"reading": TipoVerboProgresivo, "watching": TipoVerboProgresivo,

		// Adjetivos
		"big": TipoAdjetivo, "small": TipoAdjetivo, "good": TipoAdjetivo,
		"bad": TipoAdjetivo, "happy": TipoAdjetivo, "sad": TipoAdjetivo,
		"new": TipoAdjetivo, "old": TipoAdjetivo,

		// Complementos
		"book": TipoComplemento, "books": TipoComplemento, "food": TipoComplemento,
		"game": TipoComplemento, "games": TipoComplemento, "music": TipoComplemento,
		"movie": TipoComplemento, "movies": TipoComplemento, "house": TipoComplemento,
		"car": TipoComplemento, "dog": TipoComplemento, "cat": TipoComplemento,
		"school": TipoComplemento, "grandparents": TipoComplemento,
		"store": TipoComplemento, "friends": TipoComplemento, "soccer": TipoComplemento,
	}

	if tipo, existe := palabrasPorTipo[palabra]; existe {
		return tipo
	}
	return TipoDesconocido
}

// Transicionar realiza la transición del autómata según la palabra de entrada
func (a *Automata) Transicionar(palabra string) bool {
	tipoPalabra := clasificarPalabra(palabra)

	// Tabla de transición
	transiciones := map[int]map[TiposPalabra]int{
		EstadoInicio: {
			TipoSujeto: EstadoSujeto,
		},
		EstadoSujeto: {
			TipoVerboAuxiliar:   EstadoVerboAuxiliar,
			TipoVerboSimple:     EstadoVerboSimple,
			TipoVerboProgresivo: EstadoVerboSimple,
		},
		EstadoVerboAuxiliar: {
			TipoVerboProgresivo: EstadoVerboProgresivo,
		},
		EstadoVerboProgresivo: {
			TipoArticulo:    EstadoArticulo,
			TipoAdjetivo:    EstadoAdjetivo,
			TipoComplemento: EstadoComplemento,
			TipoPreposicion: EstadoPreposicion,
		},
		EstadoVerboSimple: {
			TipoArticulo:    EstadoArticulo,
			TipoAdjetivo:    EstadoAdjetivo,
			TipoComplemento: EstadoComplemento,
			TipoPreposicion: EstadoPreposicion,
		},
		EstadoArticulo: {
			TipoAdjetivo:    EstadoAdjetivo,
			TipoComplemento: EstadoComplemento,
		},
		EstadoAdjetivo: {
			TipoComplemento: EstadoComplemento,
		},
		EstadoComplemento: {
			// No hay transiciones desde el estado final
		},
		EstadoPreposicion: {
			TipoArticulo:    EstadoArticulo,
			TipoComplemento: EstadoComplemento,
		},
	}

	if nuevaEstado, existe := transiciones[a.estado][tipoPalabra]; existe {
		a.estado = nuevaEstado
		return true
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
