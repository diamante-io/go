package hal

// Pageable implementors can be added to hal.Page collections
type Pageable interface {
	PagingToken() string
}
