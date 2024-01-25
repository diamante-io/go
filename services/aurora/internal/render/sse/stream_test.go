//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package sse

import (
	"context"
	"database/sql"
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	hProblem "go/services/aurora/internal/render/problem"
	"go/support/render/problem"
	"go/support/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (s *Stream) SentCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sent
}

// IsDone is safe to call concurrently and is exported.
func (s *Stream) IsDone() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.limit == 0 {
		return s.done
	}

	return s.done || s.sent >= s.limit
}

type StreamTestSuite struct {
	suite.Suite
	ctx    context.Context
	w      *httptest.ResponseRecorder
	stream *Stream
}

// Helper method to check that the preamble has been sent and all HTTP response headers are correctly set.
func (suite *StreamTestSuite) checkHeadersAndPreamble() {
	if suite.stream.SentCount() == 0 {
		assert.Equal(suite.T(), "application/problem+json; charset=utf-8", suite.w.Header().Get("Content-Type"))
		assert.Equal(suite.T(), 500, suite.w.Code)
		return
	}

	assert.Equal(suite.T(), "text/event-stream; charset=utf-8", suite.w.Header().Get("Content-Type"))
	assert.Equal(suite.T(), "no-cache", suite.w.Header().Get("Cache-Control"))
	assert.Equal(suite.T(), "keep-alive", suite.w.Header().Get("Connection"))
	assert.Equal(suite.T(), "*", suite.w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), 200, suite.w.Code)
	assert.Contains(suite.T(), suite.w.Body.String(), "retry: 1000\nevent: open\ndata: \"hello\"\n\n")
}

func (suite *StreamTestSuite) SetupTest() {
	suite.ctx, _ = test.ContextWithLogBuffer()
	suite.w = httptest.NewRecorder()
	suite.stream = NewStream(suite.ctx, suite.w)
}

// Tests that the stream sends the preamble before any events and that events are correctly sent.
func (suite *StreamTestSuite) TestStream_Send() {
	e := Event{Data: "test message"}
	suite.stream.Send(e)
	// Before sending, it should have sent the preamble first and set the headers.
	suite.checkHeadersAndPreamble()
	// Now check that the data got written
	assert.Contains(suite.T(), suite.w.Body.String(), "data: \"test message\"\n\n")
	suite.stream.Done()
	assert.Equal(suite.T(), 1, suite.stream.SentCount())
}

// Tests that Stream can send error events.
func (suite *StreamTestSuite) TestStream_Err() {
	err := errors.New("example error")
	// If we encounter an error before sending any event, we should just
	// return the error without the hello message.
	suite.stream.Err(err)
	suite.checkHeadersAndPreamble()

	// Reset the stream to test the scenario where an event has been sent.
	suite.w = httptest.NewRecorder()
	suite.stream = NewStream(suite.ctx, suite.w)
	suite.stream.sent++
	suite.stream.Err(err)
	suite.checkHeadersAndPreamble()
	assert.Contains(suite.T(), suite.w.Body.String(), "event: error\ndata: Unexpected stream error\n\n")
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that Stream can send handled registered errors
func (suite *StreamTestSuite) TestStream_ErrRegisterError() {
	problem.RegisterError(context.DeadlineExceeded, hProblem.Timeout)
	defer problem.UnRegisterErrors()

	suite.w = httptest.NewRecorder()
	suite.stream = NewStream(suite.ctx, suite.w)
	suite.stream.sent++
	suite.stream.Err(context.DeadlineExceeded)
	suite.checkHeadersAndPreamble()
	assert.Contains(suite.T(), suite.w.Body.String(), "event: error\ndata: problem: timeout\n\n")
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that Stream can send handled ErrNoRows
func (suite *StreamTestSuite) TestStream_ErrNoRows() {
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
	defer problem.UnRegisterErrors()

	suite.w = httptest.NewRecorder()
	suite.stream = NewStream(suite.ctx, suite.w)
	suite.stream.sent++
	suite.stream.Err(sql.ErrNoRows)
	suite.checkHeadersAndPreamble()
	assert.Contains(suite.T(), suite.w.Body.String(), "event: error\ndata: problem: not_found\n\n")
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that SetLimit sets stream.done to true after the limit has been reached.
func (suite *StreamTestSuite) TestStream_SetLimit() {
	suite.stream.SetLimit(3)
	// Send more than the limit
	for i := 0; i < 5; i++ {
		message := "test message " + strconv.Itoa(i)
		suite.stream.Send(Event{Data: message})
	}
	assert.True(suite.T(), suite.stream.IsDone())
}

// Tests that SentCount reports the correct number.
func (suite *StreamTestSuite) TestStream_SentCount() {
	for i := 0; i < 5; i++ {
		message := "test message " + strconv.Itoa(i)
		suite.stream.Send(Event{Data: message})
	}
	assert.Equal(suite.T(), 5, suite.stream.SentCount())
	suite.stream.Err(errors.New("example error"))
	// Make sure that errors don't contribute to the send count
	assert.Equal(suite.T(), 5, suite.stream.SentCount())
}

// Runs the test suite.
func TestStreamTestSuite(t *testing.T) {
	suite.Run(t, new(StreamTestSuite))
}
