package clock

import "time"

// Clock provides access to the current time from an underlying source.
type Clock struct {
	Source Source
}

func (c *Clock) getSource() Source {
	if c == nil || c.Source == nil {
		return RealSource{}
	}
	return c.Source
}

// Now returns the current time as defined by the Clock's Source.
func (c *Clock) Now() time.Time {
	return c.getSource().Now()
}

// Source is any type providing a Now function that returns the current time.
type Source interface {
	// Now returns the current time.
	Now() time.Time
}

// RealSource is a Source that uses the real time as provided by the stdlib
// time.Now() function as the current time.
type RealSource struct{}

// Now returns the real system time as reported by time.Now().
func (RealSource) Now() time.Time {
	return time.Now()
}
