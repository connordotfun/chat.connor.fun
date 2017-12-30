package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPermissionSet(t *testing.T) {
	ps := NewPermissionSet()

	assert.NotEmpty(t, ps)
}

func TestPermissionSet_Add(t *testing.T) {
	ps := NewPermissionSet()

	p1 := Permission{Path: "/foo/bar", Code: 0x00F0}
	p2 := Permission{Path: "/foo/bar/1", Code: 0x00FF}
	p3 := Permission{Path: "/foo/bar/2", Code: 0x0FF0}
	p4 := Permission{Path: "/foo/bar/3", Code: 0xF0F0}

	ps.Add(p1)

	assert.Equal(t, 1, ps.Length())

	ps.Add(p2)

	assert.Equal(t, 2, ps.Length())

	ps.Add(p3, p4)

	assert.Equal(t, 4, ps.Length())

	ps.Add(p1)
	assert.Equal(t, 4, ps.Length()) //no duplicates
}

func TestPermissionSet_Contains(t *testing.T) {
	ps := NewPermissionSet()

	p1 := Permission{Path: "/foo/bar", Code: 0x00F0}
	p2 := Permission{Path: "/foo/bar/1", Code: 0x00FF}

	assert.False(t, ps.Contains(p1))

	ps.Add(p1)

	assert.True(t, ps.Contains(p1))
	assert.False(t, ps.Contains(p2))
}

func TestPermissionSet_Permissions(t *testing.T) {
	ps := NewPermissionSet()

	p1 := Permission{Path: "/foo/bar", Code: 0x00F0}
	p2 := Permission{Path: "/foo/bar/1", Code: 0x00FF}
	p3 := Permission{Path: "/foo/bar/2", Code: 0x0FF0}
	p4 := Permission{Path: "/foo/bar/3", Code: 0xF0F0}

	ps.Add(p1, p2, p3, p4)

	permissions := ps.Permissions()

	assert.Len(t, permissions, 4)

	expectedArray := []Permission{p1, p2, p3, p4}

	for exp := range expectedArray {
		found := false
		for permission := range permissions {
			if exp == permission {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}
