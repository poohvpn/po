package pooh

import (
	"fmt"
)

func Strings(values ...interface{}) []string {
	ss := make([]string, len(values))
	for i, v := range values {
		ss[i] = fmt.Sprint(v)
	}
	return ss
}
