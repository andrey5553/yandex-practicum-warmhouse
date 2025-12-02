-- Таблица устройств
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    device_type_id INTEGER NOT NULL,
    room_id INTEGER,
    serial_number VARCHAR UNIQUE NOT NULL,
    status VARCHAR DEFAULT 'offline',
    configuration JSONB DEFAULT '{}',
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_devices_room_id ON devices(room_id);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_serial_number ON devices(serial_number);
CREATE INDEX IF NOT EXISTS idx_devices_is_active ON devices(is_active);
CREATE INDEX IF NOT EXISTS idx_devices_created_at ON devices(created_at);

-- Вставляем тестовые данные (опционально)
INSERT INTO devices (name, type, device_type_id, room_id, serial_number, status, configuration) 
VALUES 
('Умная лампа гостиная', 'light', 1, 1, 'LAMP-001', 'online', '{"brightness": 80}'),
('Кондиционер спальня', 'ac', 2, 2, 'AC-001', 'online', '{"temperature": 22}'),
('Дачтик движения коридор', 'sensor', 3, 3, 'SENSOR-001', 'offline', '{"sensitivity": "high"}')
ON CONFLICT (serial_number) DO NOTHING;
