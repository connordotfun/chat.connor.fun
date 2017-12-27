package model


type PermissionSet struct {
	data map[Permission]bool
}

func NewPermissionSet() *PermissionSet {
	return &PermissionSet{data: make(map[Permission]bool)}
}

func (s *PermissionSet) Add(p... Permission) {
	for _, perm := range p {
		s.data[perm] = true
	}
}

func (s *PermissionSet) Contains(p Permission) bool {
	_, found := s.data[p]
	return found
}

func (s *PermissionSet) Permissions() []Permission {
	keys := make([]Permission, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *PermissionSet) Length() int {
	return len(s.data)
}