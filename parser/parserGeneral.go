package validators

/*
import (
	"strings"
	"sync"
	"unicode"
)

// TipoPalabra representa el tipo de palabra con más categorías
type TipoPalabra uint8

const (
	TipoDesconocido TipoPalabra = iota
	TipoSujeto
	TipoVerboSimple
	TipoVerboAuxiliar
	TipoComplemento
	TipoTiempo
	TipoPreposicion
	TipoArticulo
	TipoAdjetivo
	TipoAdverbio
	TipoConjuncion
	TipoPronombre
	TipoPuntuacion
)

// Palabra representa una palabra con su tipo y metadata adicional
type Palabra struct {
	Tipo     TipoPalabra
	Texto    string
	Original string   // Mantiene la palabra original antes de procesamiento
	Posicion int      // Posición en la oración
	Metadata Metadata // Información adicional sobre la palabra
}

// Metadata almacena información adicional sobre la palabra
type Metadata struct {
	EsNombrePropio bool
	EsAbreviatura  bool
	EsContraccion  bool
	SubTipo        string // Para clasificación más específica
}

// Token representa un token de entrada con metadata
type Token struct {
	Tipo     TipoPalabra
	Texto    string
	Original string
	Posicion int
	Metadata Metadata
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

// Contexto almacena información sobre el contexto de análisis
type Contexto struct {
	PalabraAnterior   string
	PalabraSiguiente  string
	TipoAnterior      TipoPalabra
	TipoSiguiente     TipoPalabra
	PosicionEnOracion int
}

// diccionario es un singleton thread-safe para el mapa de palabras
var (
	diccionario map[string]Palabra
	once        sync.Once
	mu          sync.RWMutex
)

// inicializarDiccionario crea el mapa de palabras una sola vez
func inicializarDiccionario() {
	once.Do(func() {
		diccionario = make(map[string]Palabra, 1000)

		// Sujetos y pronombres
		agregarPalabras([]string{
			"i", "you", "he", "she", "it", "we", "they",
			"me", "him", "her", "us", "them", "myself", "yourself",
			"himself", "herself", "itself", "ourselves", "themselves",
		}, TipoSujeto)

		// Nombres propios comunes (con metadata de nombre propio)
		nombresPropios := []string{
			"john", "mary", "peter", "julia", "mike", "ann",
			"monday", "tuesday", "wednesday", "thursday", "friday",
			"january", "february", "march", "april", "may", "june",
		}
		for _, nombre := range nombresPropios {
			diccionario[nombre] = Palabra{
				Tipo:     TipoSujeto,
				Texto:    nombre,
				Metadata: Metadata{EsNombrePropio: true},
			}
		}

		// Verbos regulares e irregulares en pasado
		agregarPalabras([]string{
			"played", "visited", "walked", "talked", "worked", "studied",
			"went", "saw", "ate", "drove", "wrote", "read", "slept",
			"bought", "sold", "taught", "caught", "thought", "brought",
		}, TipoVerboSimple)

		// Verbos auxiliares
		agregarPalabras([]string{
			"was", "were", "had", "did", "could", "would", "should",
			"might", "must", "shall", "will", "can", "may",
		}, TipoVerboAuxiliar)

		// Complementos (sustantivos comunes)
		agregarPalabras([]string{
			"football", "music", "movie", "book", "school", "house",
			"car", "dog", "cat", "computer", "phone", "food", "water",
			"time", "day", "night", "morning", "evening", "afternoon",
		}, TipoComplemento)

		// Expresiones de tiempo
		agregarPalabras([]string{
			"yesterday", "today", "tomorrow", "last week", "last year",
			"last month", "next week", "next year", "next month",
			"ago", "later", "soon", "now", "then",
		}, TipoTiempo)

		// Preposiciones
		agregarPalabras([]string{
			"in", "on", "at", "by", "for", "with", "to", "from",
			"under", "over", "between", "among", "through", "during",
		}, TipoPreposicion)

		// Artículos
		agregarPalabras([]string{
			"a", "an", "the",
		}, TipoArticulo)

		// Adjetivos comunes
		agregarPalabras([]string{
			"big", "small", "good", "bad", "happy", "sad", "new",
			"old", "young", "fast", "slow", "hot", "cold", "beautiful",
		}, TipoAdjetivo)

		// Adverbios
		agregarPalabras([]string{
			"quickly", "slowly", "carefully", "happily", "sadly",
			"very", "really", "quite", "almost", "always", "never",
		}, TipoAdverbio)

		// Conjunciones
		agregarPalabras([]string{
			"and", "but", "or", "nor", "for", "yet", "so",
			"because", "although", "unless", "since", "while",
		}, TipoConjuncion)
	})
}

// agregarPalabras agrega palabras al diccionario con thread safety
func agregarPalabras(palabras []string, tipo TipoPalabra) {
	mu.Lock()
	defer mu.Unlock()

	for _, palabra := range palabras {
		diccionario[palabra] = Palabra{
			Tipo:  tipo,
			Texto: palabra,
		}
	}
}

// preprocesarTexto limpia y prepara el texto para análisis
func preprocesarTexto(texto string) string {
	// Convertir a minúsculas manteniendo las mayúsculas en nombres propios
	palabras := strings.Fields(texto)
	for i, palabra := range palabras {
		if !esPosibleNombrePropio(palabra) {
			palabras[i] = strings.ToLower(palabra)
		}
	}

	return strings.Join(palabras, " ")
}

// esPosibleNombrePropio verifica si una palabra podría ser un nombre propio
func esPosibleNombrePropio(palabra string) bool {
	if len(palabra) == 0 {
		return false
	}
	return unicode.IsUpper(rune(palabra[0]))
}

// obtenerContextoPalabra obtiene el contexto de una palabra en la oración
func obtenerContextoPalabra(palabras []string, tokens []Token, posicion int) Contexto {
	contexto := Contexto{
		PosicionEnOracion: posicion,
	}

	if posicion > 0 {
		contexto.PalabraAnterior = palabras[posicion-1]
		if len(tokens) > 0 {
			contexto.TipoAnterior = tokens[posicion-1].Tipo
		}
	}

	if posicion < len(palabras)-1 {
		contexto.PalabraSiguiente = palabras[posicion+1]
	}

	return contexto
}

// ClasificarPalabra determina el tipo de una palabra con más contexto
func ClasificarPalabra(palabra string, ctx Contexto) Palabra {
	inicializarDiccionario()

	palabraOriginal := palabra
	palabra = strings.ToLower(strings.TrimSpace(palabra))

	mu.RLock()
	clasificacion, existe := diccionario[palabra]
	mu.RUnlock()

	if existe {
		clasificacion.Original = palabraOriginal
		clasificacion.Posicion = ctx.PosicionEnOracion
		return clasificacion
	}

	// Análisis morfológico básico para palabras desconocidas
	switch {
	case strings.HasSuffix(palabra, "ly"):
		return Palabra{TipoAdverbio, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{}}
	case strings.HasSuffix(palabra, "ed"):
		return Palabra{TipoVerboSimple, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{}}
	case strings.HasSuffix(palabra, "ing"):
		return Palabra{TipoVerboSimple, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{}}
	case esPosibleNombrePropio(palabraOriginal):
		return Palabra{TipoSujeto, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{EsNombrePropio: true}}
	}

	// Análisis basado en contexto
	if ctx.TipoAnterior == TipoArticulo {
		return Palabra{TipoComplemento, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{}}
	}

	return Palabra{TipoDesconocido, palabra, palabraOriginal, ctx.PosicionEnOracion, Metadata{}}
}

// AnalizarLexico recibe una oración y devuelve una lista de tokens con análisis mejorado
func AnalizarLexico(oracion string) ([]Token, error) {
	if strings.TrimSpace(oracion) == "" {
		return nil, &ErrorAnalisis{
			Mensaje:  "la oración está vacía",
			Posicion: 0,
			Contexto: "",
		}
	}

	// Preprocesar el texto
	oracion = preprocesarTexto(oracion)
	palabras := strings.Fields(oracion)
	tokens := make([]Token, 0, len(palabras))

	for i, palabra := range palabras {
		ctx := obtenerContextoPalabra(palabras, tokens, i)
		p := ClasificarPalabra(palabra, ctx)

		token := Token{
			Tipo:     p.Tipo,
			Texto:    p.Texto,
			Original: p.Original,
			Posicion: i,
			Metadata: p.Metadata,
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// ValidarTokens valida los tokens y la estructura de la oración con reglas mejoradas
func ValidarTokens(tokens []Token) (string, string) {
	if len(tokens) == 0 {
		return "Inválida", "No se encontraron tokens."
	}

	// Estructura para seguimiento de elementos
	type ElementoOracion struct {
		encontrado bool
		posicion   int
		cantidad   int
	}

	elementos := map[TipoPalabra]*ElementoOracion{
		TipoArticulo:      {false, -1, 0},
		TipoSujeto:        {false, -1, 0},
		TipoVerboAuxiliar: {false, -1, 0},
		TipoVerboSimple:   {false, -1, 0},
		TipoComplemento:   {false, -1, 0},
		TipoTiempo:        {false, -1, 0},
		TipoPreposicion:   {false, -1, 0},
	}

	// Analizar la estructura
	for i, token := range tokens {
		if elemento, existe := elementos[token.Tipo]; existe {
			if !elemento.encontrado {
				elemento.encontrado = true
				elemento.posicion = i
			}
			elemento.cantidad++
		}
	}

	// Validaciones de estructura
	if !elementos[TipoSujeto].encontrado {
		return "Inválida", "Falta el sujeto en la oración."
	}

	if !elementos[TipoVerboSimple].encontrado && !elementos[TipoVerboAuxiliar].encontrado {
		return "Inválida", "Falta el verbo en la oración."
	}

	// Validar orden de elementos
	if elementos[TipoSujeto].posicion > elementos[TipoVerboSimple].posicion {
		return "Inválida", "El sujeto debe preceder al verbo principal."
	}

	if elementos[TipoComplemento].encontrado &&
		elementos[TipoVerboSimple].posicion > elementos[TipoComplemento].posicion {
		return "Inválida", "El verbo debe preceder al complemento."
	}

	// Validar coherencia de tiempos verbales
	if elementos[TipoVerboAuxiliar].encontrado && elementos[TipoVerboSimple].encontrado {
		if elementos[TipoVerboAuxiliar].posicion > elementos[TipoVerboSimple].posicion {
			return "Inválida", "El verbo auxiliar debe preceder al verbo principal."
		}
	}

	// Validar uso de preposiciones
	if elementos[TipoPreposicion].encontrado {
		if elementos[TipoPreposicion].posicion == len(tokens)-1 {
			return "Inválida", "La oración no puede terminar en preposición."
		}
	}

	// Validar artículos
	if elementos[TipoArticulo].encontrado {
		siguientePosicion := elementos[TipoArticulo].posicion + 1
		if siguientePosicion >= len(tokens) ||
			(tokens[siguientePosicion].Tipo != TipoComplemento &&
				tokens[siguientePosicion].Tipo != TipoAdjetivo) {
			return "Inválida", "El artículo debe ser seguido por un sustantivo o adjetivo."
		}
	}

	return "Válida", "La oración tiene una estructura válida y coherente."
}

// ValidarOracion valida una oración con análisis mejorado
func ValidarOracion(oracion string) (string, string) {
	tokens, err := AnalizarLexico(oracion)
	if err != nil {
		return "Inválida", "Error en análisis léxico: " + err.Error()
	}
	return ValidarTokens(tokens)
}

// Funciones auxiliares
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
*/
