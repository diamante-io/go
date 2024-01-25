package auroraclient

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/diamcircle/go/support/errors"
)

// BuildURL returns the url for submitting transactions to a running aurora instance
func (sr submitRequest) BuildURL() (endpoint string, err error) {
	if sr.endpoint == "" || sr.transactionXdr == "" {
		return endpoint, errors.New("invalid request: too few parameters")
	}

	query := url.Values{}
	query.Set("tx", sr.transactionXdr)

	endpoint = fmt.Sprintf("%s?%s", sr.endpoint, query.Encode())
	return endpoint, err
}

// HTTPRequest returns the http request for submitting transactions to a running aurora instance
func (sr submitRequest) HTTPRequest(auroraURL string) (*http.Request, error) {
	form := url.Values{}
	form.Set("tx", sr.transactionXdr)
	request, err := http.NewRequest("POST", auroraURL+sr.endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}
