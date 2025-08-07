package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// CedulaRequest representa la estructura de la petición de consulta por cédula
type CedulaRequest struct {
	Cedula string `json:"cedula"`
}

// NombresRequest representa la estructura de la petición de consulta por nombres
type NombresRequest struct {
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
}

// CedulaResponse representa la respuesta exitosa con los datos
type CedulaResponse struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// NombresResponse representa la respuesta exitosa con la cédula encontrada
type NombresResponse struct {
	Cedula    string `json:"cedula"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
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

// consultarCedula realiza la consulta a la API del SRI para obtener los datos de la cédula
func consultarCedula(cedula string) (*CedulaResponse, error) {
	// Construir la URL de la API del SRI
	timestamp := time.Now().UnixMilli()
	url := fmt.Sprintf("https://srienlinea.sri.gob.ec/movil-servicios/api/v1.0/deudas/porIdentificacion/%s/?tipoPersona=N&_=%d", cedula, timestamp)

	log.Printf("Consultando API del SRI: %s", url)

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Crear petición HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	// Configurar headers para simular un navegador real
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://srienlinea.sri.gob.ec/")

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar la petición: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %v", err)
	}

	log.Printf("Respuesta de la API (primeros 500 caracteres): %s", string(body)[:min(500, len(body))])

	// Verificar el código de estado HTTP
	if resp.StatusCode != 200 {
		log.Printf("Código de estado HTTP: %d", resp.StatusCode)
		return nil, fmt.Errorf("cédula no encontrada")
	}

	// Estructura para parsear la respuesta JSON del SRI
	type SRIResponse struct {
		Contribuyente struct {
			Identificacion  string `json:"identificacion"`
			Denominacion    string `json:"denominacion"`
			NombreComercial string `json:"nombreComercial"`
			Clase           string `json:"clase"`
		} `json:"contribuyente"`
	}

	var sriData SRIResponse
	if err := json.Unmarshal(body, &sriData); err != nil {
		log.Printf("Error al parsear JSON: %v", err)
		log.Printf("Respuesta completa: %s", string(body))
		return nil, fmt.Errorf("error al procesar la respuesta del servidor")
	}

	// Verificar que se encontraron datos
	nombreCompleto := ""
	if sriData.Contribuyente.Denominacion != "" {
		nombreCompleto = sriData.Contribuyente.Denominacion
	} else if sriData.Contribuyente.NombreComercial != "" {
		nombreCompleto = sriData.Contribuyente.NombreComercial
	}

	if nombreCompleto == "" {
		log.Printf("No se encontró información del nombre en la respuesta")
		return nil, fmt.Errorf("cédula no encontrada")
	}

	log.Printf("Datos encontrados - Identificación: %s, Nombre: %s, Clase: %s",
		sriData.Contribuyente.Identificacion, nombreCompleto, sriData.Contribuyente.Clase)

	// Procesar el nombre completo para separar nombre y apellido
	nombreCompleto = strings.TrimSpace(nombreCompleto)
	palabras := strings.Fields(nombreCompleto)

	var nombre, apellido string

	if len(palabras) >= 2 {
		// Asumimos que las primeras palabras son nombres y las últimas apellidos
		// Para nombres ecuatorianos, generalmente: PRIMER_NOMBRE SEGUNDO_NOMBRE PRIMER_APELLIDO SEGUNDO_APELLIDO
		if len(palabras) == 2 {
			nombre = palabras[0]
			apellido = palabras[1]
		} else if len(palabras) == 3 {
			nombre = palabras[0]
			apellido = strings.Join(palabras[1:], " ")
		} else {
			// 4 o más palabras
			mitad := len(palabras) / 2
			nombre = strings.Join(palabras[:mitad], " ")
			apellido = strings.Join(palabras[mitad:], " ")
		}
	} else {
		nombre = nombreCompleto
		apellido = ""
	}

	return &CedulaResponse{
		Nombre:   nombre,
		Apellido: apellido,
	}, nil
}

// Función auxiliar para min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// consultarPorNombres informa sobre las alternativas legales disponibles para búsqueda por nombres
func consultarPorNombres(nombres, apellidos string) (*NombresResponse, error) {
	log.Printf("Consulta por nombres solicitada: %s %s", nombres, apellidos)

	// En lugar de intentar scraping no autorizado, informamos sobre las alternativas legales
	log.Printf("INFORMACIÓN: Existen alternativas legales oficiales para consultas por nombres en Ecuador")

	// Simular un tiempo de procesamiento mientras "evaluamos" las opciones
	time.Sleep(2 * time.Second)

	// Retornar error con información educativa sobre las alternativas legales
	return nil, fmt.Errorf(`consulta por nombres no disponible a través de APIs públicas gratuitas.

ALTERNATIVAS LEGALES DISPONIBLES:

🏛️ FUNCIÓN JUDICIAL (SATJE)
• Consulta de procesos judiciales por nombre
• URL: https://procesosjudiciales.funcionjudicial.gob.ec/busqueda
• Permite buscar si una persona tiene procesos judiciales registrados

🗳️ CONSEJO NACIONAL ELECTORAL (CNE) 
• Consulta de personas registradas para votar
• Búsqueda por nombre y apellido
• Solo para ciudadanos habilitados para elecciones

🏥 IESS (Instituto Ecuatoriano de Seguridad Social)
• Consulta de afiliados (protegida con captcha)
• No tiene API pública abierta
• URL: https://www.iess.gob.ec/

💰 SERVICIOS DE PAGO DISPONIBLES:
• EcuadorLegalOnline: Consulta por nombres y apellidos
• URL: https://tramites.ecuadorlegalonline.com/
• Incluye datos completos: cédula, estado civil, profesión, etc.
• Servicio de pago con garantía

RECOMENDACIÓN: Use el servicio de consulta por cédula que funciona con datos oficiales del SRI (gratuito y confiable)`)
} // manejarConsulta maneja las peticiones POST al endpoint /api/consultar
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

// manejarConsultaPorNombres maneja las peticiones POST al endpoint /api/consultar-nombres
func manejarConsultaPorNombres(w http.ResponseWriter, r *http.Request) {
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
	var req NombresRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "JSON inválido"})
		return
	}

	// Validar que se proporcionen nombres y apellidos
	if strings.TrimSpace(req.Nombres) == "" || strings.TrimSpace(req.Apellidos) == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Se requieren nombres y apellidos"})
		return
	}

	// Realizar la "consulta" (que en realidad retorna información sobre alternativas legales)
	resultado, err := consultarPorNombres(req.Nombres, req.Apellidos)
	if err != nil {
		// En lugar de retornar error, enviamos una respuesta informativa
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":           false,
			"nombres":           req.Nombres,
			"apellidos":         req.Apellidos,
			"message":           "Consulta por nombres no disponible a través de APIs públicas gratuitas",
			"alternatives_info": true,
			"error_details":     err.Error(),
		})
		return
	}

	// Responder con los datos encontrados (si alguna vez funcionara)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultado)
}

func main() {
	// Configurar el servidor de archivos estáticos
	fs := http.FileServer(http.Dir("./ui/static/"))
	http.Handle("/", fs)

	// Configurar los endpoints de la API
	http.HandleFunc("/api/consultar", manejarConsulta)
	http.HandleFunc("/api/consultar-nombres", manejarConsultaPorNombres)

	// Configurar el puerto
	puerto := ":8085"

	fmt.Printf("🚀 Servidor iniciado en http://localhost%s\n", puerto)
	fmt.Println("📁 Sirviendo archivos estáticos desde ./ui/static/")
	fmt.Println("🔍 Endpoint de consulta por cédula disponible en /api/consultar")
	fmt.Println("👤 Endpoint de consulta por nombres disponible en /api/consultar-nombres")

	// Iniciar el servidor
	if err := http.ListenAndServe(puerto, nil); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
