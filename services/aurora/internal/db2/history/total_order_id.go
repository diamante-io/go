package history

import (
	"fmt"
)

// PagingToken returns a cursor for this record
func (r *TotalOrderID) PagingToken() string {
	return fmt.Sprintf("%d", r.ID)
}
