package main

import (
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

	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-publisher")
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

	// Publicar mensajes cada 2 segundos
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
		token := client.Publish(currentTopic, 0, false, string(jsonMsg))
		token.Wait()

		log.Info().
			Str("topic", currentTopic).
			Str("json_message", string(jsonMsg)).
			Msg("📤 Mensaje JSON publicado")

		time.Sleep(1 * time.Second)
	}
}
