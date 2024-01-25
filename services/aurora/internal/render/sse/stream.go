package sse

import (
	"context"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/diamcircle/go/support/log"
	"github.com/diamcircle/go/support/render/problem"
)

var (
	// default error
	errBadStream = errors.New("Unexpected stream error")

	// known errors
	ErrRateLimited = errors.New("Rate limit exceeded")
)

type Stream struct {
	ctx      context.Context
	initSync sync.Once  // Variable to ensure that Init only writes the preamble once.
	mu       sync.Mutex // Mutex protects the following fields
	w        http.ResponseWriter
	done     bool
	sent     int
	limit    int
}

// NewStream creates a new stream against the provided response writer.
func NewStream(ctx context.Context, w http.ResponseWriter) *Stream {
	return &Stream{
		ctx: ctx,
		w:   w,
	}
}

// Init function is only executed once. It writes the preamble event which includes the HTTP response code and a
// hello message. This should be called before any method that writes to the client to ensure that the preamble
// has been sent first.
func (s *Stream) Init() {
	s.initSync.Do(func() {
		ok := WritePreamble(s.ctx, s.w)
		if !ok {
			s.done = true
		}
	})
}

func (s *Stream) Send(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Init()
	WriteEvent(s.ctx, s.w, e)
	s.sent++
}

func (s *Stream) SetLimit(limit int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.limit = limit
}

func (s *Stream) Done() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Init()
	WriteEvent(s.ctx, s.w, goodbyeEvent)
	s.done = true
}

func (s *Stream) Err(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If we haven't sent an event, we should simply return the normal HTTP
	// error because it means that we haven't sent the preamble.
	if s.sent == 0 {
		problem.Render(s.ctx, s.w, err)
		return
	}

	if knownErr := problem.IsKnownError(err); knownErr != nil {
		err = knownErr
	} else {
		log.Ctx(s.ctx).WithStack(err).Error(err)
		err = errBadStream
	}

	s.Init()
	WriteEvent(s.ctx, s.w, Event{Error: err})
	s.done = true
}
