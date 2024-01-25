package txsub

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go/support/log"

	"github.com/go-errors/errors"
)

// NewDefaultSubmissionList returns a list that manages open submissions purely
// in memory.
func NewDefaultSubmissionList() OpenSubmissionList {
	return &submissionList{
		submissions: map[string]*openSubmission{},
		log:         log.DefaultLogger.WithField("service", "txsub.submissionList"),
	}
}

// openSubmission tracks a slice of channels that should be emitted to when we
// know the result for the transactions with the provided hash
type openSubmission struct {
	Hash        string
	SubmittedAt time.Time
	Listeners   []Listener
}

type submissionList struct {
	sync.Mutex
	submissions map[string]*openSubmission // hash => `*openSubmission`
	log         *log.Entry
}

func (s *submissionList) Add(ctx context.Context, hash string, l Listener) error {
	s.Lock()
	defer s.Unlock()

	if cap(l) == 0 {
		panic("Unbuffered listener cannot be added to OpenSubmissionList")
	}

	if len(hash) != 64 {
		return errors.New("Unexpected transaction hash length: must be 64 hex characters")
	}

	os, ok := s.submissions[hash]

	if !ok {
		os = &openSubmission{
			Hash:        hash,
			SubmittedAt: time.Now(),
			Listeners:   []Listener{},
		}
		s.submissions[hash] = os
		s.log.WithField("hash", hash).Info("Created a new submission for a transaction")
	} else {
		s.log.WithField("hash", hash).Info("Adding listener to existing submission")
	}

	os.Listeners = append(os.Listeners, l)

	return nil
}

func (s *submissionList) Finish(ctx context.Context, hash string, r Result) error {
	s.Lock()
	defer s.Unlock()

	os, ok := s.submissions[hash]
	if !ok {
		return nil
	}

	s.log.WithFields(log.F{
		"hash":      hash,
		"listeners": len(os.Listeners),
		"result":    fmt.Sprintf("%+v", r),
	}).Info("Sending submission result to listeners")

	for _, l := range os.Listeners {
		l <- r
		close(l)
	}

	delete(s.submissions, hash)
	return nil
}

func (s *submissionList) Clean(ctx context.Context, maxAge time.Duration) (int, error) {
	s.Lock()
	defer s.Unlock()

	for _, os := range s.submissions {
		if time.Since(os.SubmittedAt) > maxAge {
			s.log.WithFields(log.F{
				"hash":      os.Hash,
				"listeners": len(os.Listeners),
			}).Warn("Cleared submission due to timeout")
			r := Result{Err: ErrTimeout}
			delete(s.submissions, os.Hash)
			for _, l := range os.Listeners {
				l <- r
				close(l)
			}
		}
	}

	return len(s.submissions), nil
}

func (s *submissionList) Pending(ctx context.Context) []string {
	s.Lock()
	defer s.Unlock()
	results := make([]string, 0, len(s.submissions))

	for hash := range s.submissions {
		results = append(results, hash)
	}

	return results
}
