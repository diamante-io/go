package diamcirclecore

import (
	"context"
	"net/http"
	"testing"

	proto "go/protocols/diamcirclecore"
	"go/support/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestSubmitTransaction(t *testing.T) {
	hmock := httptest.NewClient()
	c := &Client{HTTP: hmock, URL: "http://localhost:11626"}

	// happy path - new transaction
	hmock.On("GET", "http://localhost:11626/tx?blob=foo").
		ReturnJSON(http.StatusOK, proto.TXResponse{
			Status: proto.TXStatusPending,
		})

	resp, err := c.SubmitTransaction(context.Background(), "foo")

	if assert.NoError(t, err) {
		assert.Equal(t, proto.TXStatusPending, resp.Status)
	}
}

func TestManualClose(t *testing.T) {
	hmock := httptest.NewClient()
	c := &Client{HTTP: hmock, URL: "http://localhost:11626"}

	// happy path - new transaction
	hmock.On("GET", "http://localhost:11626/manualclose").
		ReturnString(http.StatusOK, "Manually triggered a ledger close with sequence number 7")

	err := c.ManualClose(context.Background())

	assert.NoError(t, err)
}

func TestManualClose_NotAvailable(t *testing.T) {
	hmock := httptest.NewClient()
	c := &Client{HTTP: hmock, URL: "http://localhost:11626"}

	// happy path - new transaction
	hmock.On("GET", "http://localhost:11626/manualclose").
		ReturnString(http.StatusOK,
			`{"exception": "Set MANUAL_CLOSE=true"}`,
		)

	err := c.ManualClose(context.Background())

	assert.EqualError(t, err, "exception in response: Set MANUAL_CLOSE=true")
}
