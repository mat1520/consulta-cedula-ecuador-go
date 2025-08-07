package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

// CedulaRequest representa la estructura de la petición de consulta
type CedulaRequest struct {
	Cedula string `json:"cedula"`
}

// CedulaResponse representa la respuesta exitosa con los datos
type CedulaResponse struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// ErrorResponse representa la respuesta de error
type ErrorResponse struct {
	Error string `json:"error"`
}

// validarCedula valida que la cédula sea un número de 10 dígitos
func validarCedula(cedula string) bool {
	// Verificar que tenga exactamente 10 dígitos
	if len(cedula) != 10 {
		return false
	}

	// Verificar que todos los caracteres sean números
	match, _ := regexp.MatchString("^[0-9]+$", cedula)
	return match
}

// consultarCedula realiza el web scraping para obtener los datos de la cédula
func consultarCedula(cedula string) (*CedulaResponse, error) {
	// Crear nuevo colector con configuración para debugging
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Configurar User-Agent para evitar bloqueos
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var nombre, apellido string
	var encontrado bool

	// Configurar el selector para capturar los resultados
	c.OnHTML("table", func(e *colly.HTMLElement) {
		// Buscar en todas las filas de la tabla
		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			// Obtener todas las celdas de la fila
			cells := row.ChildTexts("td")

			// Si la fila tiene al menos 2 celdas y contiene datos relevantes
			if len(cells) >= 2 {
				// Buscar patrones que indiquen nombre y apellido
				for j, cell := range cells {
					cellText := strings.TrimSpace(cell)
					cellTextUpper := strings.ToUpper(cellText)

					// Si encontramos un texto que parece un nombre (contiene letras y espacios)
					if len(cellText) > 2 && regexp.MustCompile(`^[A-ZÁÉÍÓÚÑ\s]+$`).MatchString(cellTextUpper) {
						// Si es el primer nombre encontrado, considerarlo como nombre completo
						if !encontrado {
							// Dividir en palabras para separar nombre y apellido
							palabras := strings.Fields(cellTextUpper)
							if len(palabras) >= 2 {
								// Asumimos que las primeras palabras son nombres y las últimas apellidos
								mitad := len(palabras) / 2
								nombre = strings.Join(palabras[:mitad], " ")
								apellido = strings.Join(palabras[mitad:], " ")
								encontrado = true
							} else if len(palabras) == 1 {
								nombre = palabras[0]
								// Buscar en la siguiente celda para el apellido
								if j+1 < len(cells) {
									nextCell := strings.TrimSpace(strings.ToUpper(cells[j+1]))
									if regexp.MustCompile(`^[A-ZÁÉÍÓÚÑ\s]+$`).MatchString(nextCell) {
										apellido = nextCell
									}
								}
								encontrado = true
							}
						}
					}
				}
			}
		})
	})

	// Configurar manejo de errores HTTP
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error durante el scraping: %s", err.Error())
	})

	// Configurar el POST a la página de consulta
	err := c.Post("https://www.ecuadorlegalonline.com/modulo/consultar-cedula.php", map[string]string{
		"tipo": "cedula",
		"term": cedula,
	})

	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición: %v", err)
	}

	// Si no se encontraron datos, retornar error
	if !encontrado || nombre == "" {
		return nil, fmt.Errorf("cédula no encontrada")
	}

	return &CedulaResponse{
		Nombre:   nombre,
		Apellido: apellido,
	}, nil
}

// manejarConsulta maneja las peticiones POST al endpoint /api/consultar
func manejarConsulta(w http.ResponseWriter, r *http.Request) {
	// Configurar headers CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Manejar preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Verificar que sea una petición POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Método no permitido"})
		return
	}

	// Decodificar el JSON de la petición
	var req CedulaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "JSON inválido"})
		return
	}

	// Validar la cédula
	if !validarCedula(req.Cedula) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Cédula inválida. Debe contener exactamente 10 dígitos"})
		return
	}

	// Realizar la consulta
	resultado, err := consultarCedula(req.Cedula)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Cédula no encontrada"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Error interno del servidor al consultar"})
		}
		return
	}

	// Responder con los datos encontrados
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

func main() {
	// Configurar el servidor de archivos estáticos
	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/", fs)

	// Configurar el endpoint de la API
	http.HandleFunc("/api/consultar", manejarConsulta)

	// Configurar el puerto
	puerto := ":8081"

	fmt.Printf("🚀 Servidor iniciado en http://localhost%s\n", puerto)
	fmt.Println("📁 Sirviendo archivos estáticos desde ./ui/static/")
	fmt.Println("🔍 Endpoint de consulta disponible en /api/consultar")

	// Iniciar el servidor
	if err := http.ListenAndServe(puerto, nil); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
