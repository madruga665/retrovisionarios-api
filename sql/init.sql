CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL,
    flyer TEXT
);

INSERT INTO events (date, name, flyer) VALUES
('2025-12-08 18:30:00', 'Aniversário Moto Club Dragões', 'http://aws.bucket.com/foto/1'),
(NOW(), 'Aniversário Moto Club Piratas', 'http://aws.bucket.com/foto/2');
