package user


const (
	createIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL UNIQUE NOT NULL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			secret VARCHAR(255) NOT NULL
		);
	`
	getAllUsersQuery = `
		SELECT * FROM users;
	`

	getUserByIdQuery = `
		SELECT * FROM users
			WHERE id = :id;
	`
	getUserByUsernameQuery = `
		SELECT * FROM users
			WHERE username = :username;
	`

	insertUserQuery = `
		INSERT INTO users (username, secret) VALUES (:username, :secret);
	`

	getLastInsertedQuery = `
		SELECT currval(pg_get_serial_sequence('users','id'));
	`
)