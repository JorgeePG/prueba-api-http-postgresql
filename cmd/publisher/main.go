package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configurar flags
	port := flag.Int("port", 1883, "Puerto del broker MQTT")
	host := flag.String("host", "localhost", "Host del broker MQTT")
	flag.Parse()

	// Configurar logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("🚀 Iniciando MQTT Publisher")

	brokerURL := fmt.Sprintf("tcp://%s:%d", *host, *port)
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("go-publisher")
	client := mqtt.NewClient(opts)

	log.Info().Str("broker", brokerURL).Msg("🔗 Conectando al broker...")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Error conectando al broker MQTT")
	}
	log.Info().Msg("🔵 Conectado al broker MQTT como publicador")

	topic := "test/topic"

	// Publicar mensajes cada 2 segundos
	for i := 1; ; i++ {
		msg := fmt.Sprintf("Mensaje #%d desde Publisher - %s", i, time.Now().Format("15:04:05"))
		token := client.Publish(topic, 0, false, msg)
		token.Wait()
		log.Info().Str("mensaje", msg).Str("topic", topic).Msg("📤 Mensaje publicado")
		time.Sleep(2 * time.Second)
	}
}
