package cmp

import (
	"fmt"
	"net/url"
)

type Path struct {
	Path   string
	Level  int
	Line   int
	Stream bool
}

// ID returns a path identifier. It is saved to not repeat the same requests
// over again.
func (p Path) ID() string {
	path := removeRandomC(p.Path)
	return fmt.Sprintf("%t%s", p.Stream, path)
}

// removeRandomC removes random `c` param that is part of many requests
// and originates in js-diamcircle-sdk
func removeRandomC(path string) string {
	urlObj, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	q := urlObj.Query()
	q.Del("c")

	urlObj.RawQuery = q.Encode()
	return urlObj.String()
}
