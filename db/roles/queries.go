package roles


const (
	createIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id UUID NOT NULL REFERENCES users (id),
			role VARCHAR(255) NOT NULL,
			PRIMARY KEY (user_id, role)
		);
	`

	insertRelationQuery = `
		INSERT INTO user_roles (user_id, role) VALUES (:user_id, :role);
	`

	getRolesByUserQuery = `
		SELECT role FROM user_roles
			WHERE user_id = :user_id;
	`

	deleteRoleQuery = `
		DELETE FROM user_roles
			WHERE user_id = :user_id
				AND role = :role;
	`
)
