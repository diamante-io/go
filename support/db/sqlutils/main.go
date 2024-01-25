// Package sqlutils contains utility functions for manipulating strings of SQL
package sqlutils

import (
	"regexp"
	"strings"
)

// AllStatements takes a sql script, possibly containing comments and multiple
// statements, and returns a slice of strings that correspond to each individual
// SQL statement within the script
func AllStatements(script string) (ret []string) {
	for _, s := range strings.Split(RemoveComments(script), ";") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		ret = append(ret, s)
	}
	return
}

// RemoveComments trims out any comment blocks or lines from the provided SQL
// script
func RemoveComments(script string) string {
	withoutBlocks := sqlBlockComments.ReplaceAllString(script, "")
	return sqlLineComments.ReplaceAllString(withoutBlocks, "")
}

// SQLBlockComments is a regex that matches against SQL block comments
var sqlBlockComments = regexp.MustCompile(`/\*.*?\*/`)

// SQLLineComments is a regex that matches against SQL line comments
var sqlLineComments = regexp.MustCompile("--.*?\n")
