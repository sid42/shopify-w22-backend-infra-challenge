CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS images (
    id          UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1(),
    email       VARCHAR(255) NOT NULL, 
    link_to_s3  VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id       UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1(), 
    email    VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- INSERT INTO users VALUES (1, 'test', 'test-p');