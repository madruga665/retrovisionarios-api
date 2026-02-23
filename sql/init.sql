CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    flyer TEXT
);

INSERT INTO events (date, name, location, flyer) VALUES
('2025-12-08', 'Aniversário Moto Club Dragões', 'Sede Moto Club Dragões', 'https://aws.bucket.com/foto/1'),
(CURRENT_DATE, 'Aniversário Moto Club Piratas', 'Sede Moto Club Piratas', 'https://aws.bucket.com/foto/2');
