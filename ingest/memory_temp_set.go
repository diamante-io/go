package ingest

// memoryTempSet is an in-memory implementation of TempSet interface.
// As of July 2019 this requires up to ~4GB of memory for pubnet ledger
// state processing. The internal structure is dereferenced after the
// store is closed.
type memoryTempSet struct {
	m map[string]bool
}

// Open initialize internals data structure.
func (s *memoryTempSet) Open() error {
	s.m = make(map[string]bool)
	return nil
}

// Add adds a key to TempSet.
func (s *memoryTempSet) Add(key string) error {
	s.m[key] = true
	return nil
}

// Preload does not do anything. This TempSet keeps everything in memory
// so no preloading needed.
func (s *memoryTempSet) Preload(keys []string) error {
	return nil
}

// Exist check if the key exists in a TempSet.
func (s *memoryTempSet) Exist(key string) (bool, error) {
	return s.m[key], nil
}

// Close removes reference to internal data structure.
func (s *memoryTempSet) Close() error {
	s.m = nil
	return nil
}
