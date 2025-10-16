package main

import (
	"fmt"
	"net/http"
)

func initServer() {
	// Define el puerto y muestra un mensaje en consola
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	// Inicia el servidor HTTP
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err) // Muestra error si falla el inicio
	}
}
