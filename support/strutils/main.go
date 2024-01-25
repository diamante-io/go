package strutils

import "strings"

// KebabToConstantCase converts a string from lower-case-dashes to UPPER_CASE_UNDERSCORES.
func KebabToConstantCase(kebab string) (constantCase string) {
	constantCase = strings.ToUpper(kebab)
	constantCase = strings.Replace(constantCase, "-", "_", -1)

	return constantCase
}
