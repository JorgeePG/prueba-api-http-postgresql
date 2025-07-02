package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configurar logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("🚀 Iniciando MQTT Publisher")

	// Cargar el certificado CA del broker
	caCert, err := os.ReadFile("./cert/ca.crt")
	if err != nil {
		log.Fatal().Err(err).Msg("No se pudo leer el certificado CA")
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatal().Msg("No se pudo agregar el certificado CA al pool")
	}

	// Intentar cargar certificado de cliente (opcional)
	var tlsConfig *tls.Config
	clientCert, err := tls.LoadX509KeyPair("./cert/client.crt", "./cert/client.key")
	if err != nil {
		log.Warn().Err(err).Msg("No se encontraron certificados de cliente, usando solo CA")
		tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true, // Solo para pruebas, ponlo en false en producción
			ServerName:         "localhost",
		}
	} else {
		tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true, // Solo para pruebas, ponlo en false en producción
			ServerName:         "localhost",
		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker("ssl://localhost:8883").
		SetClientID("go-publisher").
		SetTLSConfig(tlsConfig).
		SetUsername("publisher").
		SetPassword("publisher")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Error conectando al broker MQTT")
	}
	log.Info().Msg("🔵 Conectado al broker MQTT como publicador")

	// Definir múltiples topics para pruebas
	topics := []string{
		"test/topic",
		"sensors/temperature",
		"notifications/alerts",
	}

	// Estructura para el mensaje JSON
	type MQTTMessage struct {
		ID        int       `json:"id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Topic     string    `json:"topic"`
		Source    string    `json:"source"`
	}

	for i := 1; ; i++ {
		// Rotar entre los topics disponibles
		currentTopic := topics[(i-1)%len(topics)]

		// Crear mensaje estructurado
		msgData := MQTTMessage{
			ID:        i,
			Message:   fmt.Sprintf("Mensaje #%d desde Publisher", i),
			Timestamp: time.Now(),
			Topic:     currentTopic,
			Source:    "go-publisher",
		}

		// Convertir a JSON
		jsonMsg, err := json.Marshal(msgData)
		if err != nil {
			log.Error().Err(err).Msg("Error al serializar mensaje JSON")
			continue
		}

		// Publicar mensaje
		token := client.Publish(currentTopic, 1, false, string(jsonMsg))
		token.Wait()

		log.Info().
			Str("topic", currentTopic).
			Str("json_message", string(jsonMsg)).
			Msg("📤 Mensaje JSON publicado")

		time.Sleep(1 * time.Second)
	}
}
