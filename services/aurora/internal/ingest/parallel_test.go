package ingest

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/diamcircle/go/services/aurora/internal/db2/history"
	"github.com/diamcircle/go/support/errors"
)

func TestCalculateParallelLedgerBatchSize(t *testing.T) {
	assert.Equal(t, uint32(6656), calculateParallelLedgerBatchSize(20096, 20096, 3))
	assert.Equal(t, uint32(4992), calculateParallelLedgerBatchSize(20096, 20096, 4))
	assert.Equal(t, uint32(4992), calculateParallelLedgerBatchSize(20096, 0, 4))
	assert.Equal(t, uint32(64), calculateParallelLedgerBatchSize(64, 256, 4))
	assert.Equal(t, uint32(64), calculateParallelLedgerBatchSize(64, 32, 4))
	assert.Equal(t, uint32(64), calculateParallelLedgerBatchSize(2, 256, 4))
	assert.Equal(t, uint32(64), calculateParallelLedgerBatchSize(20096, 64, 1))
}

func TestParallelReingestRange(t *testing.T) {
	config := Config{}
	var (
		rangesCalled []history.LedgerRange
		m            sync.Mutex
	)
	result := &mockSystem{}
	result.On("ReingestRange", mock.AnythingOfType("[]history.LedgerRange"), mock.AnythingOfType("bool")).Run(
		func(args mock.Arguments) {
			m.Lock()
			defer m.Unlock()
			rangesCalled = append(rangesCalled, args.Get(0).([]history.LedgerRange)...)
			// simulate call
			time.Sleep(time.Millisecond * time.Duration(10+rand.Int31n(50)))
		}).Return(error(nil))
	factory := func(c Config) (System, error) {
		return result, nil
	}
	system, err := newParallelSystems(config, 3, factory)
	assert.NoError(t, err)
	err = system.ReingestRange([]history.LedgerRange{{1, 2050}}, 258)
	assert.NoError(t, err)

	sort.Slice(rangesCalled, func(i, j int) bool {
		return rangesCalled[i].StartSequence < rangesCalled[j].StartSequence
	})
	expected := []history.LedgerRange{
		{StartSequence: 1, EndSequence: 256}, {StartSequence: 257, EndSequence: 512}, {StartSequence: 513, EndSequence: 768}, {StartSequence: 769, EndSequence: 1024}, {StartSequence: 1025, EndSequence: 1280},
		{StartSequence: 1281, EndSequence: 1536}, {StartSequence: 1537, EndSequence: 1792}, {StartSequence: 1793, EndSequence: 2048}, {StartSequence: 2049, EndSequence: 2050},
	}
	assert.Equal(t, expected, rangesCalled)

	rangesCalled = nil
	system, err = newParallelSystems(config, 1, factory)
	assert.NoError(t, err)
	err = system.ReingestRange([]history.LedgerRange{{1, 1024}}, 64)
	result.AssertExpectations(t)
	expected = []history.LedgerRange{
		{StartSequence: 1, EndSequence: 64}, {StartSequence: 65, EndSequence: 128}, {StartSequence: 129, EndSequence: 192}, {StartSequence: 193, EndSequence: 256}, {StartSequence: 257, EndSequence: 320},
		{StartSequence: 321, EndSequence: 384}, {StartSequence: 385, EndSequence: 448}, {StartSequence: 449, EndSequence: 512}, {StartSequence: 513, EndSequence: 576}, {StartSequence: 577, EndSequence: 640},
		{StartSequence: 641, EndSequence: 704}, {StartSequence: 705, EndSequence: 768}, {StartSequence: 769, EndSequence: 832}, {StartSequence: 833, EndSequence: 896}, {StartSequence: 897, EndSequence: 960},
		{StartSequence: 961, EndSequence: 1024},
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, rangesCalled)
}

func TestParallelReingestRangeError(t *testing.T) {
	config := Config{}
	result := &mockSystem{}
	// Fail on the second range
	result.On("ReingestRange", []history.LedgerRange{{1537, 1792}}, mock.AnythingOfType("bool")).Return(errors.New("failed because of foo"))
	result.On("ReingestRange", mock.AnythingOfType("[]history.LedgerRange"), mock.AnythingOfType("bool")).Return(error(nil))
	factory := func(c Config) (System, error) {
		return result, nil
	}
	system, err := newParallelSystems(config, 3, factory)
	assert.NoError(t, err)
	err = system.ReingestRange([]history.LedgerRange{{1, 2050}}, 258)
	result.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, "job failed, recommended restart range: [1537, 2050]: error when processing [1537, 1792] range: failed because of foo", err.Error())
}

func TestParallelReingestRangeErrorInEarlierJob(t *testing.T) {
	config := Config{}
	var wg sync.WaitGroup
	wg.Add(1)
	result := &mockSystem{}
	// Fail on an lower subrange after the first error
	result.On("ReingestRange", []history.LedgerRange{{1025, 1280}}, mock.AnythingOfType("bool")).Run(func(mock.Arguments) {
		// Wait for a more recent range to error
		wg.Wait()
		// This sleep should help making sure the result of this range is processed later than the one below
		// (there are no guarantees without instrumenting ReingestRange(), but that's too complicated)
		time.Sleep(50 * time.Millisecond)
	}).Return(errors.New("failed because of foo"))
	result.On("ReingestRange", []history.LedgerRange{{1537, 1792}}, mock.AnythingOfType("bool")).Run(func(mock.Arguments) {
		wg.Done()
	}).Return(errors.New("failed because of bar"))
	result.On("ReingestRange", mock.AnythingOfType("[]history.LedgerRange"), mock.AnythingOfType("bool")).Return(error(nil))

	factory := func(c Config) (System, error) {
		return result, nil
	}
	system, err := newParallelSystems(config, 3, factory)
	assert.NoError(t, err)
	err = system.ReingestRange([]history.LedgerRange{{1, 2050}}, 258)
	result.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, "job failed, recommended restart range: [1025, 2050]: error when processing [1025, 1280] range: failed because of foo", err.Error())

}
