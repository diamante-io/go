package httptest

// On is the entrypoint method into this packages "client" mocking system.
func (c *Client) On(method string, url string) *ClientExpectation {
	return &ClientExpectation{
		Method: method,
		URL:    url,
		Client: c,
	}
}
