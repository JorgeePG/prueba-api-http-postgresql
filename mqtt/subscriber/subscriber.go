package subscriber

import (
	"context"
	"fmt"
	"os"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SubscriberInfo contiene información sobre un suscriptor activo
type SubscriberInfo struct {
	Topic      string
	Client     mqtt.Client
	CancelFunc context.CancelFunc
}

// SubscriberManager gestiona múltiples suscriptores MQTT
type SubscriberManager struct {
	subscribers map[string]*SubscriberInfo
	mu          sync.RWMutex
	brokerURL   string
}

// Instancia global del manager
var globalManager *SubscriberManager
var once sync.Once

// GetSubscriberManager retorna la instancia global del manager
func GetSubscriberManager() *SubscriberManager {
	once.Do(func() {
		globalManager = &SubscriberManager{
			subscribers: make(map[string]*SubscriberInfo),
			brokerURL:   "tcp://localhost:1883",
		}
	})
	return globalManager
}

func AddTopicSubscriber(topic string) {
	manager := GetSubscriberManager()

	// Verificar si ya existe un suscriptor para este topic
	if manager.IsSubscribed(topic) {
		log.Warn().Str("topic", topic).Msg("Ya existe un suscriptor para este topic")
		return
	}

	go func(topic string) {
		if topic == "" {
			log.Fatal().Msg("El topic no puede estar vacío")
		}
		// Configurar logger
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		log.Info().Msg("🚀 Iniciando MQTT Subscriber")

		clientID := fmt.Sprintf("go-subscriber-%s", topic)
		opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID(clientID)
		client := mqtt.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatal().Err(token.Error()).Msg("Error conectando al broker MQTT")
		}
		log.Info().Msg("🟢 Conectado al broker MQTT como suscriptor")

		// Crear contexto cancelable
		ctx, cancel := context.WithCancel(context.Background())

		// Registrar el suscriptor en el manager
		manager.AddSubscriber(topic, client, cancel)

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

		// ELIMINAR estas líneas que capturan Ctrl+C
		// c := make(chan os.Signal, 1)
		// signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		// Solo esperar la cancelación del contexto
		<-ctx.Done()
		log.Info().Str("topic", topic).Msg("🛑 Cancelación solicitada para el topic")

		// Desuscribirse del topic antes de desconectar
		if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", topic).Msg("Error al desuscribirse del topic")
		}

		client.Disconnect(250)
		manager.RemoveSubscriber(topic)
		log.Info().Str("topic", topic).Msg("👋 Subscriber finalizado")
	}(topic)
}

func DeleteTopicSubscriber(topic string) {
	manager := GetSubscriberManager()

	if topic == "" {
		log.Error().Msg("El topic no puede estar vacío")
		return
	}

	// Verificar si existe el suscriptor
	if !manager.IsSubscribed(topic) {
		log.Warn().Str("topic", topic).Msg("No existe un suscriptor para este topic")
		return
	}

	log.Info().Str("topic", topic).Msg("🚀 Desuscribiendo del topic")

	// Remover el suscriptor del manager (esto cancelará el contexto)
	if err := manager.RemoveSubscriber(topic); err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("Error al remover suscriptor")
		return
	}

	log.Info().Str("topic", topic).Msg("✅ Desuscrito del topic")
}

// GetActiveTopics devuelve los topics activos (función de conveniencia)
func GetActiveTopics() []string {
	manager := GetSubscriberManager()
	return manager.GetActiveSubscribers()
}

// DisconnectAllSubscribers desconecta todos los suscriptores (función de conveniencia)
func DisconnectAllSubscribers() {
	manager := GetSubscriberManager()
	manager.DisconnectAll()
}

// GetActiveSubscribers devuelve una lista de topics activos
func (sm *SubscriberManager) GetActiveSubscribers() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	topics := make([]string, 0, len(sm.subscribers))
	for topic := range sm.subscribers {
		topics = append(topics, topic)
	}
	return topics
}

// IsSubscribed verifica si ya existe un suscriptor para un topic
func (sm *SubscriberManager) IsSubscribed(topic string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.subscribers[topic]
	return exists
}

// RemoveSubscriber elimina un suscriptor del manager
func (sm *SubscriberManager) RemoveSubscriber(topic string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	subscriber, exists := sm.subscribers[topic]
	if !exists {
		return fmt.Errorf("no existe un suscriptor para el topic: %s", topic)
	}

	// Cancelar el contexto para detener la goroutine
	subscriber.CancelFunc()

	// Eliminar del mapa
	delete(sm.subscribers, topic)

	return nil
}

// AddSubscriber agrega un suscriptor al manager
func (sm *SubscriberManager) AddSubscriber(topic string, client mqtt.Client, cancel context.CancelFunc) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.subscribers[topic] = &SubscriberInfo{
		Topic:      topic,
		Client:     client,
		CancelFunc: cancel,
	}
}

// DisconnectAll desconecta todos los suscriptores
func (sm *SubscriberManager) DisconnectAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	log.Info().Msg("🛑 Desconectando todos los suscriptores...")

	for topic, subscriber := range sm.subscribers {
		log.Info().Str("topic", topic).Msg("Desconectando suscriptor")
		subscriber.CancelFunc()
	}

	// Limpiar el mapa
	sm.subscribers = make(map[string]*SubscriberInfo)
	log.Info().Msg("👋 Todos los suscriptores desconectados")
}
