package domain

import "math"

type Page struct {
	Size        int `json:"size" doc:"조회할 갯수" minimum:"5" default:"20"`
	Page        int `json:"page" doc:"조회할 페이지" minimum:"1" default:"1"`
	RecordCount int `json:"-"` // 전체 레코드 수
}

// Offset 현재 오프셋
func (p *Page) Offset() int {
	return (p.Page - 1) * p.Size
}

// HasNext 다음 페이지 존재 여부
func (p *Page) HasNext() bool {
	return p.Page < p.TotalPageCount()
}

// TotalPageCount 총 페이지 수
func (p *Page) TotalPageCount() int {
	return int(math.Ceil(float64(p.RecordCount) / float64(p.Size)))
}
