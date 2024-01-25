package config

import (
	"fmt"
	"strings"
)

func (err *InvalidConfigError) Error() string {
	fields := make([]string, 0, len(err.InvalidFields))
	for key := range err.InvalidFields {
		fields = append(fields, key)
	}

	return fmt.Sprintf(`invalid fields: %s`, strings.Join(fields, ","))
}
