package rooms


const (
	createIfNotExistsRoomsQuery = `
		CREATE TABLE IF NOT EXISTS rooms (
			id UUID UNIQUE NOT NULL PRIMARY KEY,
			name VARCHAR(20) UNIQUE NOT NULL
		);
	`

	insertRoomQuery = `
		INSERT INTO rooms (id, name) VALUES (:id, :name);
	`

	selectRoomByNameQuery = `
		SELECT * FROM rooms
			WHERE name = :name;
	`

	selectRoomByIdQuery = `
		SELECT * FROM rooms
			WHERE id = :id;
	`
)
