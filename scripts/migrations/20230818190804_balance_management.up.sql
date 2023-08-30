CREATE TYPE session_type AS (
    refresh_token   VARCHAR(255),
    expires_at      TIMESTAMP
);

CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    login           VARCHAR(255) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password        VARCHAR(255) NOT NULL,
    session         session_type,
    registered_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE accounts (
    id              SERIAL PRIMARY KEY,
    user_id         INT UNIQUE NOT NULL,
    balance         NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE products (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) UNIQUE NOT NULL,
    price           NUMERIC(10, 2) NOT NULL
);

CREATE TYPE order_status AS ENUM ('accepted', 'processing', 'confirmed', 'cancelled');

CREATE TABLE orders (
    id              SERIAL PRIMARY KEY,
    account_id      INT NOT NULL,
    products        INT[] DEFAULT ARRAY[]::INT[],
    amount          NUMERIC(10, 2) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    status          order_status DEFAULT 'accepted',
    FOREIGN KEY (account_id) REFERENCES accounts (id)
--     FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE reserves (
    id              SERIAL PRIMARY KEY,
    order_id        INT UNIQUE NOT NULL,
    amount          NUMERIC(10, 2) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

-- SELECT o.account_id, p.name, p.price, o.created_at FROM orders o
-- JOIN reserves r on o.id = r.order_id
-- JOIN products p on p.id = o.product_id;

CREATE TYPE operation_type AS ENUM ('revenue', 'refund');

CREATE TABLE operations (
    id             SERIAL PRIMARY KEY,
    account_id     INT NOT NULL,
    order_id       INT NOT NULL,
    amount         NUMERIC(10, 2) NOT NULL,
    type           operation_type NOT NULL,
    order_date     TIMESTAMP NOT NULL,
    description    VARCHAR(255) DEFAULT NULL,
    FOREIGN KEY (account_id) REFERENCES accounts (id)
);
