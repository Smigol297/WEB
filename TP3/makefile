BINARY=myapp

# Compilar
build:
	go build -o $(BINARY) .

# Ejecutar servidor
run: build
	./$(BINARY) & echo $$! > server.pid
	sleep 1
# Ejecutar getProducts
get-products:
	curl http://localhost:8080/products
add-products:
	# Agregar productos de ejemplo
	curl -X POST http://localhost:8080/products \
	-H "Content-Type: application/json" \
	-d '{"name":"LaptopNueva","description":"Laptop 16GB RAM","price":3500,"quantity":10}'

	curl -X POST http://localhost:8080/products \
	-H "Content-Type: application/json" \
	-d '{"name":"MouseNuevo","description":"Mouse inalámbrico","price":50,"quantity":100}'

# Test endpoints con distintos IDs
test-ids:
	# GET por ID
	curl -X GET http://localhost:8080/products/1

	curl -X GET http://localhost:8080/products/2

	# PUT por ID (actualización de ejemplo)
	curl -X PUT http://localhost:8080/products/1 \
	-H "Content-Type: application/json" \
	-d '{"name":"Laptop Updated","description":"Laptop 16GB RAM","price":3600,"quantity":12}'
	

	curl -X PUT http://localhost:8080/products/2 \
	-H "Content-Type: application/json" \
	-d '{"name":"Mouse Updated","description":"Mouse inalámbrico","price":55,"quantity":90}'
	

	# DELETE por ID
	curl -X DELETE http://localhost:8080/products/3
	
	curl -X DELETE http://localhost:8080/products/4
	

# Ejecutar flujo completo
all: run get-products add-productstest-ids get-products stop clean

# Limpiar binarios
clean:
	rm -f $(BINARY)

stop:
	@if [ -f server.pid ]; then \
		kill $$(cat server.pid); \
		rm -f server.pid; \
	else \
		 "No se encontró server.pid, utilizar lsoft -i :8080 para buscar el proceso y luego kill 'PID' para matarlo"; \
	fi

reboot: stop run