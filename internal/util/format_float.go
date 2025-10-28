package util

import (
	"strconv"
	"strings"
)

// FormatFloat 소수점 정밀도 제거 유틸
//
// prec를 통해 소수점 자릿수를 지정한다.
func FormatFloat(f float64, prec int) string {
	s := strconv.FormatFloat(f, 'f', prec, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")

	return s
}
