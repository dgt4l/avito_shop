CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_salt VARCHAR(255) NOT NULL,
    coins INT CHECK (coins >= 0) NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INT NOT NULL
);

CREATE TABLE IF NOT EXISTS inventory (
    user_id INT NOT NULL,
    item_id INT NOT NULL,
    quantity INT CHECK (quantity >= 0) NOT NULL,
    PRIMARY KEY (user_id, item_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (item_id) REFERENCES items(id)
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY NOT NULL,
    from_user_id INT,
    to_user_id INT,
    amount INT CHECK (amount >= 0) NOT NULL,
    FOREIGN KEY (from_user_id) REFERENCES users(id),
    FOREIGN KEY (to_user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users USING HASH (username);
CREATE INDEX IF NOT EXISTS idx_inventory_user ON inventory (user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item ON inventory (item_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions (from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to ON transactions (to_user_id);

INSERT INTO items (name, price)
VALUES ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500)
ON CONFLICT DO NOTHING;


