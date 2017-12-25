package user


const (
	CreateIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL UNIQUE NOT NULL PRIMARY KEY
			username NVARCHAR(255) UNIQUE NOT NULL
			secret VARCHAR(255) NOT NULL
		);
	`
	GetAllUsersQuery = `
		SELECT * FROM users;
	`

	GetUserByIdQuery = `
		SELECT * FROM users
			WHERE id = :id;
	`
	GetUserByUsernameQuery = `
		SELECT * FROM users
			WHERE username = :username;
	`

	InsertUserQuery = `
		INSERT INTO users (username, secret) VALUES (:username, :secret);
	`
)