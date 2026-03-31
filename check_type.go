package logging

import (
	"fmt"
	"strings"
)

// CheckType renders variadic values as short debug details in logs.
func CheckType(datas ...any) string {
	if len(datas) == 0 {
		return ""
	}

	var result []string
	for _, data := range datas {
		switch v := data.(type) {
		case string:
			result = append(result, fmt.Sprintf("[%s]", v))
		case int, int64, bool:
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
