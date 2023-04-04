package file

import (
	"strings"

	"mango-admin/pkg/config/encoder"
)

func format(p string, e encoder.Encoder) string {
	parts := strings.Split(p, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return e.String()
}
