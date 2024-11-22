package validators

import (
	"encoding/json"
	"os"
	"testing"
)

type Palabras struct {
	Palabras []string `json:"palabras"`
}

// Función para cargar palabras desde JSON (asumiendo que ya tienes esta función definida)
func CargarPalabrasDesdeJSON(filename string) ([]string, error) {
	// Abre el archivo JSON
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decodifica el JSON
	var palabras Palabras
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&palabras)
	if err != nil {
		return nil, err
	}

	return palabras.Palabras, nil
}

func TestCargarPalabrasDesdeJSON(t *testing.T) {
	// Intentamos cargar el archivo de prueba
	wordsData, err := CargarPalabrasDesdeJSON("../words_test.json")
	if err != nil {
		t.Fatalf("Error loading words from JSON: %v", err)
	}

	// Verificamos que el archivo se cargó correctamente y contiene las palabras esperadas
	expectedWords := []string{"jugador", "fútbol", "jugar", "partido"}

	// Comprobamos que las palabras cargadas coinciden con las esperadas
	for i, word := range expectedWords {
		if wordsData[i] != word {
			t.Errorf("Se esperaba '%v' pero se obtuvo '%v'", word, wordsData[i])
		}
	}
}
