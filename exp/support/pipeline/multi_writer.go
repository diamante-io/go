package pipeline

import (
	"io"

	"github.com/diamcircle/go/support/errors"
)

func (m *multiWriter) Write(entry interface{}) error {
	m.mutex.Lock()
	m.wroteEntries++
	m.mutex.Unlock()

	results := make(chan error, len(m.writers))

	for _, w := range m.writers {
		go func(w Writer) {
			// We can keep sending entries even when io.ErrClosedPipe is returned
			// as bufferedStateReadWriter will ignore them (won't add them to
			// a channel).
			results <- w.Write(entry)
		}(w)
	}

	countClosedPipes := 0
	for range m.writers {
		err := <-results
		if err != nil {
			if err == io.ErrClosedPipe {
				countClosedPipes++
			} else {
				return err
			}
		}
	}

	// When all pipes are closed return `io.ErrClosedPipe` because there are no
	// active readers anymore.
	if countClosedPipes == len(m.writers) {
		return io.ErrClosedPipe
	}

	return nil
}

func (m *multiWriter) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.closeAfter--
	if m.closeAfter > 0 {
		return nil
	} else if m.closeAfter < 0 {
		return errors.New("Close() called more times than closeAfter")
	}

	for _, w := range m.writers {
		err := w.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

var _ Writer = &multiWriter{}
