CREATE TABLE IF NOT EXISTS mqtt_messages (
    id SERIAL PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    payload TEXT NOT NULL,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    qos INTEGER DEFAULT 0,
    retained BOOLEAN DEFAULT FALSE
);

-- Índices para mejorar el rendimiento
CREATE INDEX IF NOT EXISTS idx_mqtt_messages_topic ON mqtt_messages(topic);
CREATE INDEX IF NOT EXISTS idx_mqtt_messages_received_at ON mqtt_messages(received_at);