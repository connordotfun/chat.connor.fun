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
		INSERT INTO rooms (id, name, lat, lon, radius) VALUES (:id, :name, :lat, :lon, :radius);
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
		SELECT rooms.*, earth_distance(ll_to_earth(:lat, :lon), ll_to_earth(lat, lon)) as distance FROM rooms
			WHERE earth_box(ll_to_earth(:lat, :lon), :radius + radius) @> ll_to_earth(rooms.lat, rooms.lon)
			ORDER BY distance ASC;
	`
)
