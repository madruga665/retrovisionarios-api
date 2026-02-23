CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    name VARCHAR(255) NOT NULL,
    flyer TEXT
);

INSERT INTO events (date, name, flyer) VALUES
('2025-12-08', 'Aniversário Moto Club Dragões', 'https://aws.bucket.com/foto/1'),
(CURRENT_DATE, 'Aniversário Moto Club Piratas', 'https://aws.bucket.com/foto/2');
