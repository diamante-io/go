package auroraclient

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/diamcircle/go/support/errors"
)

// decodeResponse decodes the response from a request to a aurora server
func decodeResponse(resp *http.Response, object interface{}, hc *Client) (err error) {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	u, err := url.Parse(hc.AuroraURL)
	if err != nil {
		return errors.Errorf("unable to parse the provided aurora url: %s", hc.AuroraURL)
	}
	setCurrentServerTime(u.Hostname(), resp.Header["Date"], hc)

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		auroraError := &Error{
			Response: resp,
		}
		decodeError := decoder.Decode(&auroraError.Problem)
		if decodeError != nil {
			return errors.Wrap(decodeError, "error decoding aurora.Problem")
		}
		return auroraError
	}

	err = decoder.Decode(&object)
	if err != nil {
		return errors.Wrap(err, "error decoding response")
	}
	return
}

// countParams counts the number of parameters provided
func countParams(params ...interface{}) int {
	counter := 0
	for _, param := range params {
		switch param := param.(type) {
		case string:
			if param != "" {
				counter++
			}
		case int:
			if param > 0 {
				counter++
			}
		case uint:
			if param > 0 {
				counter++
			}
		case bool:
			counter++
		default:
			panic("Unknown parameter type")
		}

	}
	return counter
}

// addQueryParams sets query parameters for a url
func addQueryParams(params ...interface{}) string {
	query := url.Values{}

	for _, param := range params {
		switch param := param.(type) {
		case cursor:
			if param != "" {
				query.Add("cursor", string(param))
			}
		case Order:
			if param != "" {
				query.Add("order", string(param))
			}
		case limit:
			if param != 0 {
				query.Add("limit", strconv.Itoa(int(param)))
			}
		case assetCode:
			if param != "" {
				query.Add("asset_code", string(param))
			}
		case assetIssuer:
			if param != "" {
				query.Add("asset_issuer", string(param))
			}
		case includeFailed:
			if param {
				query.Add("include_failed", "true")
			}
		case join:
			if param != "" {
				query.Add("join", string(param))
			}
		case reserves:
			if len(param) > 0 {
				query.Add("reserves", strings.Join(param, ","))
			}
		case map[string]string:
			for key, value := range param {
				if value != "" {
					query.Add(key, value)
				}
			}
		default:
			panic("Unknown parameter type")
		}
	}

	return query.Encode()
}

// setCurrentServerTime saves the current time returned by a aurora server
func setCurrentServerTime(host string, serverDate []string, hc *Client) {
	if len(serverDate) == 0 {
		return
	}
	st, err := time.Parse(time.RFC1123, serverDate[0])
	if err != nil {
		return
	}
	serverTimeMapMutex.Lock()
	ServerTimeMap[host] = ServerTimeRecord{ServerTime: st.UTC().Unix(), LocalTimeRecorded: hc.clock.Now().UTC().Unix()}
	serverTimeMapMutex.Unlock()
}

// currentServerTime returns the current server time for a given aurora server
func currentServerTime(host string, currentTimeUTC int64) int64 {
	serverTimeMapMutex.Lock()
	st := ServerTimeMap[host]
	serverTimeMapMutex.Unlock()
	if &st == nil {
		return 0
	}

	// if it has been more than 5 minutes from the last time, then return 0 because the saved
	// server time is behind.
	if currentTimeUTC-st.LocalTimeRecorded > 60*5 {
		return 0
	}
	return currentTimeUTC - st.LocalTimeRecorded + st.ServerTime
}
