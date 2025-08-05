package main

import (
	"fmt"
	"log"
	"net/http"
)

// manejar las rutas de las peticiones a la ruta /
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "servidor en funcionamiento!")
}

func main() {
	// Le decimos a Go que la función "handler" se encargará de las peticiones a la raíz del sitio
	http.HandleFunc("/", handler)
	// Iniciamos el servidor en el puerto 8080
	fmt.Println("Servidor iniciando en el puerto 8080...")
	fmt.Println("Visita http://localhost:8080 en tu navegador.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
