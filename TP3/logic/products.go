package TP3

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// INICIO CON MINUSCULA = NO PUBLICO
// INICIO CON MAYUSCULA = PUBLICO
type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
	Quantity    int
}

var productSlice = []Product{
	{1, "Laptop", "Laptop gamer 16GB RAM", 3500.00, 10},
	{2, "Mouse", "Mouse inalámbrico", 50.00, 100},
	{3, "Teclado", "Teclado mecánico retroiluminado", 120.00, 75},
	{4, "Monitor", "Monitor 27 pulgadas 144Hz", 900.00, 20},
	{5, "Impresora", "Impresora multifuncional", 650.00, 15},
}

func ValidateProduct(p Product) error {
	if p.Name == "" {
		return fmt.Errorf("El nombre del producto no puede estar vacío")
	}
	if p.Price < 0 {
		return fmt.Errorf("El precio del producto no puede ser negativo")
	}
	return nil
}
func ApplyDiscount(p Product, percentage float64) Product {
	if percentage < 0 {
		fmt.Println("Invalid discount percentage")
		return p // No discount applied
	}
	p.Price -= p.Price * (percentage / 100)
	return p
}

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		addProduct(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //Define el encabezado HTTP Content-Type de la respuesta. Indica que el contenido de la respuesta es de tipo JSON.
	json.NewEncoder(w).Encode(productSlice)            //Codifica el slice productSlice en formato JSON y lo escribe en el ResponseWriter w, que envía la respuesta HTTP al cliente.
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	var p Product

	// decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// validar producto
	if err := ValidateProduct(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// asignar ID y agregar al slice
	p.Id = len(productSlice) + 1
	productSlice = append(productSlice, p)

	// devolver producto creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func ProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer el ID del producto de la URL
	var id int
	_, err := fmt.Sscanf(r.URL.Path, "/products/%d", &id)
	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProductByID(w, r, id)
	case http.MethodPut:
		putProductByID(w, r, id)
	case http.MethodDelete:
		deleteProductByID(w, r, id)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// retorna índice y puntero al producto, o nil si no existe
func findProductByID(id int) (int, *Product) {
	for i := range productSlice {
		if productSlice[i].Id == id {
			return i, &productSlice[i]
		}
	}
	return -1, nil
}

func getProductByID(w http.ResponseWriter, r *http.Request, id int) {
	if _, p := findProductByID(id); p != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
		return
	}
	http.Error(w, "Producto no encontrado", http.StatusNotFound)
}

func putProductByID(w http.ResponseWriter, r *http.Request, id int) {
	var updatedProduct Product

	// Decodificar JSON
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validar producto
	if err := ValidateProduct(updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Buscar y actualizar el producto
	if i, p := findProductByID(id); p != nil {
		updatedProduct.Id = id // Asegurar que el ID no cambie
		productSlice[i] = updatedProduct
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedProduct)
		return
	}
	http.Error(w, "Producto no encontrado", http.StatusNotFound)
}

func deleteProductByID(w http.ResponseWriter, r *http.Request, id int) {
	if i, p := findProductByID(id); p != nil {
		productSlice = append(productSlice[:i], productSlice[i+1:]...) // Eliminar el producto del slice
		w.WriteHeader(http.StatusNoContent)                            // Responder con 204 No Content
		return
	}
	http.Error(w, "Producto no encontrado", http.StatusNotFound)
}
