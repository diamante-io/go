package ledgerbackend

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"

	"go/support/errors"
	"go/xdr"
)

// PrepareRangeResponse describes the status of the pending PrepareRange operation.
type PrepareRangeResponse struct {
	LedgerRange   Range     `json:"ledgerRange"`
	StartTime     time.Time `json:"startTime"`
	Ready         bool      `json:"ready"`
	ReadyDuration int       `json:"readyDuration"`
}

// LatestLedgerSequenceResponse is the response for the GetLatestLedgerSequence command.
type LatestLedgerSequenceResponse struct {
	Sequence uint32 `json:"sequence"`
}

// LedgerResponse is the response for the GetLedger command.
type LedgerResponse struct {
	Ledger Base64Ledger `json:"ledger"`
}

// Base64Ledger extends xdr.LedgerCloseMeta with JSON encoding and decoding
type Base64Ledger xdr.LedgerCloseMeta

func (r *Base64Ledger) UnmarshalJSON(b []byte) error {
	var base64 string
	if err := json.Unmarshal(b, &base64); err != nil {
		return err
	}

	var parsed xdr.LedgerCloseMeta
	if err := xdr.SafeUnmarshalBase64(base64, &parsed); err != nil {
		return err
	}
	*r = Base64Ledger(parsed)

	return nil
}

func (r Base64Ledger) MarshalJSON() ([]byte, error) {
	base64, err := xdr.MarshalBase64(xdr.LedgerCloseMeta(r))
	if err != nil {
		return nil, err
	}
	return json.Marshal(base64)
}

// RemoteCaptivediamcircleCore is an http client for interacting with a remote captive core server.
type RemoteCaptivediamcircleCore struct {
	url                      *url.URL
	client                   *http.Client
	lock                     *sync.Mutex
	prepareRangePollInterval time.Duration
}

// RemoteCaptiveOption values can be passed into NewRemoteCaptive to customize a RemoteCaptivediamcircleCore instance.
type RemoteCaptiveOption func(c *RemoteCaptivediamcircleCore)

// PrepareRangePollInterval configures how often the captive core server will be polled when blocking
// on the PrepareRange operation.
func PrepareRangePollInterval(d time.Duration) RemoteCaptiveOption {
	return func(c *RemoteCaptivediamcircleCore) {
		c.prepareRangePollInterval = d
	}
}

// NewRemoteCaptive returns a new RemoteCaptivediamcircleCore instance.
//
// Only the captiveCoreURL parameter is required.
func NewRemoteCaptive(captiveCoreURL string, options ...RemoteCaptiveOption) (RemoteCaptivediamcircleCore, error) {
	u, err := url.Parse(captiveCoreURL)
	if err != nil {
		return RemoteCaptivediamcircleCore{}, errors.Wrap(err, "unparseable url")
	}

	client := RemoteCaptivediamcircleCore{
		prepareRangePollInterval: time.Second,
		url:                      u,
		client:                   &http.Client{Timeout: 10 * time.Second},
		lock:                     &sync.Mutex{},
	}
	for _, option := range options {
		option(&client)
	}
	return client, nil
}

func decodeResponse(response *http.Response, payload interface{}) error {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.New(string(body))
	}

	if err := json.NewDecoder(response.Body).Decode(payload); err != nil {
		return errors.Wrap(err, "failed to decode json payload")
	}
	return nil
}

// GetLatestLedgerSequence returns the sequence of the latest ledger available
// in the backend. This method returns an error if not in a session (start with
// PrepareRange).
//
// Note that for UnboundedRange the returned sequence number is not necessarily
// the latest sequence closed by the network. It's always the last value available
// in the backend.
func (c RemoteCaptivediamcircleCore) GetLatestLedgerSequence(ctx context.Context) (sequence uint32, err error) {
	// TODO: Have a context on this request so we can cancel all outstanding
	// requests, not just PrepareRange.
	u := *c.url
	u.Path = path.Join(u.Path, "latest-sequence")
	request, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return 0, errors.Wrap(err, "cannot construct http request")
	}

	response, err := c.client.Do(request)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute request")
	}

	var parsed LatestLedgerSequenceResponse
	if err = decodeResponse(response, &parsed); err != nil {
		return 0, err
	}

	return parsed.Sequence, nil
}

// Close cancels any pending PrepareRange requests.
func (c RemoteCaptivediamcircleCore) Close() error {
	return nil
}

// PrepareRange prepares the given range (including from and to) to be loaded.
// Captive diamcircle-core backend needs to initalize diamcircle-Core state to be
// able to stream ledgers.
// diamcircle-Core mode depends on the provided ledgerRange:
//   * For BoundedRange it will start diamcircle-Core in catchup mode.
//   * For UnboundedRange it will first catchup to starting ledger and then run
//     it normally (including connecting to the diamcircle network).
// Please note that using a BoundedRange, currently, requires a full-trust on
// history archive. This issue is being fixed in diamcircle-Core.
func (c RemoteCaptivediamcircleCore) PrepareRange(ctx context.Context, ledgerRange Range) error {
	// TODO: removing createContext call here means we could technically have
	// multiple prepareRange requests happening at the same time. Do we still
	// need to enforce that?

	timer := time.NewTimer(c.prepareRangePollInterval)
	defer timer.Stop()

	for {
		ready, err := c.IsPrepared(ctx, ledgerRange)
		if err != nil {
			return err
		}
		if ready {
			return nil
		}

		select {
		case <-ctx.Done():
			return errors.Wrap(ctx.Err(), "shutting down")
		case <-timer.C:
			timer.Reset(c.prepareRangePollInterval)
		}
	}
}

// IsPrepared returns true if a given ledgerRange is prepared.
func (c RemoteCaptivediamcircleCore) IsPrepared(ctx context.Context, ledgerRange Range) (bool, error) {
	// TODO: Have some way to cancel all outstanding requests, not just
	// PrepareRange.
	u := *c.url
	u.Path = path.Join(u.Path, "prepare-range")
	rangeBytes, err := json.Marshal(ledgerRange)
	if err != nil {
		return false, errors.Wrap(err, "cannot serialize range")
	}
	body := bytes.NewReader(rangeBytes)
	request, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return false, errors.Wrap(err, "cannot construct http request")
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

	var response *http.Response
	response, err = c.client.Do(request)
	if err != nil {
		return false, errors.Wrap(err, "failed to execute request")
	}

	var parsed PrepareRangeResponse
	if err = decodeResponse(response, &parsed); err != nil {
		return false, err
	}

	return parsed.Ready, nil
}

// GetLedger long-polls a remote diamcircle core backend, until the requested
// ledger is ready.

// Call PrepareRange first to instruct the backend which ledgers to fetch.
//
// Requesting a ledger on non-prepared backend will return an error.
//
// Because data is streamed from diamcircle-Core ledger after ledger user should
// request sequences in a non-decreasing order. If the requested sequence number
// is less than the last requested sequence number, an error will be returned.
func (c RemoteCaptivediamcircleCore) GetLedger(ctx context.Context, sequence uint32) (xdr.LedgerCloseMeta, error) {
	for {
		// TODO: Have some way to cancel all outstanding requests, not just
		// PrepareRange.
		u := *c.url
		u.Path = path.Join(u.Path, "ledger", strconv.FormatUint(uint64(sequence), 10))
		request, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
		if err != nil {
			return xdr.LedgerCloseMeta{}, errors.Wrap(err, "cannot construct http request")
		}

		response, err := c.client.Do(request)
		if err != nil {
			return xdr.LedgerCloseMeta{}, errors.Wrap(err, "failed to execute request")
		}

		if response.StatusCode == http.StatusRequestTimeout {
			response.Body.Close()
			// This request timed out. Retry.
			continue
		}

		var parsed LedgerResponse
		if err = decodeResponse(response, &parsed); err != nil {
			return xdr.LedgerCloseMeta{}, err
		}

		return xdr.LedgerCloseMeta(parsed.Ledger), nil
	}
}
