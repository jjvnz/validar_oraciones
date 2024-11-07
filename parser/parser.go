package validators

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

// ElementoOracion representa el estado de un elemento dentro de una oración
type ElementoOracion struct {
	encontrado bool // Indica si el elemento ha sido encontrado
	posicion   int  // Posición del elemento en la oración
	cantidad   int  // Cuántas veces aparece el elemento (aunque usualmente será 1)
}

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

		// Verbos en pasado simple
		agregarPalabras([]string{
			"played", "visited", "walked", "talked", "worked", "studied",
			"went", "saw", "ate", "drove", "wrote", "slept",
			"bought", "sold", "taught", "caught", "thought", "brought",
		}, TipoVerboSimple)

		// Verbos auxiliares (No deben estar presentes en pasados simples afirmativos)
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

// ValidarTokens valida los tokens y la estructura de la oración para pasado simple afirmativo
func ValidarTokens(tokens []Token) (string, string) {
	if len(tokens) == 0 {
		return "Inválida", "No se encontraron tokens."
	}

	// Estructura para seguimiento de elementos clave
	elementos := map[TipoPalabra]*ElementoOracion{
		TipoSujeto:      {false, -1, 0},
		TipoVerboSimple: {false, -1, 0},
		TipoComplemento: {false, -1, 0},
	}

	// Analizar la estructura de la oración
	for i, token := range tokens {
		if elemento, existe := elementos[token.Tipo]; existe {
			if !elemento.encontrado {
				elemento.encontrado = true
				elemento.posicion = i
			}
			elemento.cantidad++
		}
	}

	// Validaciones específicas para pasado simple afirmativo
	// Verificar que el sujeto esté presente
	if !elementos[TipoSujeto].encontrado {
		return "Inválida", "Falta el sujeto en la oración."
	}

	// Verificar que el verbo en pasado esté presente
	if !elementos[TipoVerboSimple].encontrado {
		return "Inválida", "Falta el verbo en pasado simple en la oración."
	}

	// Verificar que no haya verbos auxiliares (como "did", "was", "were")
	for _, token := range tokens {
		if token.Tipo == TipoVerboAuxiliar {
			return "Inválida", "La oración no debe contener verbos auxiliares."
		}
	}

	// Validar la posición del sujeto y el verbo (el verbo debe seguir al sujeto)
	if elementos[TipoSujeto].posicion > elementos[TipoVerboSimple].posicion {
		return "Inválida", "El verbo debe seguir al sujeto."
	}

	// Validar que los complementos (si existen) no precedan al verbo
	if elementos[TipoComplemento].encontrado &&
		elementos[TipoVerboSimple].posicion > elementos[TipoComplemento].posicion {
		return "Inválida", "El verbo debe preceder al complemento."
	}

	// Si se pasa todas las validaciones, la oración es válida
	return "Válida", "La oración tiene una estructura válida en pasado simple afirmativo."
}

// ValidarOracion valida una oración en pasado simple afirmativo
func ValidarOracion(oracion string) (string, string) {
	tokens, err := AnalizarLexico(oracion)
	if err != nil {
		return "Inválida", "Error en análisis léxico: " + err.Error()
	}
	return ValidarTokens(tokens)
}
