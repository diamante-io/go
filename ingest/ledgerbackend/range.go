package ledgerbackend

import (
	"encoding/json"
	"fmt"
)

// Range represents a range of ledger sequence numbers.
type Range struct {
	from    uint32
	to      uint32
	bounded bool
}

type jsonRange struct {
	From    uint32 `json:"from"`
	To      uint32 `json:"to"`
	Bounded bool   `json:"bounded"`
}

func (r *Range) UnmarshalJSON(b []byte) error {
	var s jsonRange
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	r.from = s.From
	r.to = s.To
	r.bounded = s.Bounded

	return nil
}

func (r Range) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonRange{
		From:    r.from,
		To:      r.to,
		Bounded: r.bounded,
	})
}

func (r Range) String() string {
	if r.bounded {
		return fmt.Sprintf("[%d,%d]", r.from, r.to)
	}
	return fmt.Sprintf("[%d,latest)", r.from)
}

func (r Range) Contains(other Range) bool {
	if r.bounded && !other.bounded {
		return false
	}
	if r.bounded && other.bounded {
		return r.from <= other.from && r.to >= other.to
	}
	return r.from <= other.from
}

// SingleLedgerRange constructs a bounded range containing a single ledger.
func SingleLedgerRange(ledger uint32) Range {
	return Range{from: ledger, to: ledger, bounded: true}
}

// BoundedRange constructs a bounded range of ledgers with a fixed starting ledger and ending ledger.
func BoundedRange(from uint32, to uint32) Range {
	return Range{from: from, to: to, bounded: true}
}

// BoundedRange constructs a unbounded range of ledgers with a fixed starting ledger.
func UnboundedRange(from uint32) Range {
	return Range{from: from, bounded: false}
}
