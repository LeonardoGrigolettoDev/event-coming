package repositories

import (
	"fmt"
	"strings"
)

func buildWhere(filters map[string]interface{}, startIdx int) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(filters))
	args := make([]interface{}, 0, len(filters))
	i := startIdx
	for k, v := range filters {
		parts = append(parts, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	return "WHERE " + strings.Join(parts, " AND "), args
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }
