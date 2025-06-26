package subscriber

import (
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func AddTopicSubscriber(topic string) {
	go func(topic string) {
		if topic == "" {
			log.Fatal().Msg("El topic no puede estar vacío")
		}
		// Configurar logger
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("🚀 Iniciando MQTT Subscriber")

		opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-subscriber")
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal().Err(token.Error()).Msg("Error conectando al broker MQTT")
		}
		log.Info().Msg("🟢 Conectado al broker MQTT como suscriptor")

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
	}(topic)
}
