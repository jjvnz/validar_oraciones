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