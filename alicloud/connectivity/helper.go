package connectivity

import "strings"

func ConvertKebabToSnake(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
