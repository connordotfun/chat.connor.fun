package roles


const (
	createIfNotExistsQuery = `
		CREATE TABLE user_roles IF NOT EXISTS (
			user_id INTEGER NOT NULL REFERENCES users (id),
			role VARCHAR(255) NOT NULL,
			PRIMARY KEY (user_id, role)
		);
	`

	insertRelationQuery = `
		INSERT INTO user_roles (user_id, role) VALUES (:user_id, :role);
	`
)
