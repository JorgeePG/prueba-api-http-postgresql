package publisher

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func StartPublisher() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("go-publisher")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("🔵 Conectado al broker MQTT como publicador")

	topic := "test/topic"

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Mensaje #%d", i)
		token := client.Publish(topic, 0, false, msg)
		token.Wait()
		fmt.Printf("📤 Publicado: %s\n", msg)
		time.Sleep(1 * time.Second)
	}

	client.Disconnect(250)
}
