package validators

import (
	"errors"
	"strings"
)

// Tipos de palabras
type TipoPalabra uint8

const (
	TipoDesconocido TipoPalabra = iota
	TipoSujeto
	TipoVerboSimple
	TipoComplemento
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

// ClasificarPalabra determina el tipo de una palabra
func ClasificarPalabra(palabra string) Palabra {
	palabra = strings.ToLower(strings.TrimSpace(palabra))

	// Mapa de palabras por tipo
	palabrasPorTipo := map[string]TipoPalabra{
		// Sujetos
		"i":            TipoSujeto,
		"you":          TipoSujeto,
		"he":           TipoSujeto,
		"she":          TipoSujeto,
		"it":           TipoSujeto,
		"we":           TipoSujeto,
		"they":         TipoSujeto,
		"john":         TipoSujeto,
		"mary":         TipoSujeto,
		"peter":        TipoSujeto,
		"julia":        TipoSujeto,
		"mike":         TipoSujeto,
		"ann":          TipoSujeto,
		"the dog":      TipoSujeto,
		"the cat":      TipoSujeto,
		"my friend":    TipoSujeto,
		"the teacher":  TipoSujeto,
		"the students": TipoSujeto,

		// Verbos en pasado simple
		"played":     TipoVerboSimple,
		"visited":    TipoVerboSimple,
		"walked":     TipoVerboSimple,
		"talked":     TipoVerboSimple,
		"worked":     TipoVerboSimple,
		"studied":    TipoVerboSimple,
		"watched":    TipoVerboSimple,
		"listened":   TipoVerboSimple,
		"ate":        TipoVerboSimple,
		"went":       TipoVerboSimple,
		"saw":        TipoVerboSimple,
		"bought":     TipoVerboSimple,
		"made":       TipoVerboSimple,
		"read":       TipoVerboSimple,
		"cleaned":    TipoVerboSimple,
		"called":     TipoVerboSimple,
		"finished":   TipoVerboSimple,
		"liked":      TipoVerboSimple,
		"traveled":   TipoVerboSimple,
		"wrote":      TipoVerboSimple,
		"spoke":      TipoVerboSimple,
		"ran":        TipoVerboSimple,
		"swam":       TipoVerboSimple,
		"drank":      TipoVerboSimple,
		"gave":       TipoVerboSimple,
		"took":       TipoVerboSimple,
		"flew":       TipoVerboSimple,
		"thought":    TipoVerboSimple,
		"came":       TipoVerboSimple,
		"found":      TipoVerboSimple,
		"felt":       TipoVerboSimple,
		"broke":      TipoVerboSimple,
		"chose":      TipoVerboSimple,
		"held":       TipoVerboSimple,
		"left":       TipoVerboSimple,
		"taught":     TipoVerboSimple,
		"built":      TipoVerboSimple,
		"sent":       TipoVerboSimple,
		"met":        TipoVerboSimple,
		"lost":       TipoVerboSimple,
		"said":       TipoVerboSimple,
		"slept":      TipoVerboSimple,
		"understood": TipoVerboSimple,
		"wore":       TipoVerboSimple,
		"kept":       TipoVerboSimple,
		"grew":       TipoVerboSimple,
		"threw":      TipoVerboSimple,
		"gained":     TipoVerboSimple,
		"began":      TipoVerboSimple,
		"ended":      TipoVerboSimple,
		"arrived":    TipoVerboSimple,
		"departed":   TipoVerboSimple,
		"founded":    TipoVerboSimple,
		"proved":     TipoVerboSimple,
		"remained":   TipoVerboSimple,
		"attended":   TipoVerboSimple,
		"celebrated": TipoVerboSimple,
		"enjoyed":    TipoVerboSimple,
		"helped":     TipoVerboSimple,
		"created":    TipoVerboSimple,
		"improved":   TipoVerboSimple,
		"discussed":  TipoVerboSimple,
		"explained":  TipoVerboSimple,
		"described":  TipoVerboSimple,
		"answered":   TipoVerboSimple,
		"continued":  TipoVerboSimple,
		"researched": TipoVerboSimple,

		// Complementos
		"football":     TipoComplemento,
		"music":        TipoComplemento,
		"movie":        TipoComplemento,
		"book":         TipoComplemento,
		"school":       TipoComplemento,
		"home":         TipoComplemento,
		"park":         TipoComplemento,
		"store":        TipoComplemento,
		"homework":     TipoComplemento,
		"food":         TipoComplemento,
		"game":         TipoComplemento,
		"tv":           TipoComplemento,
		"party":        TipoComplemento,
		"meeting":      TipoComplemento,
		"friend":       TipoComplemento,
		"family":       TipoComplemento,
		"house":        TipoComplemento,
		"garden":       TipoComplemento,
		"city":         TipoComplemento,
		"beach":        TipoComplemento,
		"restaurant":   TipoComplemento,
		"concert":      TipoComplemento,
		"trip":         TipoComplemento,
		"vacation":     TipoComplemento,
		"project":      TipoComplemento,
		"presentation": TipoComplemento,
		"exercise":     TipoComplemento,
		"lesson":       TipoComplemento,
		"activity":     TipoComplemento,
		"event":        TipoComplemento,
		"test":         TipoComplemento,
		"competition":  TipoComplemento,
		"adventure":    TipoComplemento,
		"challenge":    TipoComplemento,
		"celebration":  TipoComplemento,
		"gathering":    TipoComplemento,
		"ceremony":     TipoComplemento,
		"discussion":   TipoComplemento,
		"session":      TipoComplemento,
		"assignment":   TipoComplemento,
		"work":         TipoComplemento,
		"research":     TipoComplemento,
		"field":        TipoComplemento,
		"tour":         TipoComplemento,
		"exploration":  TipoComplemento,
		"training":     TipoComplemento,
		"seminar":      TipoComplemento,
	}

	if tipo, existe := palabrasPorTipo[palabra]; existe {
		return Palabra{tipo, palabra}
	}
	return Palabra{TipoDesconocido, palabra}
}

// AnalizarLexico recibe una oración y devuelve una lista de tokens
func AnalizarLexico(oracion string) ([]Token, error) {
	if oracion == "" {
		return nil, errors.New("la oración está vacía")
	}

	var tokens []Token
	palabras := strings.Fields(oracion)

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

	// Validación de estructura: al menos un sujeto y un verbo
	tieneSujeto := false
	tieneVerbo := false

	for _, token := range tokens {
		if token.Tipo == TipoSujeto {
			tieneSujeto = true
		} else if token.Tipo == TipoVerboSimple {
			tieneVerbo = true
		}
	}

	if tieneSujeto && tieneVerbo {
		return "Válida", "La oración tiene una estructura válida."
	}

	return "Inválida", "La oración debe contener al menos un sujeto y un verbo."
}
