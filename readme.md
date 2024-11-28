# Validación de Oraciones en Pasado Simple Afirmativo

## Descripción del Proyecto

Este proyecto implementa un sistema de validación gramatical desarrollado en Go, especializado en analizar oraciones en inglés en **pasado simple afirmativo**. La aplicación utiliza técnicas avanzadas de análisis sintáctico para examinar la estructura gramatical de las oraciones, verificando:

- Correcta conjugación de verbos
- Uso adecuado de pronombres
- Estructura sintáctica del pasado simple
- Reglas de concordancia gramatical

### Características Principales

- 🔍 Análisis sintáctico computacional
- 📝 Validación de oraciones en pasado simple
- 🧠 Clasificación léxica avanzada
- 🚦 Detección precisa de errores gramaticales

## Arquitectura Técnica

### Estructura del Proyecto

```
validar_oraciones
├─ .dockerignore
├─ .gitignore
├─ Dockerfile
├─ go.mod
├─ go.sum
├─ handlers
│  ├─ oracion.go
│  └─ oracion_test.go
├─ main.go
├─ middleware
│  └─ middleware.go
├─ models
│  ├─ config_test.go
│  ├─ models.go
│  └─ models_test.go
├─ package-lock.json
├─ package.json
├─ parser
│  ├─ json_charge_test.go
│  ├─ parser.go
│  └─ parser_test.go
├─ readme.md
├─ static
│  ├─ css
│  │  └─ tailwind.css
│  └─ js
│     └─ scripts.js
├─ tailwind.config.js
├─ templates
│  └─ index.html
├─ words.json
└─ words_test.json

```

### Componentes del Sistema

1. **Tokenización**: Divide la oración en unidades mínimas
2. **Clasificación Léxica**: Categoriza cada palabra
3. **Validación Sintáctica**: Verifica la estructura gramatical
4. **Reglas de Conjugación**: Valida el uso correcto de verbos

### Diagrama de Flujo de Validación de Oraciones

![Diagrama de Flujo de Validación de Oraciones](https://github.com/user-attachments/assets/ff448646-e242-4471-b0f2-84ecb2fe2e0c)



### Tecnologías Utilizadas

- Lenguaje: Go (Golang)
- Análisis: Procesamiento de lenguaje natural (NLP)
- Estructuras de Datos: Mapas, Slices
- Concurrencia: sync.Mutex, sync.Once

## Requisitos Previos

- Go go 1.23.2
- Docker (opcional)

## Instalación y Configuración

### Instalación Directa

```bash
# Clonar repositorio
git clone https://github.com/jjvnz/validar_oraciones.git
cd validar_oraciones

# Instalar dependencias
go mod download

# Ejecutar proyecto
go run main.go
```

### Instalación con Docker

```bash
# Construir imagen
docker build -t validar_oraciones .

# Ejecutar contenedor
docker run -d -p 8080:8080 validar_oraciones
```

## Ejemplos de Uso

### Oraciones Válidas

✅ "I visited my grandmother last weekend."
✅ "She was happy yesterday."
✅ "They were at the park."

### Oraciones Inválidas

❌ "I visit my grandmother last weekend."
❌ "She were happy yesterday."
❌ "They was at the park."

## Funcionalidades Detalladas

- Validación de conjugaciones verbales
- Verificación de pronombres
- Detección de estructuras incorrectas
- Retroalimentación descriptiva de errores

## Desafíos Técnicos Resueltos

- Manejo concurrente de diccionarios
- Clasificación contextual de palabras
- Implementación de reglas gramaticales complejas

## Contribuciones

Las contribuciones son bienvenidas. Por favor, lea las directrices de contribución antes de enviar un pull request.

## Licencia

Este proyecto está bajo la Licencia MIT.

## Contacto

- Repositorio: https://github.com/jjvnz/validar_oraciones
- Desarrollador:
  - Juan Jair Villalobos Núñez

## Próximos Pasos

- [ ] Soporte para más tiempos verbales
- [ ] Mejora del sistema de clasificación léxica
- [ ] Implementación de machine learning

### Inicialización del diccionario de palabras

El código incluye una función `inicializarDiccionario()` que se ejecuta una sola vez utilizando `sync.Once`:

```go
var once sync.Once

func inicializarDiccionario() {
    once.Do(func() {
        diccionario = make(map[string]models.Palabra, 1000)
        
        // Carga y procesamiento de palabras desde JSON
        wordsData, err := cargarPalabrasDesdeJSON("words.json")
        if err != nil {
            log.Fatal("Error loading words from JSON:", err)
            return
        }
        
        // Agrega palabras a diferentes categorías del diccionario
        agregarPalabras(wordsData.Sujeto, models.TipoSujeto)
        agregarPalabras(wordsData.Verbos.Regulares, models.TipoVerboSimple)
        agregarPalabras(wordsData.Verbos.Irregulares.VerbosComunes, models.TipoVerboSimple)
        // ...
    })
}
```

Esta función carga palabras desde un archivo JSON y las agrega a un diccionario compartido utilizando un mutex para garantizar la concurrencia.

### Función de clasificación de palabras

La función `ClasificarPalabra()` utiliza el diccionario compartido para clasificar palabras:

```go
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

    // Clasificación basada en sufijos y contexto
    // ...
}
```

Esta función utiliza un mutex de solo lectura (`mu.RLock()`) para acceder al diccionario compartido.

### Función de análisis léxico

La función `AnalizarLexico()` analiza una oración y crea tokens:

```go
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
```

Esta función divide la oración en palabras y las clasifica utilizando la función `ClasificarPalabra()`.

### Función de validación de tokens

`ValidarTokens()` valida los tokens obtenidos del análisis léxico:

```go
func ValidarTokens(tokens []models.Token) (string, string) {
    // ... (validaciones detalladas)

    // Verificaciones estrictas
    if primeraAparicionWasWere != -1 {
        // Verificar uso correcto de was/were
        // Verificar estructura de la frase
    }

    // Verificar que haya sujeto
    if !elementos[models.TipoSujeto].Encontrado {
        return "Invalid", "El sujeto falta en la oración."
    }

    // Verificar presencia de verbo
    tieneVerboSimple := elementos[models.TipoVerboSimple].Encontrado
    tieneVerboEstado := elementos[models.TipoVerboEstado].Encontrado
    tieneVerboModalPasado := elementos[models.TipoVerboModalPasado].Encontrado

    if !tieneVerboSimple && !tieneVerboEstado && !tieneVerboModalPasado {
        return "Invalid", "Falta un verbo en tiempo pasado en la oración."
    }

    // Verificar construcciones negativas incorrectas
    if elementos[models.TipoNegativo].Encontrado && elementos[models.TipoVerboSimple].Encontrado {
        return "Invalid", "La oración no puede contener ambas construcciones modales y negativas en la misma estructura."
    }

    return "Valid", "La oración tiene una estructura válida en presente simple afirmativo."
}
```

Esta función realiza varias validaciones estrictas sobre la estructura de la oración, incluyendo el uso correcto de verbos regulares e irregulares, el papel del sujeto y otros elementos. Estos fragmentos muestran cómo se manejan recursos compartidos, se implementa el análisis léxico y se realizan validaciones gramaticales complejas en Go utilizando concurrencia y sincronización con mutexes.
