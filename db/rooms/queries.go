package rooms


const (
	createIfNotExistsRoomsQuery = `
		CREATE TABLE IF NOT EXISTS rooms (
			id UUID UNIQUE NOT NULL PRIMARY KEY,
			name VARCHAR(20) UNIQUE NOT NULL,
			lat Float NOT NULL,
			lon Float NOT NULL,
			radius Float NOT NULL
		);
	`

	createCubeQuery = `
		CREATE EXTENSION IF NOT EXISTS cube;
	`
	createEarthDistQuery = `
		CREATE EXTENSION IF NOT EXISTS earthdistance;
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

	selectWithinRadiusQuery = `
		SELECT rooms.*, earth_distance(ll_to_earth(:lat, :lng), ll_to_earth(lat, lon)) as distance FROM rooms
			WHERE earth_box(ll_to_earth(:lat, :lng), :radius + radius) @> ll_to_earth(rooms.lat, rooms.lon)
			ORDER BY distance ASC;
	`
)
