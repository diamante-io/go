package clocktest

import (
	"time"
)

// FixedSource is a clock source that has its current time stopped and fixed at
// a specific time.
type FixedSource time.Time

// Now returns the fixed source's constant time.
func (s FixedSource) Now() time.Time {
	return time.Time(s)
}
