package stacklog

import (
	"fmt"
	"strings"
)

// CheckType builds a short detail string for variadic values.
func CheckType(datas ...any) string {
	if len(datas) == 0 {
		return ""
	}
	var result []string
	for _, data := range datas {
		switch v := data.(type) {
		case string:
			result = append(result, fmt.Sprintf("[%s]", v))
		case int, int64, bool, float64, float32, uint, uint64:
			result = append(result, fmt.Sprintf("[%v]", v))
		default:
			if v == nil {
				result = append(result, "[nil]")
			} else {
				result = append(result, fmt.Sprintf("[%+v]", v))
			}
		}
	}
	return " | Details: " + strings.Join(result, " , ")
}
