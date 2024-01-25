//lint:file-ignore U1001 Ignore all unused code, staticcheck doesn't understand testify/suite

package sequence

import (
	"context"
	"testing"
	"time"

	"go/services/aurora/internal/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type QueueTestSuite struct {
	suite.Suite
	ctx   context.Context
	queue *Queue
}

func (suite *QueueTestSuite) SetupTest() {
	suite.ctx = test.Context()
	suite.queue = NewQueue()
}

//Push adds the provided channel on to the priority queue
func (suite *QueueTestSuite) TestQueue_Push() {
	ctx := test.Context()
	_ = ctx

	assert.Equal(suite.T(), 0, suite.queue.Size())

	suite.queue.Push(2)
	assert.Equal(suite.T(), 1, suite.queue.Size())
	_, s := suite.queue.head()
	assert.Equal(suite.T(), uint64(2), s)

	suite.queue.Push(1)
	assert.Equal(suite.T(), 2, suite.queue.Size())
	_, s = suite.queue.head()
	assert.Equal(suite.T(), uint64(1), s)
}

// Tests the update method
func (suite *QueueTestSuite) TestQueue_Update() {
	// Update removes sequences that are submittable or in the past
	results := []<-chan error{
		suite.queue.Push(1),
		suite.queue.Push(2),
		suite.queue.Push(3),
		suite.queue.Push(4),
	}

	suite.queue.Update(2)

	// the update above signifies that 2 is the accounts current sequence,
	// meaning that 3 is submittable, and so only 4 should still be queued
	assert.Equal(suite.T(), 1, suite.queue.Size())
	_, s := suite.queue.head()
	assert.Equal(suite.T(), uint64(4), s)

	suite.queue.Update(4)
	assert.Equal(suite.T(), 0, suite.queue.Size())

	assert.Equal(suite.T(), ErrBadSequence, <-results[0])
	assert.Equal(suite.T(), ErrBadSequence, <-results[1])
	assert.Equal(suite.T(), nil, <-results[2])
	assert.Equal(suite.T(), ErrBadSequence, <-results[3])

	// Update clears the queue if the head has not been released within the time limit
	suite.queue.timeout = 1 * time.Millisecond
	result := suite.queue.Push(2)
	<-time.After(10 * time.Millisecond)
	suite.queue.Update(0)

	assert.Equal(suite.T(), 0, suite.queue.Size())
	assert.Equal(suite.T(), ErrBadSequence, <-result)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
