package auroraclient

import (
	"net/http"

	"go/support/errors"
)

// BuildURL returns the url for getting fee stats about a running aurora instance
func (fr feeStatsRequest) BuildURL() (endpoint string, err error) {
	endpoint = fr.endpoint
	if endpoint == "" {
		err = errors.New("invalid request: too few parameters")
	}

	return
}

// HTTPRequest returns the http request for the fee stats endpoint
func (fr feeStatsRequest) HTTPRequest(auroraURL string) (*http.Request, error) {
	endpoint, err := fr.BuildURL()
	if err != nil {
		return nil, err
	}

	return http.NewRequest("GET", auroraURL+endpoint, nil)
}
