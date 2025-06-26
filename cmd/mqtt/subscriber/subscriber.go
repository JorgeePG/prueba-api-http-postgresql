package subscriber

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartSubscriber() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-subscriber")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("🟢 Conectado al broker MQTT como suscriptor")

	topic := "test/topic"
	client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📥 Mensaje recibido en %s: %s\n", msg.Topic(), msg.Payload())
	})

	// Mantener el suscriptor vivo
	for {
		time.Sleep(1 * time.Second)
	}
}
