package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	log.Info().Msg("🚀 Iniciando MQTT Subscriber")

	brokerURL := fmt.Sprintf("tcp://%s:%d", *host, *port)
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("go-subscriber")
	client := mqtt.NewClient(opts)

	log.Info().Str("broker", brokerURL).Msg("🔗 Conectando al broker...")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Error conectando al broker MQTT")
	}
	log.Info().Msg("🟢 Conectado al broker MQTT como suscriptor")

	topic := "test/topic"
	token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		log.Info().
			Str("topic", msg.Topic()).
			Str("payload", string(msg.Payload())).
			Msg("📥 Mensaje recibido")
	})

	if token.Wait() && token.Error() != nil {
		log.Fatal().Err(token.Error()).Msg("Error suscribiéndose al topic")
	}

	log.Info().Str("topic", topic).Msg("✅ Suscrito al topic")

	// Esperar señal de interrupción
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info().Msg("🛑 Desconectando...")
	client.Disconnect(250)
	log.Info().Msg("👋 Subscriber finalizado")
}
