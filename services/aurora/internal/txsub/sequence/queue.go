package sequence

import (
	"container/heap"
	"time"
)

// Queue manages the submission queue for a single source account. The
// transaction system uses Push to enqueue submissions for given sequence
// numbers.
//
// Queue maintains a priority queue of pending submissions, and when updated
// (via the Update() method) with the current sequence number of the account
// being managed, queued submissions that can be acted upon will be unblocked.
//
type Queue struct {
	lastActiveAt time.Time
	timeout      time.Duration
	nextSequence uint64
	queue        pqueue
}

// NewQueue creates a new *Queue
func NewQueue() *Queue {
	result := &Queue{
		lastActiveAt: time.Now(),
		timeout:      10 * time.Second,
		queue:        nil,
	}

	heap.Init(&result.queue)

	return result
}

// Size returns the count of currently buffered submissions in the queue.
func (q *Queue) Size() int {
	return len(q.queue)
}

// Push enqueues the intent to submit a transaction at the provided sequence
// number and returns a channel that will emit when it is safe for the client
// to do so.
//
// Push does not perform any triggering (which
// occurs in Update(), even if the current sequence number for this queue is
// the same as the provided sequence, to keep internal complexity much lower.
// Given that, the recommended usage pattern is:
//
// 1. Push the submission onto the queue
// 2. Load the current sequence number for the source account from the DB
// 3. Call Update() with the result from step 2 to trigger the submission if
//		possible
func (q *Queue) Push(sequence uint64) <-chan error {
	ch := make(chan error, 1)
	heap.Push(&q.queue, item{sequence, ch})
	return ch
}

// Update notifies the queue that the provided sequence number is the latest
// seen value for the account that this queue manages submissions for.
//
// This function is monotonic... calling it with a sequence number lower than
// the latest seen sequence number is a noop.
func (q *Queue) Update(sequence uint64) {
	if q.nextSequence <= sequence {
		q.nextSequence = sequence + 1
	}

	wasChanged := false

	for {
		if q.Size() == 0 {
			break
		}

		ch, hseq := q.head()
		// if the next queued transaction has a sequence higher than the account's
		// current sequence, stop removing entries
		if hseq > q.nextSequence {
			break
		}

		// since this entry is unlocked (i.e. it's sequence is the next available
		// or in the past we can remove it an mark the queue as changed
		q.pop()
		wasChanged = true

		if hseq < q.nextSequence {
			ch <- ErrBadSequence
			close(ch)
		} else if hseq == q.nextSequence {
			ch <- nil
			close(ch)
		}
	}

	// if we modified the queue, bump the timeout for this queue
	if wasChanged {
		q.lastActiveAt = time.Now()
		return
	}

	// if the queue wasn't changed, see if it is too old, clear
	// it and make room for other's
	if time.Since(q.lastActiveAt) > q.timeout {
		for q.Size() > 0 {
			ch, _ := q.pop()
			ch <- ErrBadSequence
			close(ch)
		}
	}
}

// helper function for interacting with the priority queue
func (q *Queue) head() (chan error, uint64) {
	if len(q.queue) == 0 {
		return nil, uint64(0)
	}

	return q.queue[0].Chan, q.queue[0].Sequence
}

// helper function for interacting with the priority queue
func (q *Queue) pop() (chan error, uint64) {
	i := heap.Pop(&q.queue).(item)

	return i.Chan, i.Sequence
}

// item is a element of the priority queue
type item struct {
	Sequence uint64
	Chan     chan error
}

// pqueue is a priority queue used by Queue to manage buffered submissions.  It
// implements heap.Interface.
type pqueue []item

func (pq pqueue) Len() int { return len(pq) }

func (pq pqueue) Less(i, j int) bool {
	return pq[i].Sequence < pq[j].Sequence
}

func (pq pqueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *pqueue) Push(x interface{}) {
	*pq = append(*pq, x.(item))
}

func (pq *pqueue) Pop() interface{} {
	old := *pq
	n := len(old)
	result := old[n-1]
	*pq = old[0 : n-1]
	return result
}
