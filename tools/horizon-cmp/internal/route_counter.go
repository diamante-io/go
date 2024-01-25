package cmp

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// MakeRoute translates route with * wildcards into a regexp. It adds start/end
// of line asserts, changes * into any characters except when in used for lists
// (like /accounts/*) - in such case it ensures there is no more `/` characters.
// This is not ideal and requires checking routes in correct order but is enough
// for aurora-cmp.
// More here: https://regex101.com/r/EhBRtS/1
func MakeRoute(route string) *Route {
	name := route
	route = "^" + route
	route = strings.ReplaceAll(route, "*", "[^?/]+") // everything except `/` and `?`
	route = route + "[?/]?[^/]*$"                    // ? or / or nothing and then everything except `/``
	return &Route{
		name:   name,
		regexp: regexp.MustCompile(route),
	}
}

type Route struct {
	name    string
	regexp  *regexp.Regexp
	counter int
}

type Routes struct {
	List      []*Route
	unmatched []string
	mutex     sync.Mutex
}

func (r *Routes) Count(path string) {
	for _, route := range r.List {
		if route.regexp.Match([]byte(path)) {
			r.mutex.Lock()
			route.counter++
			r.mutex.Unlock()
			return
		}
	}

	r.mutex.Lock()
	r.unmatched = append(r.unmatched, path)
	r.mutex.Unlock()
}

func (r *Routes) Print() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	fmt.Println("Routes coverage:")
	for _, route := range r.List {
		fmt.Println(route.counter, route.name)
	}

	fmt.Println("Unmatched paths:")
	for _, path := range r.unmatched {
		fmt.Println(path)
	}
}
