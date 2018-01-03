package messages


const (
	createIfNotExistsQuery = `
		CREATE TABLE IF NOT EXISTS messages (
			id UUID NOT NULL PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users (id),
			room_id UUID NOT NULL REFERENCES rooms (id),
			text TEXT NOT NULL,
			create_date Integer NOT NULL
		);
	`

	insertMessageQuery = `
		INSERT INTO messages (id, user_id, room_id, text, create_date) VALUES
			(:id, :user_id, :room_id, :text, :create_date);
	`

	updateMessageTextQuery = `
		UPDATE messages SET
			text = :text
		WHERE id = :id
		RETURNING *;
	`

	selectOneByIdQuery = `
		SELECT * FROM messages
			WHERE id = :id
	`

	selectByUserIdQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.user_id = :user_id;
	`

	selectTopByUserIdQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.user_id = :user_id
		LIMIT :count;
	`

	selectByRoomIdQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.room_id = :room_id;
	`

	selectTopByRoomQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.room_id = :room_id
		LIMIT :count;
	`

	selectByUserAndRoomQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.room_id = :room_id
				AND m.user_id = :user_id;
	`

	selectTopByUserAndRoomQuery = `
		SELECT m.id, m.text, m.create_date, u.id as user_id, u.username, r.id as room_id, r.name as room_name
				FROM messages as m
			JOIN users as u ON
				u.id = m.user_id
			JOIN rooms as r ON
				m.room_id = r.id
			WHERE m.room_id = :room_id
				AND m.user_id = :user_id
		LIMIT :count;
	`
)