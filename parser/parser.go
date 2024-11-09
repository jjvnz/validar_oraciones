package validators

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
	"validar_oraciones/models"
)

// Variables globales
var (
	diccionario map[string]models.Palabra
	once        sync.Once
	mu          sync.RWMutex
)

// Estructura para leer el JSON de palabras
type WordsData struct {
	Verbosos struct {
		Regulares   []string `json:"regulares"`
		Irregulares []string `json:"irregulares"`
		Auxiliares  []string `json:"auxiliares"`
		Estado      []string `json:"estado"`
	} `json:"verbos"`
	Sujeto        []string `json:"sujeto"`
	Complementos  []string `json:"complementos"`
	Preposiciones []string `json:"preposiciones"`
	Articulos     []string `json:"articulos"`
	Adjetivos     []string `json:"adjetivos"`
	Adverbios     []string `json:"adverbios"`
	Conjunciones  []string `json:"conjunciones"`
	Tiempos       []string `json:"tiempos"`
}

// Inicializa el diccionario de palabras, asegurándose de hacerlo solo una vez
func inicializarDiccionario() {
	once.Do(func() {
		diccionario = make(map[string]models.Palabra, 1000)

		// Cargar las palabras desde el archivo JSON
		wordsData, err := cargarPalabrasDesdeJSON("words.json")
		if err != nil {
			log.Fatal("Error cargando palabras desde JSON:", err)
			return
		}

		// Agregar las palabras de cada categoría
		agregarPalabras(wordsData.Sujeto, models.TipoSujeto)
		agregarPalabras(wordsData.Verbosos.Regulares, models.TipoVerboSimple)
		agregarPalabras(wordsData.Verbosos.Irregulares, models.TipoVerboSimple)
		agregarPalabras(wordsData.Verbosos.Auxiliares, models.TipoVerboAuxiliar)
		agregarPalabras(wordsData.Verbosos.Estado, models.TipoVerboEstado)
		agregarPalabras(wordsData.Complementos, models.TipoComplemento)
		agregarPalabras(wordsData.Tiempos, models.TipoTiempo)
		agregarPalabras(wordsData.Preposiciones, models.TipoPreposicion)
		agregarPalabras(wordsData.Articulos, models.TipoArticulo)
		agregarPalabras(wordsData.Adjetivos, models.TipoAdjetivo)
		agregarPalabras(wordsData.Adverbios, models.TipoAdverbio)
		agregarPalabras(wordsData.Conjunciones, models.TipoConjuncion)
	})
}

// Función para cargar palabras desde el archivo JSON
func cargarPalabrasDesdeJSON(filepath string) (WordsData, error) {
	var wordsData WordsData
	file, err := os.ReadFile(filepath)
	if err != nil {
		return wordsData, err
	}

	err = json.Unmarshal(file, &wordsData)
	if err != nil {
		return wordsData, err
	}

	return wordsData, nil
}

// Función para agregar palabras con metadata
func agregarPalabrasConMetadata(palabras []string, tipo models.TipoPalabra, metadata models.Metadata) {
	mu.Lock()
	defer mu.Unlock()

	for _, palabra := range palabras {
		diccionario[palabra] = models.Palabra{
			Tipo:     tipo,
			Texto:    palabra,
			Metadata: metadata,
		}
	}
}

// Función para agregar palabras al diccionario sin metadata
func agregarPalabras(palabras []string, tipo models.TipoPalabra) {
	agregarPalabrasConMetadata(palabras, tipo, models.Metadata{})
}

// Función de preprocesamiento del texto
func preprocesarTexto(texto string) string {
	palabras := strings.Fields(texto)
	for i, palabra := range palabras {
		if !esPosibleNombrePropio(palabra) {
			palabras[i] = strings.ToLower(palabra)
		}
	}
	return strings.Join(palabras, " ")
}

// Verificar si una palabra puede ser un nombre propio
func esPosibleNombrePropio(palabra string) bool {
	return len(palabra) > 0 && unicode.IsUpper(rune(palabra[0]))
}

// Obtención del contexto de la palabra en la oración
func obtenerContextoPalabra(palabras []string, tokens []models.Token, posicion int) models.Contexto {
	contexto := models.Contexto{PosicionEnOracion: posicion}

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

// Clasificar palabra basándonos en su contexto
func ClasificarPalabra(palabra string, ctx models.Contexto) models.Palabra {
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

	// Clasificación de palabras con sufijos
	switch {
	case strings.HasSuffix(palabra, "ly"):
		return models.Palabra{Tipo: models.TipoAdverbio, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion}
	case strings.HasSuffix(palabra, "ed"):
		return models.Palabra{Tipo: models.TipoVerboSimple, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion}
	case strings.HasSuffix(palabra, "ing"):
		return models.Palabra{Tipo: models.TipoVerboSimple, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion}
	case esPosibleNombrePropio(palabraOriginal):
		return models.Palabra{Tipo: models.TipoSujeto, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion, Metadata: models.Metadata{EsNombrePropio: true}}
	}

	// Si el tipo anterior fue un artículo, se clasifica como complemento
	if ctx.TipoAnterior == models.TipoArticulo {
		return models.Palabra{Tipo: models.TipoComplemento, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion}
	}

	// Si no se encuentra en el diccionario, se marca como desconocido
	return models.Palabra{Tipo: models.TipoDesconocido, Texto: palabra, Original: palabraOriginal, Posicion: ctx.PosicionEnOracion}
}

// Análisis léxico de una oración
func AnalizarLexico(oracion string) ([]models.Token, error) {
	if strings.TrimSpace(oracion) == "" {
		return nil, &models.ErrorAnalisis{
			Mensaje:  "la oración está vacía",
			Posicion: 0,
			Contexto: "",
		}
	}

	oracion = preprocesarTexto(oracion)
	palabras := strings.Fields(oracion)
	tokens := make([]models.Token, 0, len(palabras))

	for i, palabra := range palabras {
		ctx := obtenerContextoPalabra(palabras, tokens, i)
		p := ClasificarPalabra(palabra, ctx)

		token := models.Token{
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

// Validación de tokens para oraciones en pasado simple afirmativo
func ValidarTokens(tokens []models.Token) (string, string) {
	if len(tokens) == 0 {
		return "Inválida", "No se encontraron tokens."
	}

	// Inicializar elementos de la oración
	elementos := map[models.TipoPalabra]*models.ElementoOracion{
		models.TipoSujeto:      {Encontrado: false, Posicion: -1},
		models.TipoVerboSimple: {Encontrado: false, Posicion: -1},
		models.TipoVerboEstado: {Encontrado: false, Posicion: -1},
		models.TipoComplemento: {Encontrado: false, Posicion: -1},
	}

	// Recorrer tokens y actualizar elementos
	for i, token := range tokens {
		if elemento, existe := elementos[token.Tipo]; existe && !elemento.Encontrado {
			elemento.Encontrado = true
			elemento.Posicion = i
			elemento.Cantidad++
		}
	}

	// Verificar que la oración tenga un sujeto
	if !elementos[models.TipoSujeto].Encontrado {
		return "Inválida", "Falta el sujeto en la oración."
	}

	// Verificar que haya al menos un verbo
	tieneVerboSimple := elementos[models.TipoVerboSimple].Encontrado
	tieneVerboEstado := elementos[models.TipoVerboEstado].Encontrado

	if !tieneVerboSimple && !tieneVerboEstado {
		return "Inválida", "Falta un verbo en pasado en la oración."
	}

	// Verificar que no haya verbos auxiliares
	for _, token := range tokens {
		if token.Tipo == models.TipoVerboAuxiliar {
			return "Inválida", "La oración no debe contener verbos auxiliares."
		}
	}

	// Determinar la posición del verbo
	var posicionVerbo int
	if tieneVerboSimple {
		posicionVerbo = elementos[models.TipoVerboSimple].Posicion
	} else {
		posicionVerbo = elementos[models.TipoVerboEstado].Posicion
	}

	// Verificar que el verbo siga al sujeto
	if elementos[models.TipoSujeto].Posicion > posicionVerbo {
		return "Inválida", "El verbo debe seguir al sujeto."
	}

	// Verificar que el complemento siga al verbo, si existe
	if elementos[models.TipoComplemento].Encontrado && posicionVerbo > elementos[models.TipoComplemento].Posicion {
		return "Inválida", "El complemento debe ir después del verbo."
	}

	return "Válida", "La oración tiene una estructura válida en pasado simple afirmativo."
}

// Función para validar toda la oración
func ValidarOracion(oracion string) (string, string) {
	tokens, err := AnalizarLexico(oracion)
	if err != nil {
		return "Inválida", "Error en análisis léxico: " + err.Error()
	}
	return ValidarTokens(tokens)
}
