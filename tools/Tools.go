package tools

import (
	"fmt"
	"strings"
)

func ArrayToString(a [3]int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", "", -1), "[]")
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func DivCheck(x, y int) (q int, ok bool) {
	defer func() {
		recover()
	}()
	q = x / y
	return q, true
}
