package rooms


const (
	createIfNotExistsRoomsQuery = `
		CREATE TABLE IF NOT EXISTS rooms (
			id SERIAL UNIQUE NOT NULL PRIMARY KEY,
			name VARCHAR(20) UNIQUE NOT NULL
		);
	`

	insertRoomQuery = `
		INSERT INTO rooms (name) VALUES (:name);
	`

	selectRoomByNameQuery = `
		SELECT * FROM rooms
			WHERE name = :name;
	`
)
