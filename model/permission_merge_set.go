package model


type PermissionSet struct {
	data map[string]Permission
}

func NewPermissionSet() *PermissionSet {
	return &PermissionSet{data: make(map[string]Permission)}
}

func (s *PermissionSet) Add(p Permission) {
	if inMap, ok := s.data[p.Path]; ok {
		s.data[p.Path], _  = inMap.Merge(p)
	} else {
		s.data[p.Path] = p
	}
}

func (s *PermissionSet) AddAll(p... Permission) {
	for _, perm := range p {
		s.Add(perm)
	}
}

func (s *PermissionSet) Contains(p Permission) bool {
	perm, found := s.data[p.Path]
	return found && perm == p
}

func (s *PermissionSet) Permissions() []Permission {
	keys := make([]Permission, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, s.data[k])
	}
	return keys
}

func (s *PermissionSet) Length() int {
	return len(s.data)
}
