CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS images (
    id          UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_email  VARCHAR(255) NOT NULL, 
    key         VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id       UUID NOT NULL PRIMARY KEY DEFAULT uuid_generate_v1(), 
    email    VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);
