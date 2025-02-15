package repository

const (
	getFromUsers = `SELECT id, username, password_salt FROM users WHERE username=$1;`

	insertToUsers = `INSERT INTO users (username, password_salt, coins) values ($1, $2, $3) RETURNING id;`

	getFromItems = `SELECT id, name, price FROM items WHERE name = $1`

	getCoinsFromUser = `SELECT coins from users WHERE id = $1 FOR UPDATE`

	updateCoinsFromUser = `UPDATE users SET coins = coins - $1 WHERE id = $2`

	insertToInventory = `INSERT INTO inventory (user_id, item_id, quantity) VALUES ($1, $2, 1) ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = inventory.quantity + 1`

	getCoins = `SELECT coins from users where id = $1`

	getUserInventory = `SELECT i.name, quantity from inventory INNER JOIN items i ON i.id = inventory.item_id WHERE user_id = $1`

	getUserRecieved = `SELECT username as from_user, amount FROM transactions INNER JOIN users ON users.id = transactions.from_user_id WHERE to_user_id = $1`

	getUserSent = `SELECT username as to_user, amount FROM transactions INNER JOIN users ON users.id = transactions.to_user_id WHERE from_user_id = $1`

	getIdFromUsers = `SELECT id from users WHERE username = $1`

	getCoinsToUser = `SELECT coins from users WHERE username = $1 FOR UPDATE`

	updateCoinsToUser = `UPDATE users SET coins = coins + $1 WHERE username = $2`

	insertToTransactions = `INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)`
)
