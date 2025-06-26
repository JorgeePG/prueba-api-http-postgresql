package repository

import (
	"database/sql"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"
)

type MqttMessageRepository struct {
	db *sql.DB
}

func NewMqttMessageRepository(db *sql.DB) *MqttMessageRepository {
	return &MqttMessageRepository{db: db}
}

func (r *MqttMessageRepository) Create(message *models.MqttMessage) error {
	query := `
        INSERT INTO mqtt_messages (topic, payload, qos, retained)
        VALUES ($1, $2, $3, $4)
        RETURNING id, received_at
    `

	err := r.db.QueryRow(
		query,
		message.Topic,
		message.Payload,
		message.QOS,
		message.Retained,
	).Scan(&message.ID, &message.ReceivedAt)

	return err
}

func (r *MqttMessageRepository) GetByTopic(topic string, limit int) ([]models.MqttMessage, error) {
	query := `
        SELECT id, topic, payload, received_at, qos, retained
        FROM mqtt_messages
        WHERE topic = $1
        ORDER BY received_at DESC
        LIMIT $2
    `

	rows, err := r.db.Query(query, topic, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.MqttMessage
	for rows.Next() {
		var msg models.MqttMessage
		err := rows.Scan(
			&msg.ID,
			&msg.Topic,
			&msg.Payload,
			&msg.ReceivedAt,
			&msg.QOS,
			&msg.Retained,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (r *MqttMessageRepository) GetAll(limit int) ([]models.MqttMessage, error) {
	query := `
        SELECT id, topic, payload, received_at, qos, retained
        FROM mqtt_messages
        ORDER BY received_at DESC
        LIMIT $1
    `

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.MqttMessage
	for rows.Next() {
		var msg models.MqttMessage
		err := rows.Scan(
			&msg.ID,
			&msg.Topic,
			&msg.Payload,
			&msg.ReceivedAt,
			&msg.QOS,
			&msg.Retained,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
