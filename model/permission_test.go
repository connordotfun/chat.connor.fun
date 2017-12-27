package model

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

const (
	testPermissionStr1 = `{"path":"a/b/foo/bar","verbs":"crud"}`
)

func TestPermissionJSONMarshal(t *testing.T) {
	p := Permission{
		Path: "a/b/foo/bar",
		code: 0xFFFF,
	}

	jsonP, err := json.Marshal(p)

	assert.NoError(t, err)
	assert.Equal(t, testPermissionStr1, string(jsonP))
}

func TestPermissionJSONUnmarshal(t *testing.T) {
	var p Permission
	err := json.Unmarshal([]byte(testPermissionStr1), &p)

	pActual := Permission{
		Path: "a/b/foo/bar",
		code: 0xFFFF,
	}

	assert.NoError(t, err)
	assert.Equal(t, pActual, p)
}
