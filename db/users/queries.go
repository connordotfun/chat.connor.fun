package users


const (
	createIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID UNIQUE NOT NULL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
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
		INSERT INTO users (id, username, email, secret) VALUES (:id, :username, :email, :secret);
	`

	getUserByEmailQuery = `
		SELECT * FROM users
			WHERE username = :email;
	`
)