package internal

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go/ingest/ledgerbackend"
	supporthttp "go/support/http"
	"go/support/http/httpdecode"
	supportlog "go/support/log"
)

func serializeResponse(
	logger *supportlog.Entry,
	w http.ResponseWriter,
	r *http.Request,
	response interface{},
	err error,
) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.WithContext(r.Context()).WithError(err).Warn("could not serialize response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type GetLedgerRequest struct {
	Sequence uint32 `path:"sequence"`
}

// Handler returns an HTTP handler which exposes captive core operations via HTTP endpoints.
func Handler(api CaptiveCoreAPI) http.Handler {
	mux := supporthttp.NewMux(api.log)

	mux.Get("/latest-sequence", func(w http.ResponseWriter, r *http.Request) {
		response, err := api.GetLatestLedgerSequence(r.Context())
		serializeResponse(api.log, w, r, response, err)
	})

	mux.Get("/ledger/{sequence}", func(w http.ResponseWriter, r *http.Request) {
		req := GetLedgerRequest{}
		if err := httpdecode.Decode(r, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// must be shorter than the RemoteCaptiveCore http client timeout.
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel() // release timer

		var response ledgerbackend.LedgerResponse
		var err error
		done := make(chan struct{})
		go func() {
			response, err = api.GetLedger(ctx, req.Sequence)
			close(done)
		}()

		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusRequestTimeout)
		case <-done:
			serializeResponse(api.log, w, r, response, err)
		}
	})

	mux.Post("/prepare-range", func(w http.ResponseWriter, r *http.Request) {
		ledgerRange := ledgerbackend.Range{}
		if err := json.NewDecoder(r.Body).Decode(&ledgerRange); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		response, err := api.PrepareRange(r.Context(), ledgerRange)
		serializeResponse(api.log, w, r, response, err)
	})

	return mux
}
