# Variables
APP_NAME=api-http-postgresql
APP_BIN=./bin/$(APP_NAME)
MAIN_PATH=./cmd/main.go

# Construir la aplicación
build:
	@echo "Construyendo la aplicación..."
	@mkdir -p bin
	@go build -o $(APP_BIN) $(MAIN_PATH)


dev:
	@echo "Iniciando servicios con Docker Compose..."
	@docker-compose up -d
	@echo "Servicios iniciados correctamente."
	@sleep 2
	@echo "Ejecutando la aplicación en modo desarrollo..."
	@bash -c "trap 'echo \"Deteniendo servicios...\"; docker-compose down' EXIT; go run $(MAIN_PATH)"



# Ejecutar la aplicación
run: build
	@echo "Iniciando servicios con Docker Compose..."
	@docker-compose up -d
	@echo "Servicios iniciados correctamente."
	@sleep 2
	@echo "Ejecutando la aplicación..."
	@bash -c "trap 'echo \"Deteniendo servicios...\"; docker-compose down' EXIT; $(APP_BIN)"