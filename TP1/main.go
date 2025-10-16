package main

import (
	"compress/gzip" // Importa el paquete para comprimir datos usando gzip
	"fmt"           // Importa funciones de formato de texto
	"net/http"      // Importa el paquete para crear servidores HTTP
	"strings"       // Importa utilidades para manipular strings
)

// Handler para la ruta /validacion
func validacion(w http.ResponseWriter, req *http.Request) {
	// Verifica que la ruta y el método sean correctos
	if req.URL.Path != "/validacion" || req.Method != http.MethodGet {
		http.NotFound(w, req) // Responde 404 si no coincide
		return
	}
	// Intenta parsear el formulario
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Error al parsear", http.StatusBadRequest) // Error si falla el parseo
		return
	}
	// Obtiene los valores del formulario
	nombre := req.FormValue("nombre")
	email := req.FormValue("email")
	mensaje := req.FormValue("mensaje")
	// Verifica si algún campo está vacío
	if nombre == "" || email == "" || mensaje == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // Define el tipo de contenido
		// Responde con HTML indicando que faltan datos
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Bienvenido</title></head> <body><h1>Los datos ESTAN vacios</h1>
            <a href="/">Volver</a></body></html>`)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8") // Define el tipo de contenido
		// Responde con HTML indicando que los datos están completos
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Bienvenido</title></head> <body><h1>Los datos NO estan vacios</h1>
            <a href="/">Volver</a><br></body></html>`)
	}
	infoHandler(w, req) // Llama al handler de info para mostrar detalles de la petición
}

// Estructura que implementa http.ResponseWriter y redirige la escritura al gzip.Writer
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

// Sobrescribe el método Write para comprimir los datos antes de enviarlos
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	// Escribe los datos comprimidos
	return w.Writer.Write(b)
}

// Middleware que comprime la respuesta si el cliente acepta gzip
func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Accept-Encoding") // Indica a los caches que la respuesta depende de Accept-Encoding
		// Si el cliente no acepta gzip, responde normalmente
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")               // Indica que la respuesta está comprimida
		gz := gzip.NewWriter(w)                                  // Crea un escritor gzip
		defer gz.Close()                                         // Asegura que se cierre al final
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w} // Crea el ResponseWriter personalizado
		h.ServeHTTP(gzw, r)                                      // Llama al handler original usando el writer comprimido
	})
}

// Handler que muestra información de la petición
func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Muestra el método HTTP
	fmt.Fprintf(w, "Método: %s\n", r.Method)

	// Muestra todas las cabeceras recibidas
	fmt.Fprintf(w, "\nCabeceras:\n")
	for k, v := range r.Header {
		fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ", "))
	}

	// Procesa parámetros GET y POST
	r.ParseForm()
	fmt.Fprintf(w, "\nParámetros GET/POST:\n")
	for k, v := range r.Form {
		fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ", "))
	}

	// Muestra la dirección IP del cliente
	ip := r.RemoteAddr
	// Si hay proxy, intenta obtener IP real
	if realIP := r.Header.Get("X-Forwarded-For"); realIP != "" {
		ip = realIP
	}
	fmt.Fprintf(w, "\nDirección IP del cliente: %s\n", ip)
}

func main() {
	/// Crea un FileServer para servir archivos estáticos desde ./static
	fs := http.FileServer(http.Dir("./static"))

	// Registra un handler para la raíz que sirve index.html o archivos estáticos
	http.Handle("/", gzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html") // Sirve index.html en la raíz
			return
		}
		fs.ServeHTTP(w, r) // Sirve archivos estáticos para otras rutas
	})))
	// Registra el handler para /validacion, también comprimido
	http.Handle("/validacion", gzipHandler(http.HandlerFunc(validacion)))
	http.Handle("/info", gzipHandler(http.HandlerFunc(infoHandler)))

	initServer()
}
