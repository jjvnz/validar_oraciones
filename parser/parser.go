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
	Verbos struct {
		Regulares   []string `json:"regulares"`
		Irregulares struct {
			VerbosComunes []string `json:"verbos_comunes"`
			Auxiliares    []string `json:"verbos_auxiliares"`
		} `json:"irregulares"`
	} `json:"verbos"`
	Sujeto            []string            `json:"sujeto"`
	Complementos      Complementos        `json:"complementos"`
	Preposiciones     []string            `json:"preposiciones"`
	Articulos         []string            `json:"articulos"`
	Adjetivos         map[string][]string `json:"adjetivos"`
	Adverbios         map[string][]string `json:"adverbios"`
	ExpresionesTiempo []string            `json:"expresiones_tiempo"`
	ModalesPasados    []string            `json:"modales_pasados"` // Campo agregado para los verbos modales pasados
}

// Nuevo struct para modelar Complementos como un objeto en lugar de una lista
type Complementos struct {
	Objetos []string `json:"objetos"`
	Lugares []string `json:"lugares"`
	Comida  []string `json:"comida"`
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
		agregarPalabras(wordsData.Verbos.Regulares, models.TipoVerboSimple)
		agregarPalabras(wordsData.Verbos.Irregulares.VerbosComunes, models.TipoVerboSimple) // Verbos comunes
		agregarPalabras(wordsData.Verbos.Irregulares.Auxiliares, models.TipoVerboAuxiliar)  // Verbos auxiliares
		agregarPalabras(wordsData.Adjetivos["estado"], models.TipoVerboEstado)              // Si hay una categoría para Estado en la estructura
		agregarPalabras(wordsData.ModalesPasados, models.TipoVerboModalPasado)              // Nuevos verbos modales
		agregarPalabras(wordsData.ExpresionesTiempo, models.TipoTiempo)
		agregarPalabras(wordsData.Preposiciones, models.TipoPreposicion)
		agregarPalabras(wordsData.Articulos, models.TipoArticulo)
		agregarPalabras(wordsData.Adjetivos["apariencia"], models.TipoAdjetivo)
		agregarPalabras(wordsData.Adjetivos["personalidad"], models.TipoAdjetivo)
		agregarPalabras(wordsData.Adjetivos["estado"], models.TipoAdjetivo)
		agregarPalabras(wordsData.Adverbios["tiempo"], models.TipoAdverbio)
		agregarPalabras(wordsData.Adverbios["modo"], models.TipoAdverbio)
		agregarPalabras(wordsData.Adverbios["frecuencia"], models.TipoAdverbio)

		// Agregar complementos
		agregarPalabras(wordsData.Complementos.Objetos, models.TipoComplemento)
		agregarPalabras(wordsData.Complementos.Lugares, models.TipoComplemento)
		agregarPalabras(wordsData.Complementos.Comida, models.TipoComplemento)

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

// Función para agregar palabras al diccionario sin metadata
func agregarPalabras(palabras []string, tipo models.TipoPalabra) {
	agregarPalabrasConMetadata(palabras, tipo, models.Metadata{})
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

func ValidarTokens(tokens []models.Token) (string, string) {
	if len(tokens) == 0 {
		return "Inválida", "No se encontraron tokens."
	}

	// Inicializar elementos de la oración
	elementos := map[models.TipoPalabra]*models.ElementoOracion{
		models.TipoSujeto:           {Encontrado: false, Posicion: -1},
		models.TipoVerboSimple:      {Encontrado: false, Posicion: -1},
		models.TipoVerboEstado:      {Encontrado: false, Posicion: -1},
		models.TipoVerboModalPasado: {Encontrado: false, Posicion: -1},
		models.TipoComplemento:      {Encontrado: false, Posicion: -1},
		models.TipoNegativo:         {Encontrado: false, Posicion: -1},
	}

	// Lista de auxiliares no permitidos
	auxiliaresNoPermitidos := map[string]bool{
		"has":  true,
		"have": true,
		"had":  true,
		"do":   true,
		"does": true,
		"did":  true,
		"am":   true,
		"is":   true,
		"are":  true,
	}

	// Palabras negativas
	palabrasNegativas := map[string]bool{
		"not":   true,
		"never": true,
		"no":    true,
	}

	// Variables para rastrear detalles importantes
	primeraAparicionWasWere := -1
	primerSujeto := -1
	posicionesPermitidas := map[models.TipoPalabra]bool{
		models.TipoPreposicion: true,
		models.TipoComplemento: true,
		models.TipoArticulo:    true,
		models.TipoAdjetivo:    true,
	}

	// Recorrer tokens y actualizar elementos
	for i, token := range tokens {
		// Verificar palabras negativas
		if palabrasNegativas[token.Texto] {
			return "Inválida", "No se permiten construcciones negativas en oraciones afirmativas."
		}

		// Verificar auxiliares no permitidos
		if auxiliaresNoPermitidos[token.Texto] {
			return "Inválida", "No se permiten verbos auxiliares en pasado simple afirmativo."
		}

		// Actualizar la primera aparición de was/were
		if token.Texto == "was" || token.Texto == "were" {
			if primeraAparicionWasWere == -1 {
				primeraAparicionWasWere = i
			}
			elementos[models.TipoVerboSimple].Encontrado = true
			elementos[models.TipoVerboSimple].Posicion = i
			elementos[models.TipoVerboSimple].Cantidad++
		}

		// Actualizar el primer sujeto
		if token.Tipo == models.TipoSujeto {
			if primerSujeto == -1 {
				primerSujeto = i
			}
			elementos[models.TipoSujeto].Encontrado = true
			elementos[models.TipoSujeto].Posicion = i
			elementos[models.TipoSujeto].Cantidad++
		}

		// Actualizar otros elementos de la oración
		if elemento, existe := elementos[token.Tipo]; existe && !elemento.Encontrado {
			elemento.Encontrado = true
			elemento.Posicion = i
			elemento.Cantidad++
		}

		// Marcar complementos
		if token.Tipo == models.TipoComplemento {
			elementos[models.TipoComplemento].Encontrado = true
			elementos[models.TipoComplemento].Posicion = i
			elementos[models.TipoComplemento].Cantidad++
		}
	}

	// Validaciones estrictas
	// Verificar que haya un sujeto antes de was/were
	if primeraAparicionWasWere != -1 {
		if primerSujeto == -1 || primerSujeto >= primeraAparicionWasWere {
			return "Inválida", "Falta un sujeto antes del verbo 'was' o 'were'."
		}

		// Verificar que no haya tokens no permitidos entre sujeto y verbo
		for i := primerSujeto + 1; i < primeraAparicionWasWere; i++ {
			if !posicionesPermitidas[tokens[i].Tipo] {
				return "Inválida", "El verbo debe seguir inmediatamente al sujeto."
			}
		}
	}

	// Verificar que la oración tenga un sujeto
	if !elementos[models.TipoSujeto].Encontrado {
		return "Inválida", "Falta el sujeto en la oración."
	}

	// Verificar que haya al menos un verbo (incluyendo was/were)
	tieneVerboSimple := elementos[models.TipoVerboSimple].Encontrado
	tieneVerboEstado := elementos[models.TipoVerboEstado].Encontrado
	tieneVerboModalPasado := elementos[models.TipoVerboModalPasado].Encontrado

	if !tieneVerboSimple && !tieneVerboEstado && !tieneVerboModalPasado {
		return "Inválida", "Falta un verbo en pasado en la oración."
	}

	// Verificar que no haya construcciones negativas incorrectas
	if elementos[models.TipoNegativo].Encontrado && elementos[models.TipoVerboSimple].Encontrado {
		return "Inválida", "La oración no puede contener verbos modales y negativos en la misma estructura."
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
