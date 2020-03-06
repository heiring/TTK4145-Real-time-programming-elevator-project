package tools

import (
	"fmt"
	"strings"
)

func ArrayToString(a [3]int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", "", -1), "[]")
}
