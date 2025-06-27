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

# Ejecutar con SSL
run-ssl: build ssl-certs
	@echo "Iniciando servicios con Docker Compose..."
	@docker-compose up -d
	@echo "Servicios iniciados correctamente."
	@sleep 2
	@echo "Ejecutando la aplicación con SSL..."
	@USE_SSL=true SSL_CERT_PATH=certs/ssl/server.crt SSL_KEY_PATH=certs/ssl/server.key bash -c "trap 'echo \"Deteniendo servicios...\"; docker-compose down' EXIT; $(APP_BIN)"

# Desarrollo con SSL
dev-ssl: ssl-certs
	@echo "Iniciando servicios con Docker Compose..."
	@docker-compose up -d
	@echo "Servicios iniciados correctamente."
	@sleep 2
	@echo "Ejecutando la aplicación en modo desarrollo con SSL..."
	@USE_SSL=true SSL_CERT_PATH=certs/ssl/server.crt SSL_KEY_PATH=certs/ssl/server.key bash -c "trap 'echo \"Deteniendo servicios...\"; docker-compose down' EXIT; go run $(MAIN_PATH)"

mqtt-setup:
	@echo "=== CONFIGURACIÓN DE MOSQUITTO ==="
	@echo "Configurando MQTT con seguridad básica..."
	@mkdir -p config/mqtt
	@echo "Creando archivos de configuración MQTT..."

mqtt-certs:
	@echo "Generando certificados para MQTT..."
	@mkdir -p certs/mqtt
	@openssl req -new -x509 -days 365 -extensions v3_ca -keyout certs/mqtt/ca.key -out certs/mqtt/ca.crt -subj "/C=ES/ST=Madrid/L=Madrid/O=TestOrg/CN=localhost"
	
mqtt-status:
	@echo "=== ESTADO DE MOSQUITTO ==="
	@ps aux | grep mosquitto | grep -v grep || echo "Mosquitto no está corriendo"
	@echo "=== PUERTOS MQTT ==="
	@netstat -tlnp 2>/dev/null | grep 1883 || echo "Puerto 1883 no encontrado"

mqtt-auth-setup:
	@echo "=== CONFIGURANDO AUTENTICACIÓN MQTT ==="
	@echo "Creando usuario 'mqttuser' con contraseña..."
	@sudo mosquitto_passwd -c /etc/mosquitto/passwd mqttuser
	@echo "Configurando permisos del archivo de contraseñas..."
	@sudo chown mosquitto:mosquitto /etc/mosquitto/passwd
	@sudo chmod 600 /etc/mosquitto/passwd
	@echo "Configurando Mosquitto para usar autenticación..."
	@sudo bash -c 'echo "allow_anonymous false" > /etc/mosquitto/conf.d/auth.conf'
	@sudo bash -c 'echo "password_file /etc/mosquitto/passwd" >> /etc/mosquitto/conf.d/auth.conf'
	@echo "Probando configuración antes de reiniciar..."
	@sudo mosquitto -c /etc/mosquitto/mosquitto.conf -t || (echo "Error en configuración. Revirtiendo cambios..." && sudo rm -f /etc/mosquitto/conf.d/auth.conf && exit 1)
	@echo "Reiniciando Mosquitto..."
	@sudo systemctl restart mosquitto
	@echo "¡Autenticación configurada! Usuario: mqttuser"

# Diagnosticar problemas de MQTT
mqtt-diagnose:
	@echo "=== DIAGNÓSTICO MOSQUITTO ==="
	@echo "Estado del servicio:"
	@sudo systemctl status mosquitto.service --no-pager
	@echo "\n=== CONFIGURACIÓN AUTH ==="
	@cat /etc/mosquitto/conf.d/auth.conf 2>/dev/null || echo "No se encontró auth.conf"
	@echo "\n=== ARCHIVO DE CONTRASEÑAS ==="
	@ls -la /etc/mosquitto/passwd 2>/dev/null || echo "No se encontró archivo de contraseñas"
	@echo "\n=== LOGS RECIENTES ==="
	@sudo journalctl -u mosquitto.service -n 10 --no-pager

# Reparar configuración MQTT
mqtt-fix:
	@echo "=== REPARANDO CONFIGURACIÓN MQTT ==="
	@echo "Deteniendo Mosquitto..."
	@sudo systemctl stop mosquitto
	@echo "Validando configuración..."
	@sudo mosquitto -c /etc/mosquitto/mosquitto.conf -v
	@echo "Iniciando Mosquitto..."
	@sudo systemctl start mosquitto
	@sudo systemctl status mosquitto.service --no-pager

# Agregar también los targets de diagnóstico mejorados
mqtt-test-config:
	@echo "=== PROBANDO CONFIGURACIÓN MOSQUITTO ==="
	@sudo mosquitto -c /etc/mosquitto/mosquitto.conf -t
	@echo "Configuración válida ✓"

# Resetear configuración MQTT (volver al estado original)
mqtt-reset:
	@echo "=== RESETEANDO CONFIGURACIÓN MQTT ==="
	@sudo systemctl stop mosquitto
	@sudo rm -f /etc/mosquitto/conf.d/auth.conf
	@sudo rm -f /etc/mosquitto/passwd
	@echo "Configuración reseteada. Iniciando Mosquitto sin autenticación..."
	@sudo systemctl start mosquitto
	@sudo systemctl status mosquitto.service --no-pager

# Generar certificados SSL para HTTPS
ssl-certs:
	@echo "=== GENERANDO CERTIFICADOS SSL ==="
	@mkdir -p certs/ssl
	@echo "Generando clave privada..."
	@openssl genrsa -out certs/ssl/server.key 2048
	@echo "Generando certificado autofirmado..."
	@openssl req -new -x509 -key certs/ssl/server.key -out certs/ssl/server.crt -days 365 -subj "/C=ES/ST=Madrid/L=Madrid/O=TestOrg/CN=localhost"
	@echo "Certificados SSL generados en certs/ssl/"

# Configuración completa de SSL
ssl-setup: ssl-certs
	@echo "=== CONFIGURACIÓN SSL COMPLETADA ==="
	@echo "Certificados disponibles en:"
	@ls -la certs/ssl/
	@echo "El servidor usará HTTPS en el puerto configurado"
