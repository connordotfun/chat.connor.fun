package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


const (
	testJsonRoleData = `
		{
		  "admin": {
			"parent": "NONE",
			"override": "ALLOW_ALL"
		  },

		  "banned": {
			"parent": "NONE",
			"override": "ALLOW_NONE"
		  },

		  "anon_user": {
			"parent": "NONE",
			"override": "NONE",
			"permissions": [
			  {"path": "/api/v1/users",  "verbs": "c"},
			  {"path": "/api/v1/login", "verbs": "c"},
			  {"path": "/api/v1/room/*/messages",   "verbs": "r"},
			  {"path": "/api/v1/room/*/messages/*", "verbs": "r"}
			]
		  },

		  "normal_user": {
			"parent": "anon_user",
			"override": "NONE",
			"permissions": [
			  {"path": "/api/v1/users",  "verbs": "c"},
			  {"path": "/api/v1/login", "verbs": "c"},
			  {"path": "/api/v1/room/*/messages",   "verbs": "cr"},
			  {"path": "/api/v1/room/*/messages/*", "verbs": "r"}
			]
		  }
		}
	`
)

func TestInitRoleMap(t *testing.T) {
	err := InitRoleMap([]byte(testJsonRoleData))

	assert.NoError(t, err)

	assert.NotEmpty(t, Roles)

	assert.NotEmpty(t, Roles.GetRole("admin"))
	assert.NotEmpty(t, Roles.GetRole("banned"))
	assert.NotEmpty(t, Roles.GetRole("anon_user"))
	assert.NotEmpty(t, Roles.GetRole("normal_user"))

	assert.Len(t, Roles.GetRole("anon_user").Permissions, 4)
	assert.Equal(t, "ALLOW_ALL", Roles.GetRole("admin").Override)
}
