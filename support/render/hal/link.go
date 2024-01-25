package hal

import (
	"strings"
)

type Link struct {
	Href      string `json:"href"`
	Templated bool   `json:"templated,omitempty"`
}

func (l *Link) PopulateTemplated() {
	l.Templated = strings.Contains(l.Href, "{")
}

func NewLink(href string) Link {
	l := Link{Href: href}
	l.PopulateTemplated()
	return l
}
