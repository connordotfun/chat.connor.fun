package roles

import (
	"testing"
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
	"github.com/aaronaaeng/chat.connor.fun/db/users"
	"fmt"
)

const (
	testDbHost = "localhost"
	testDbPort = 5432
	testDbName = "connorfuntest"
)

var testDb sqlx.DB

func TestMain(m *testing.M) {
	//Init the connorfuntest db and reconnect to the new db
	db, err := sqlx.Open("postgres", "postgresql://localhost:5432?sslmode=disable")
	if err != nil {
		panic("failed to establish db connection")
	}

	_, err = db.Exec("CREATE DATABASE connorfuntest;")
	if err != nil {
		panic("Failed to create test db: " + err.Error())
	}

	testDb, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s sslmode=disable", testDbHost, testDbPort, testDbName))
	if err != nil {
		panic("failed to establish db connection")
	}

	//init the user table
	_, err = users.Init(testDb)
	if err != nil {
		panic(err)
	}
	m.Run()

	//clean up the connorfuntestdb
	testDb.Close()
	_, err = db.Exec("DROP DATABASE connorfuntest;")
	if err != nil {
		panic("FAILED TO DROP TEST DB: " + err.Error())
	}
	db.Close()
}

func TestInit(t *testing.T) {

}

func TestRepository_AddRole(t *testing.T) {

}

func TestRepository_GetUserRoles(t *testing.T) {

}
