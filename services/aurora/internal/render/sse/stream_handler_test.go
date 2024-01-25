package sse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go/services/aurora/internal/ledger"
)

type testingFactory struct {
	ledgerSource ledger.Source
}

func (f *testingFactory) Get() ledger.Source {
	return f.ledgerSource
}

func TestSendByeByeOnContextDone(t *testing.T) {
	ledgerSource := ledger.NewTestingSource(1)
	handler := StreamHandler{LedgerSourceFactory: &testingFactory{ledgerSource}}

	r, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	r = r.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeStream(w, r, 10, func() ([]Event, error) {
		cancel()
		return []Event{}, nil
	})

	expected := "retry: 1000\nevent: open\ndata: \"hello\"\n\n" +
		"retry: 10\nevent: close\ndata: \"byebye\"\n\n"

	if got := w.Body.String(); got != expected {
		t.Fatalf("expected '%v' but got '%v'", expected, got)
	}
}
