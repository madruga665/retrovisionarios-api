-- 1. Criação da tabela de eventos
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    flyer TEXT,
    deleted BOOLEAN DEFAULT FALSE
);

-- 2. Inserção de dados em eventos (Note o ";" no final)
INSERT INTO events (date, name, location, flyer) VALUES
('2025-12-08 20:00', 'Aniversário Moto Club Dragões', 'Sede Moto Club Dragões', 'https://aws.bucket.com/foto/1'),
(CURRENT_TIMESTAMP, 'Aniversário Moto Club Piratas', 'Sede Moto Club Piratas', 'https://aws.bucket.com/foto/2');

-- 3. Criar o tipo enumerado
-- É comum usar letras maiúsculas ou minúsculas, mas os valores devem ser exatos
CREATE TYPE video_category AS ENUM ('ORIGINAL SONG', 'COVER');

-- 4. Criação da tabela de vídeos
CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255) NOT NULL,
    video_src VARCHAR(255) NOT NULL,
    category video_category NOT NULL
);

-- 5. Inserção de dados em vídeos
INSERT INTO videos (title, subtitle, video_src, category) VALUES
('A Estrada', 'Clipe Oficial • 2012', 'https://www.youtube.com/embed/s1hzD8fUxek?si=eJdaWH5jCvUoQH8M', 'ORIGINAL SONG'),
('Resistência', 'Clipe Oficial • 2012', 'https://www.youtube.com/embed/kbBbCEVqmp8?si=MtZ_AqSqIP4H0ava', 'ORIGINAL SONG'),
('Resistência', 'Clipe Oficial • 2012', 'https://www.youtube.com/embed/kbBbCEVqmp8?si=MtZ_AqSqIP4H0ava', 'COVER');