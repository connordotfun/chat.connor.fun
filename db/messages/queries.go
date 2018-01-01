package messages


const (
	createIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL NOT NULL PRIMARY KEY,
			user_id Integer NOT NULL REFERENCES users (id),
			room_id Integer NOT NULL REFERENCES rooms (id),
			text TEXT NOT NULL,
			create_date Integer NOT NULL,
		);
	`

	insertMessageQuery = `
		INSERT INTO messages (user_id, room_id, text, create_date) VALUES
			(:user_id, :room_id, :text, :create_date);
	`

	selectByUserIdQuery = `
		SELECT * FROM messages
			WHERE user_id = :id;
	`

	selectByRoomIdQuery = `
		SELECT * FROM messages
			WHERE room_id = :id;
	`

	selectByUserAndRoomQuery = `
		SELECT * FROM messages
			WHERE user_id = :user_id;
			AND
			room_id = :room_id;
	`
)