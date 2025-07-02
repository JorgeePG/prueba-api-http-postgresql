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

	log.Info().Msg("🚀 [TEST MODE] Iniciando MQTT Publisher sin SSL")

	// Configuración simple sin SSL para testing
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
		SetClientID("test-publisher").
		SetConnectTimeout(10 * time.Second).
		SetKeepAlive(30 * time.Second).
		SetAutoReconnect(true)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("❌ [TEST MODE] Error conectando al broker MQTT")
	}
	log.Info().Msg("🔵 [TEST MODE] Conectado al broker MQTT como publicador")

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

	log.Info().Msg("🔄 [TEST MODE] Iniciando envío de mensajes cada 2 segundos...")

	for i := 1; ; i++ {
		// Rotar entre los topics disponibles
		currentTopic := topics[(i-1)%len(topics)]

		// Crear mensaje estructurado
		msgData := MQTTMessage{
			ID:        i,
			Message:   fmt.Sprintf("[TEST MODE] Mensaje #%d desde Publisher", i),
			Timestamp: time.Now(),
			Topic:     currentTopic,
			Source:    "test-publisher",
		}

		// Convertir a JSON
		jsonMsg, err := json.Marshal(msgData)
		if err != nil {
			log.Error().Err(err).Msg("❌ [TEST MODE] Error al serializar mensaje JSON")
			continue
		}

		// Publicar mensaje
		token := client.Publish(currentTopic, 1, false, string(jsonMsg))
		if token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", currentTopic).Msg("❌ [TEST MODE] Error publicando mensaje")
		} else {
			log.Info().
				Str("topic", currentTopic).
				Str("message", string(jsonMsg)).
				Msg("📤 [TEST MODE] Mensaje JSON publicado exitosamente")
		}

		time.Sleep(2 * time.Second)
	}
}
