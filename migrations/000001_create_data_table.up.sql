CREATE TABLE credentials (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    alias VARCHAR(255) UNIQUE,
    info TEXT
);

CREATE TABLE payment_card (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    exp_month SMALLINT NOT NULL,
    exp_year SMALLINT NOT NULL,
    cvv SMALLINT,
    alias VARCHAR(255) UNIQUE,
    info TEXT
);

CREATE TABLE binary_data (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    data BYTEA NOT NULL,
    alias VARCHAR(255) UNIQUE NOT NULL,
    info TEXT
);

CREATE TABLE text_data (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    data VARCHAR(255) NOT NULL,
    alias VARCHAR(255) UNIQUE NOT NULL,
    info TEXT
);
