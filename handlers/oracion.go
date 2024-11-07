package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"validar_oraciones/models"
	parser "validar_oraciones/parser"
)

// OracionHandler maneja las solicitudes relacionadas con la validación de oraciones
type OracionHandler struct {
	config    models.ValidadorConfig
	templates *template.Template
	logger    *log.Logger
}

// NewOracionHandler crea una nueva instancia del manejador
func NewOracionHandler(config models.ValidadorConfig, logger *log.Logger) (*OracionHandler, error) {
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
	vars := models.PageVariables{
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
		vars := models.PageVariables{
			ErrorMessage: fmt.Sprintf("Por favor, ingrese un máximo de %d oraciones.", h.config.MaxOraciones),
		}
		h.renderTemplate(w, vars)
		return
	}

	resultados := h.validarOraciones(oraciones)
	stats := h.calcularEstadisticas(resultados)

	vars := models.PageVariables{
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

// validarOraciones procesa y valida cada oración usando el análisis léxico
func (h *OracionHandler) validarOraciones(oraciones []string) []models.ResultadoOracion {
	var resultados []models.ResultadoOracion

	for _, oracion := range oraciones {
		// Limpiar la oración antes de validarla
		oracion = h.limpiarOracion(oracion)

		if err := h.validarLongitud(oracion); err != nil {
			resultados = append(resultados, models.ResultadoOracion{
				Oracion:     oracion,
				EsValida:    false,
				Mensaje:     "Longitud inválida",
				Explicacion: err.Error(),
			})
			continue
		}

		// Análisis léxico
		tokens, err := parser.AnalizarLexico(oracion)
		if err != nil {
			resultados = append(resultados, models.ResultadoOracion{
				Oracion:     oracion,
				EsValida:    false,
				Mensaje:     "Error en el análisis léxico",
				Explicacion: err.Error(),
			})
			continue
		}

		// Validar la estructura de la oración basada en los tokens
		validez, explicacion := parser.ValidarTokens(tokens)
		resultados = append(resultados, models.ResultadoOracion{
			Oracion:     oracion,
			EsValida:    validez == "Válida",
			Mensaje:     validez,
			Explicacion: explicacion,
		})
	}

	return resultados
}

// calcularEstadisticas genera estadísticas sobre los resultados
func (h *OracionHandler) calcularEstadisticas(resultados []models.ResultadoOracion) models.Estadisticas {
	stats := models.Estadisticas{
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
		stats.PorcentajeExito = (float64(validas) / float64(len(resultados))) * 100
	}

	return stats
}

// renderTemplate renderiza la plantilla con las variables dadas
func (h *OracionHandler) renderTemplate(w http.ResponseWriter, vars models.PageVariables) {
	if err := h.templates.Execute(w, vars); err != nil {
		h.logger.Println("Error al renderizar la plantilla:", err)
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
	}
}

// handleError maneja y registra los errores
func (h *OracionHandler) handleError(w http.ResponseWriter, message string, err error) {
	h.logger.Println(message, err)
	http.Error(w, message, http.StatusBadRequest)
}

// HandleAPIValidation maneja la validación de oraciones a través de la API
func (h *OracionHandler) HandleAPIValidation(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Oracion string `json:"oracion"`
	}

	// Decodificar el cuerpo de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Análisis léxico
	tokens, err := parser.AnalizarLexico(request.Oracion)
	if err != nil {
		h.logger.Printf("Error en el análisis léxico: %v", err)
		http.Error(w, "Error en el análisis de la oración", http.StatusInternalServerError)
		return
	}

	// Validar la estructura de la oración basada en los tokens
	validez, explicacion := parser.ValidarTokens(tokens)

	response := struct {
		Tokens      []parser.Token `json:"tokens"`
		EsValida    bool           `json:"es_valida"`
		Mensaje     string         `json:"mensaje"`
		Explicacion string         `json:"explicacion"`
	}{
		Tokens:      tokens,
		EsValida:    validez == "Válida",
		Mensaje:     validez,
		Explicacion: explicacion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
