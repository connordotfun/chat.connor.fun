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
		Code: 0xFFFF,
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
		Code: 0xFFFF,
	}

	assert.NoError(t, err)
	assert.Equal(t, pActual, p)
}

func TestPermission_IsPermitted_SimplePath(t *testing.T) {
	p := Permission{Path: "/foo/bar/foo/bar", Code: actionRead | actionDelete}

	assert.True(t, p.IsPermitted("GET", "/foo/bar/foo/bar"), "failed same method and path")
	assert.True(t, p.IsPermitted("DELETE", "/foo/bar/foo/bar"), "failed same method and path")
	assert.False(t, p.IsPermitted("GET", "/foo/bar/foo/baz"), "failed same method and wrong path")

	assert.False(t, p.IsPermitted("POST", "/foo/bar/foo/bar"), "failed wrong method same path")
}

func TestPermission_IsPermitted_(t *testing.T) {
	p := Permission{Path: "/foo/*/foo/bar", Code: actionRead}

	assert.True(t, p.IsPermitted("GET", "/foo/1/foo/bar"))
	assert.True(t, p.IsPermitted("GET", "/foo/2/foo/bar"))

	assert.False(t, p.IsPermitted("GET", "/foo/2/foo/abcd"))
}
