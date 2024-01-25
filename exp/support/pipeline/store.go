package pipeline

func (s *Store) init() {
	s.values = make(map[string]interface{})
}

func (s *Store) Put(name string, value interface{}) {
	s.initOnce.Do(s.init)
	s.values[name] = value
}

func (s *Store) Get(name string) interface{} {
	s.initOnce.Do(s.init)
	return s.values[name]
}
