package subscriber

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/repository"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	db          *sql.DB
	mqttRepo    *repository.MqttMessageRepository
	tlsConfig   *tls.Config
}

// Instancia global del manager
var globalManager *SubscriberManager
var once sync.Once

// GetSubscriberManager retorna la instancia global del manager
func GetSubscriberManager() *SubscriberManager {
	once.Do(func() {
		globalManager = &SubscriberManager{
			subscribers: make(map[string]*SubscriberInfo),
			brokerURL:   "ssl://localhost:8883", // Cambiar a SSL
		}
		// Configurar TLS al inicializar
		globalManager.setupTLS()
	})
	return globalManager
}

// SetDatabase configura la base de datos para el manager
func (sm *SubscriberManager) SetDatabase(db *sql.DB) {
	sm.db = db
	sm.mqttRepo = repository.NewMqttMessageRepository(db)
}

func AddTopicSubscriber(topic string) error {
	manager := GetSubscriberManager()

	// Validación más estricta del topic
	if topic == "" || len(topic) == 0 {
		log.Error().Msg("❌ El topic no puede estar vacío")
		return fmt.Errorf("el topic no puede estar vacío")
	}

	// Verificar si ya existe un suscriptor para este topic
	if manager.IsSubscribed(topic) {
		log.Warn().Str("topic", topic).Msg("⚠️ Ya existe un suscriptor para este topic")
		return fmt.Errorf("ya existe un suscriptor para el topic: %s", topic)
	}

	// Verificar que TLS esté configurado antes de proceder
	if manager.tlsConfig == nil {
		log.Error().Str("topic", topic).Msg("❌ TLS no está configurado. No se puede proceder con la suscripción")
		return fmt.Errorf("TLS no está configurado")
	}

	// Log para debug del topic recibido
	log.Info().
		Str("topic", topic).
		Str("topic_length", fmt.Sprintf("%d", len(topic))).
		Str("topic_bytes", fmt.Sprintf("%+v", []byte(topic))).
		Msg("🔍 Topic recibido para suscripción")

	go func(topic string) {
		if topic == "" {
			log.Error().Msg("El topic no puede estar vacío")
			return
		}

		log.Info().Str("topic", topic).Msg("🚀 Iniciando MQTT Subscriber")

		// Usar la misma configuración que el publisher
		clientID := fmt.Sprintf("go-subscriber-%s-%d", topic, time.Now().UnixNano())

		// Verificar que TLS esté configurado
		if manager.tlsConfig == nil {
			log.Error().Str("topic", topic).Msg("❌ TLS no está configurado correctamente")
			return
		}

		opts := mqtt.NewClientOptions().
			AddBroker("ssl://localhost:8883").
			SetClientID(clientID).
			SetTLSConfig(manager.tlsConfig).
			SetUsername("publisher").
			SetPassword("publisher").
			SetConnectTimeout(10 * time.Second).
			SetKeepAlive(30 * time.Second).
			SetPingTimeout(5 * time.Second).
			SetWriteTimeout(5 * time.Second).
			SetAutoReconnect(true).
			SetMaxReconnectInterval(5 * time.Second).
			SetConnectionLostHandler(func(client mqtt.Client, err error) {
				log.Error().Err(err).Str("topic", topic).Msg("🔴 Conexión MQTT perdida")
			}).
			SetOnConnectHandler(func(client mqtt.Client) {
				log.Info().Str("topic", topic).Msg("🟢 Cliente MQTT reconectado")
			})

		client := mqtt.NewClient(opts)

		log.Info().Str("broker", "ssl://localhost:8883").Msg("🔌 Intentando conectar con SSL...")

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", topic).Msg("❌ Error conectando al broker MQTT")
			return
		}
		log.Info().Str("topic", topic).Msg("🟢 Conectado al broker MQTT como suscriptor")

		// Crear contexto cancelable
		ctx, cancel := context.WithCancel(context.Background())

		// Registrar el suscriptor en el manager ANTES de suscribirse
		manager.AddSubscriber(topic, client, cancel)

		// Verificar conexión antes de suscribirse
		if !client.IsConnected() {
			log.Error().Str("topic", topic).Msg("❌ Cliente no está conectado")
			manager.RemoveSubscriber(topic)
			return
		}

		token := client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
			log.Info().Str("topic", msg.Topic()).Str("payload", string(msg.Payload())).Msg("🔥 CALLBACK ZEROLOG")

			// Obtener el manager dentro del callback
			mgr := GetSubscriberManager()

			// Crear el mensaje para guardar en BD
			mqttMessage := &models.MqttMessage{
				Topic:    msg.Topic(),
				Payload:  string(msg.Payload()),
				QOS:      int(msg.Qos()),
				Retained: msg.Retained(),
			}

			// Guardar en la base de datos
			if mgr.mqttRepo != nil {
				if err := mgr.mqttRepo.Create(mqttMessage); err != nil {
					log.Error().
						Err(err).
						Str("topic", msg.Topic()).
						Msg("❌ Error guardando mensaje en base de datos")
				} else {
					log.Info().
						Str("topic", msg.Topic()).
						Int("message_id", mqttMessage.ID).
						Str("payload", string(msg.Payload())).
						Msg("💾 Mensaje guardado en base de datos")
				}
			} else {
				log.Warn().Msg("⚠️  Base de datos no configurada, solo registrando mensaje")
				log.Info().
					Str("topic", msg.Topic()).
					Str("payload", string(msg.Payload())).
					Msg("📥 Mensaje recibido")
			}
		})

		if token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Str("topic", topic).Msg("❌ Error suscribiéndose al topic")
			manager.RemoveSubscriber(topic)
			return
		}

		log.Info().Str("topic", topic).Msg("✅ Suscrito al topic correctamente")

		// Esperar cancelación
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

	return nil
}

func DeleteTopicSubscriber(topic string) error {
	manager := GetSubscriberManager()

	if topic == "" {
		log.Error().Msg("El topic no puede estar vacío")
		return fmt.Errorf("el topic no puede estar vacío")
	}

	// Verificar si existe el suscriptor
	if !manager.IsSubscribed(topic) {
		log.Warn().Str("topic", topic).Msg("No existe un suscriptor para este topic")
		return fmt.Errorf("no existe un suscriptor para el topic: %s", topic)
	}

	log.Info().Str("topic", topic).Msg("🚀 Desuscribiendo del topic")

	// Remover el suscriptor del manager (esto cancelará el contexto)
	if err := manager.RemoveSubscriber(topic); err != nil {
		log.Error().Err(err).Str("topic", topic).Msg("Error al remover suscriptor")
		return fmt.Errorf("error al remover suscriptor: %w", err)
	}

	log.Info().Str("topic", topic).Msg("✅ Desuscrito del topic")
	return nil
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

// ListMqttMessages obtiene los mensajes MQTT guardados en la base de datos (función de conveniencia)
func ListMqttMessages(limit int) ([]models.MqttMessage, error) {
	manager := GetSubscriberManager()
	return manager.ListMqttMessages(limit)
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

// ListMqttMessages obtiene todos los mensajes MQTT guardados en la base de datos
func (sm *SubscriberManager) ListMqttMessages(limit int) ([]models.MqttMessage, error) {
	if sm.mqttRepo == nil {
		return nil, fmt.Errorf("base de datos no configurada")
	}

	if limit <= 0 {
		limit = 100 // Límite por defecto
	}

	messages, err := sm.mqttRepo.GetAll(limit)
	if err != nil {
		log.Error().Err(err).Msg("❌ Error obteniendo mensajes MQTT de la base de datos")
		return nil, err
	}

	log.Info().
		Int("count", len(messages)).
		Int("limit", limit).
		Msg("📋 Mensajes MQTT obtenidos de la base de datos")

	return messages, nil
}

// setupTLS configura los certificados TLS para el manager
func (sm *SubscriberManager) setupTLS() {
	log.Info().Msg("🔧 [TLS] Iniciando configuración TLS...")

	// 1. Usar la misma ruta que el publisher para consistencia
	basePath := os.Getenv("CERT_PATH")
	if basePath == "" {
		basePath = "mqtt/publisher/cert"
	}

	caPath := basePath + "/ca.crt"
	log.Info().Str("path", caPath).Msg("[TLS] Intentando leer certificado CA")

	caCert, err := os.ReadFile(caPath)
	if err != nil {
		log.Error().Err(err).Str("path", caPath).Msg("[TLS] ❌ No se pudo leer el certificado CA")
		// Intentar con ruta alternativa
		altPath := "./certs/ssl/server.crt"
		log.Warn().Str("alt_path", altPath).Msg("[TLS] Intentando ruta alternativa")
		caCert, err = os.ReadFile(altPath)
		if err != nil {
			log.Error().Err(err).Str("alt_path", altPath).Msg("[TLS] ❌ Tampoco se encontró certificado en ruta alternativa")
			return
		}
		caPath = altPath
	}

	// 2. Crear el pool de CA
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Error().Str("path", caPath).Msg("[TLS] ❌ No se pudo agregar el certificado CA al pool")
		return
	}
	log.Info().Msg("[TLS] ✅ CA agregada correctamente al pool")

	// 3. Intentar cargar certificado de cliente (opcional)
	clientCrtPath := basePath + "/client.crt"
	clientKeyPath := basePath + "/client.key"
	log.Info().Str("crt", clientCrtPath).Str("key", clientKeyPath).Msg("[TLS] Intentando cargar certificado de cliente")

	clientCert, err := tls.LoadX509KeyPair(clientCrtPath, clientKeyPath)
	if err != nil {
		log.Warn().Err(err).
			Str("crt", clientCrtPath).
			Str("key", clientKeyPath).
			Msg("[TLS] ⚠️ No se encontraron certificados de cliente, usando solo CA")
		sm.tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true, // Solo para pruebas, ponlo en false en producción
			ServerName:         "localhost",
		}
		log.Info().Msg("[TLS] ✅ TLS configurado solo con CA")
	} else {
		log.Info().Str("crt", clientCrtPath).Str("key", clientKeyPath).Msg("[TLS] ✅ Certificado de cliente cargado correctamente")
		sm.tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			Certificates:       []tls.Certificate{clientCert},
			InsecureSkipVerify: true, // Solo para pruebas, ponlo en false en producción
			ServerName:         "localhost",
		}
		log.Info().Msg("[TLS] ✅ TLS configurado con CA y certificado de cliente")
	}

	log.Info().Msg("✅ [TLS] Configuración TLS completada exitosamente")
}
