package models

import (
	"time"
)

type MqttMessage struct {
	ID         int       `json:"id" db:"id"`
	Topic      string    `json:"topic" db:"topic"`
	Payload    string    `json:"payload" db:"payload"`
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
	QOS        int       `json:"qos" db:"qos"`
	Retained   bool      `json:"retained" db:"retained"`
}
