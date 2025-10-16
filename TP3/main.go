package main

import (
	TP3 "TP3/logic"
	"net/http"
)

func main() {
	// Ejercicio 1
	/*product := TP3.Product{
		Id:          1,
		Name:        "Sample Product",
		Description: "This is a sample product",
		Price:       100.0,
		Quantity:    10,
	}

	err := TP3.ValidateProduct(product)
	if err != nil {
		fmt.Println("Error:", err)
		return
	} else {
		fmt.Println("Product is valid")
	}
	productPriceAfter := product.Price
	product = TP3.ApplyDiscount(product, 10)
	fmt.Printf("Product before discount: %+v\n", productPriceAfter)
	fmt.Printf("Product after discount: %+v\n", product.Price)*/

	// Ejercicio 2
	// curl -X GET http://localhost:8080/products para obtener productos
	// curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"Id":6,"Name":"Tablet","Description":"Tablet 10 pulgadas","Price":300.00,"Quantity":50}' para agregar un producto
	http.HandleFunc("/products", TP3.ProductsHandler) // Mapea la ruta /products a la función ProductsHandler del paquete TP3
	//Ejercicio 3
	http.HandleFunc("/products/{id}", TP3.ProductByIDHandler)

	initServer()
}
