package sequence

import (
	"fmt"
	"strings"
	"sync"
)

// Manager provides a system for tracking the transaction submission queue for
// a set of addresses.  Requests to submit at a certain sequence number are
// registered using the Push() method, and as the system is updated with
// account sequence information (through the Update() method) requests are
// notified that they can safely submit to diamcircle-core.
type Manager struct {
	mutex   sync.Mutex
	MaxSize int
	queues  map[string]*Queue
}

// NewManager returns a new manager
func NewManager() *Manager {
	return &Manager{
		MaxSize: 1024, //TODO: make MaxSize configurable
		queues:  map[string]*Queue{},
	}
}

func (m *Manager) String() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var addys []string

	for addy, q := range m.queues {
		addys = append(addys, fmt.Sprintf("%5s:%d", addy, q.nextSequence))
	}

	return "[ " + strings.Join(addys, ",") + " ]"
}

// Size returns the count of submissions buffered within
// this manager.
func (m *Manager) Size() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.size()
}

func (m *Manager) Addresses() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	addys := make([]string, 0, len(m.queues))

	for addy := range m.queues {
		addys = append(addys, addy)
	}

	return addys
}

// Push registers an intent to submit a transaction for the provided address at
// the provided sequence.  A channel is returned that will be written to when
// the requester should attempt the submission.
func (m *Manager) Push(address string, sequence uint64) <-chan error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.size() >= m.MaxSize {
		return m.getError(ErrNoMoreRoom)
	}

	aq, ok := m.queues[address]
	if !ok {
		aq = NewQueue()
		m.queues[address] = aq
	}

	return aq.Push(sequence)
}

// Update notifies the manager of newly loaded account sequence information.  The manager uses this information
// to notify requests to submit that they should proceed.  See Queue#Update for the actual meat of the logic.
func (m *Manager) Update(updates map[string]uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for address, seq := range updates {
		queue, ok := m.queues[address]
		if !ok {
			continue
		}

		queue.Update(seq)
		if queue.Size() == 0 {
			delete(m.queues, address)
		}
	}
}

// size returns the count of submissions buffered within this manager.  This
// internal version assumes you have locked the manager previously.
func (m *Manager) size() int {
	var result int
	for _, q := range m.queues {
		result += q.Size()
	}

	return result
}

func (m *Manager) getError(err error) <-chan error {
	ch := make(chan error, 1)
	ch <- err
	close(ch)
	return ch
}
