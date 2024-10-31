package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"validar_oraciones/validators"
)

// ResultadoOracion representa el resultado de validar una oración
type ResultadoOracion struct {
	Oracion     string `json:"oracion"`
	EsValida    bool   `json:"esValida"`
	Mensaje     string `json:"mensaje"`
	Explicacion string `json:"explicacion"`
}

// PageVariables contiene todas las variables necesarias para la plantilla
type PageVariables struct {
	Oraciones        []ResultadoOracion
	ErrorMessage     string
	TotalOraciones   int
	OracionesValidas int
	ShowResults      bool
	Estadisticas     Estadisticas
}

// Estadisticas contiene información estadística sobre las oraciones analizadas
type Estadisticas struct {
	PorcentajeExito float64
	ErroresComunes  map[string]int
	TiposValidos    map[string]int
}

// ValidadorConfig contiene la configuración para el validador
type ValidadorConfig struct {
	MaxOraciones   int
	MinPalabras    int
	MaxPalabras    int
	LimpiarEntrada bool
	ModoEstricto   bool
}

// NewValidadorConfig crea una nueva configuración con valores por defecto
func NewValidadorConfig() ValidadorConfig {
	return ValidadorConfig{
		MaxOraciones:   5,
		MinPalabras:    2,
		MaxPalabras:    15,
		LimpiarEntrada: true,
		ModoEstricto:   false,
	}
}

// OracionHandler maneja las solicitudes relacionadas con la validación de oraciones
type OracionHandler struct {
	config    ValidadorConfig
	templates *template.Template
	logger    *log.Logger
}

// NewOracionHandler crea una nueva instancia del manejador
func NewOracionHandler(config ValidadorConfig, logger *log.Logger) (*OracionHandler, error) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		return nil, err
	}

	return &OracionHandler{
		config:    config,
		templates: tmpl,
		logger:    logger,
	}, nil
}

// limpiarOracion elimina caracteres no deseados y espacios extra
func (h *OracionHandler) limpiarOracion(oracion string) string {
	// Eliminar caracteres especiales excepto puntuación básica
	oracion = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ' || r == '.' || r == ',' {
			return r
		}
		return -1
	}, oracion)

	// Normalizar espacios
	return strings.Join(strings.Fields(oracion), " ")
}

// validarLongitud verifica que la oración cumpla con los límites de palabras
func (h *OracionHandler) validarLongitud(oracion string) error {
	palabras := strings.Fields(oracion)
	if len(palabras) < h.config.MinPalabras {
		return fmt.Errorf("la oración debe tener al menos %d palabras", h.config.MinPalabras)
	}
	if len(palabras) > h.config.MaxPalabras {
		return fmt.Errorf("la oración no debe exceder %d palabras", h.config.MaxPalabras)
	}
	return nil
}

// ServeHTTP maneja las solicitudes HTTP
func (h *OracionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// handleGet maneja las solicitudes GET
func (h *OracionHandler) handleGet(w http.ResponseWriter, _ *http.Request) {
	vars := PageVariables{
		ShowResults: false,
	}
	h.renderTemplate(w, vars)
}

// handlePost maneja las solicitudes POST
func (h *OracionHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.handleError(w, "Error al procesar el formulario", err)
		return
	}

	input := r.FormValue("oraciones")
	oraciones := h.procesarEntrada(input)

	if len(oraciones) > h.config.MaxOraciones {
		vars := PageVariables{
			ErrorMessage: fmt.Sprintf("Por favor, ingrese un máximo de %d oraciones.", h.config.MaxOraciones),
		}
		h.renderTemplate(w, vars)
		return
	}

	resultados := h.validarOraciones(oraciones)
	stats := h.calcularEstadisticas(resultados)

	vars := PageVariables{
		Oraciones:        resultados,
		TotalOraciones:   len(resultados),
		OracionesValidas: stats.TiposValidos["válidas"],
		ShowResults:      true,
		Estadisticas:     stats,
	}

	h.renderTemplate(w, vars)
}

// procesarEntrada divide y limpia las oraciones de entrada
func (h *OracionHandler) procesarEntrada(input string) []string {
	oraciones := strings.Split(input, ".")
	var processed []string

	for _, o := range oraciones {
		if o = strings.TrimSpace(o); o != "" {
			if h.config.LimpiarEntrada {
				o = h.limpiarOracion(o)
			}
			if o != "" {
				processed = append(processed, o)
			}
		}
	}

	return processed
}

// validarOraciones procesa y valida cada oración
func (h *OracionHandler) validarOraciones(oraciones []string) []ResultadoOracion {
	var resultados []ResultadoOracion

	for _, oracion := range oraciones {
		if err := h.validarLongitud(oracion); err != nil {
			resultados = append(resultados, ResultadoOracion{
				Oracion:     oracion,
				EsValida:    false,
				Mensaje:     "Longitud inválida",
				Explicacion: err.Error(),
			})
			continue
		}

		validez, explicacion := validators.ValidarOracion(oracion)
		resultados = append(resultados, ResultadoOracion{
			Oracion:     oracion,
			EsValida:    validez == "Válida",
			Mensaje:     validez,
			Explicacion: explicacion,
		})
	}

	return resultados
}

// calcularEstadisticas genera estadísticas sobre los resultados
func (h *OracionHandler) calcularEstadisticas(resultados []ResultadoOracion) Estadisticas {
	stats := Estadisticas{
		ErroresComunes: make(map[string]int),
		TiposValidos:   make(map[string]int),
	}

	validas := 0
	for _, r := range resultados {
		if r.EsValida {
			validas++
			stats.TiposValidos["válidas"]++
		} else {
			stats.ErroresComunes[r.Mensaje]++
		}
	}

	if len(resultados) > 0 {
		stats.PorcentajeExito = float64(validas) / float64(len(resultados)) * 100
	}

	return stats
}

// renderTemplate renderiza la plantilla HTML con las variables proporcionadas
func (h *OracionHandler) renderTemplate(w http.ResponseWriter, vars PageVariables) {
	if err := h.templates.Execute(w, vars); err != nil {
		h.handleError(w, "Error al renderizar la plantilla", err)
	}
}

// handleError maneja y registra errores
func (h *OracionHandler) handleError(w http.ResponseWriter, message string, err error) {
	h.logger.Printf("Error: %v - %v", message, err)
	http.Error(w, message, http.StatusInternalServerError)
}

// HandleAPIValidation maneja las solicitudes de validación vía API
func (h *OracionHandler) HandleAPIValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Oraciones []string `json:"oraciones"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	resultados := h.validarOraciones(input.Oraciones)
	stats := h.calcularEstadisticas(resultados)

	response := struct {
		Resultados   []ResultadoOracion `json:"resultados"`
		Estadisticas Estadisticas       `json:"estadisticas"`
	}{
		Resultados:   resultados,
		Estadisticas: stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
