package users

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	_"github.com/lib/pq"
	"github.com/aaronaaeng/chat.connor.fun/model"
)

const (
	testDbHost = "localhost"
	testDbPort = 5432
	testDbName = "connorfuntest_userrepo"
)

var testDb *sqlx.DB

func TestMain(m *testing.M) {
	//Init the connorfuntest db and reconnect to the new db
	db, err := sqlx.Open("postgres", "postgresql://localhost:5432?sslmode=disable")
	if err != nil {
		panic("failed to establish db connection: " + err.Error())
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS connorfuntest_userrepo")
	if err != nil {
		panic(err) //Something went horribly wrong
	}

	_, err = db.Exec("CREATE DATABASE connorfuntest_userrepo;")
	if err != nil {
		panic("Failed to create test db: " + err.Error())
	}

	testDb, err = sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", testDbHost, testDbPort, testDbName))
	if err != nil {
		panic("failed to establish db connection")
	}

	m.Run()

	//clean up the connorfuntestdb
	testDb.Close()
	cleanUpTestDb(db)

	db.Close()
}

func cleanUpTestDb(db *sqlx.DB) {
	_, err:= db.Query(`SELECT
			pg_terminate_backend (pg_stat_activity.pid)
			FROM
			pg_stat_activity
			WHERE
			pg_stat_activity.datname = 'connorfuntest_userrepo'`)

	if err != nil {
		panic("FAILED TO KILL BG CONNECTIONS")
	}

	_, err = db.Exec("DROP DATABASE connorfuntest_userrepo;")
	if err != nil {
		panic("FAILED TO DROP TEST DB: " + err.Error())
	}
}

func testRowCountEquals(t *testing.T, expected int) {
	rows, err := testDb.Queryx("SELECT COUNT(*) FROM users")
	assert.NoError(t, err)

	assert.True(t, rows.Next())
	var count int
	err = rows.Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, expected, count)
}

func initTables() {
	_, _ = Init(testDb) //these must be inited in the right order
}

func cleanUpTables(t *testing.T) {
	_, err := testDb.Exec("DROP TABLE users")
	assert.NoError(t, err)
}


func TestInit(t *testing.T) {
	_, err := Init(testDb)

	assert.NoError(t, err)
	assert.NotEmpty(t, Repo)

	_, err = testDb.Exec("DROP TABLE users")
	assert.NoError(t, err)
}

func TestRepository_Create(t *testing.T) {
	initTables()

	user1 := model.User{Username: "user1", Secret: "test"}
	user2 := model.User{Username: "user2", Secret: "test"}

	validUser1, err := Repo.Create(user1)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser1)
	assert.Equal(t, user1.Username, validUser1.Username)

	testRowCountEquals(t, 1)

	validUser2, err := Repo.Create(user2)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser2)
	assert.Equal(t, user2.Username, validUser2.Username)

	testRowCountEquals(t, 2)

	assert.NotEqual(t, validUser1.Id, validUser2.Id)

	_, err = Repo.Create(user1)
	assert.Error(t, err)

	testRowCountEquals(t, 2)

	cleanUpTables(t)
}

func TestRepository_GetAll(t *testing.T) {
	initTables()

	user1 := model.User{Username: "user1", Secret: "test"}
	user2 := model.User{Username: "user2", Secret: "test"}

	validUser1, err := Repo.Create(user1)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser1)
	assert.Equal(t, user1.Username, validUser1.Username)

	validUser2, err := Repo.Create(user2)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser2)
	assert.Equal(t, user2.Username, validUser2.Username)

	allUsers, err := Repo.GetAll()
	assert.Len(t, allUsers, 2)

	cleanUpTables(t)
}

func TestRepository_GetById(t *testing.T) {
	initTables()

	user1 := model.User{Username: "user1", Secret: "test"}
	user2 := model.User{Username: "user2", Secret: "test"}

	validUser1, err := Repo.Create(user1)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser1)
	assert.Equal(t, user1.Username, validUser1.Username)

	validUser2, err := Repo.Create(user2)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser2)
	assert.Equal(t, user2.Username, validUser2.Username)

	selectedUser, err := Repo.GetById(validUser1.Id)
	assert.NoError(t, err)
	assert.Equal(t, validUser1, selectedUser)

	noUser, err := Repo.GetById(12345)
	assert.NoError(t, err)
	assert.Nil(t, noUser)

	cleanUpTables(t)
}

func TestRepository_GetByUsername(t *testing.T) {
	initTables()

	user1 := model.User{Username: "user1", Secret: "test"}
	user2 := model.User{Username: "user2", Secret: "test"}

	validUser1, err := Repo.Create(user1)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser1)
	assert.Equal(t, user1.Username, validUser1.Username)

	validUser2, err := Repo.Create(user2)
	assert.NoError(t, err)
	assert.NotEmpty(t, validUser2)
	assert.Equal(t, user2.Username, validUser2.Username)

	selectedUser, err := Repo.GetByUsername(validUser1.Username)
	assert.NoError(t, err)
	assert.Equal(t, validUser1, selectedUser)

	noUser, err := Repo.GetByUsername("not a real username")
	assert.NoError(t, err)
	assert.Nil(t, noUser)

	cleanUpTables(t)
}