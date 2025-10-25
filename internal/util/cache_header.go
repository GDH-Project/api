package util

import (
	"fmt"
	"strings"
)

type CacheType string

var (
	CacheTypePublic  CacheType = "public"
	CacheTypePrivate CacheType = "private"
)

type CacheHeader struct {
	CacheControl string `header:"Cache-Control"`
}

// CacheHeaderBuilder
//
// "cache-control" 헤더에 삽입되는 데이터 구조체 입니다.
type CacheHeaderBuilder struct {
	CacheType CacheType // 캐시 시간
	TTL       int       // TTL
}

// String
//
// "cache-control" 헤더에 삽입되는 구조체가 조합된 문자열 입니다.
func (c *CacheHeaderBuilder) String() string {
	maxAgeStr := fmt.Sprintf("max-age=%d", c.TTL)

	if c.CacheType == "" {
		return maxAgeStr
	}

	var b strings.Builder
	b.WriteString(string(c.CacheType))
	b.WriteString(", ")
	b.WriteString(maxAgeStr)

	return b.String()
}
