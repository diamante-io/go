package scraper

import "net/url"

// nextCursor finds the cursor parameter on the "next" link of a hProtocol page.
func nextCursor(nextPageURL string) (cursor string, err error) {
	u, err := url.Parse(nextPageURL)
	if err != nil {
		return
	}

	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	cursor = m["cursor"][0]

	return
}
