package domain

type Page struct {
	Size int `json:"size" doc:"조회할 갯수" default:"20"`
	Page int `json:"page" doc:"조회할 페이지" default:"1"`
}
