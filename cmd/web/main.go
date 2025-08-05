package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 1. Creamos el "manejador de archivos" que apunta a nuestra carpeta estática.
	// El path "../../ui/static" es relativo a donde ejecutamos el main.go.
	// Desde 'cmd/web/', subimos dos niveles (..) para llegar a la raíz del proyecto
	// y luego entramos a 'ui/static/'.
	fileServer := http.FileServer(http.Dir("../../ui/static"))

	// 2. Le decimos a Go que use este manejador de archivos para todas las peticiones a la raíz "/".
	// Cuando alguien visite http://localhost:8080/, Go buscará un 'index.html' en esa carpeta.
	http.Handle("/", fileServer)

	fmt.Println("Servidor iniciando en el puerto 8080...")
	fmt.Println("Visita http://localhost:8080 en tu navegador.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
