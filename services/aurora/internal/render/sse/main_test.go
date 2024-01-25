package sse

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/diamcircle/go/support/test"
	"github.com/stretchr/testify/assert"
)

func TestWriteEventOutput(t *testing.T) {
	ctx, _ := test.ContextWithLogBuffer()
	testCases := []struct {
		Event             Event
		ExpectedSubstring string
	}{
		{Event{Data: "test"}, "data: \"test\"\n\n"},
		{Event{ID: "1", Data: "test"}, "id: 1\n"},
		{Event{Retry: 1000, Data: "test"}, "retry: 1000\n"},
		{Event{Error: errors.New("busted")}, "event: error\ndata: busted\n\n"},
		{Event{Event: "test", Data: "test"}, "event: test\ndata: \"test\"\n\n"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Checking for expected substring %s", tc.ExpectedSubstring), func(t *testing.T) {
			w := httptest.NewRecorder()
			WriteEvent(ctx, w, tc.Event)
			bodyString := w.Body.String()
			assert.Contains(t, bodyString, tc.ExpectedSubstring)
		})
	}
}

func TestWriteEventLogs(t *testing.T) {
	ctx, log := test.ContextWithLogBuffer()
	w := httptest.NewRecorder()
	WriteEvent(ctx, w, Event{Error: errors.New("busted")})
	assert.NotContains(t, log.String(), "level=error")
	assert.NotContains(t, log.String(), "busted")
}

// Tests that the preamble sets the correct headers and writes the hello event.
func TestWritePreamble(t *testing.T) {
	ctx, _ := test.ContextWithLogBuffer()
	w := httptest.NewRecorder()
	WritePreamble(ctx, w)
	assert.Equal(t, "text/event-stream; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "retry: 1000\nevent: open\ndata: \"hello\"\n\n")
}
