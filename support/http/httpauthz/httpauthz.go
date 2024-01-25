// Package httpauthz contains helper functions for
// parsing the 'Authorization' header in HTTP requests.
package httpauthz

import "strings"

// ParseBearerToken parses a bearer token's value from a HTTP Authorization
// header. If the prefix of the authorization header value is 'Bearer ' (case
// ignored) the rest of the value that follows the prefix is returned,
// otherwise an empty string is returned.
func ParseBearerToken(authorizationHeaderValue string) string {
	const prefix = "Bearer "
	if hasPrefixFold(authorizationHeaderValue, prefix) {
		return authorizationHeaderValue[len(prefix):]
	}
	return ""
}

func hasPrefixFold(s string, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return strings.EqualFold(s[0:len(prefix)], prefix)
}
