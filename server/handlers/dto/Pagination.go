package dto

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// `page` is 1-based
type PaginationQuery struct {
	PageIndex    int               `json:"pageIndex"`
	ItemsPerPage int               `json:"itemsPerPage"`
	Filters      map[string]string `json:"filters"`
}

type PaginationResponse[T interface{}] struct {
	Items        []T `json:"items"`
	Total        int `json:"total"`
	PageIndex    int `json:"pageIndex"`
	ItemsPerPage int `json:"itemsPerPage"`
}

const (
	defaultItemsPerPage = 10
	defaultPageIndex    = 0
)

func ParsePaginationRequest(c echo.Context) *PaginationQuery {
	p := new(PaginationQuery)
	err := echo.QueryParamsBinder(c).
		Int("itemsPerPage", &p.ItemsPerPage).
		Int("pageIndex", &p.PageIndex).
		CustomFunc("filters", func(_ []string) []error {
			err := json.Unmarshal([]byte(c.QueryParam("filters")), &p.Filters)
			return []error{err}
		}).
		BindError()

	if err != nil {
		p.ItemsPerPage = defaultItemsPerPage
		p.PageIndex = defaultPageIndex
	}

	if p.ItemsPerPage <= 0 {
		p.ItemsPerPage = defaultItemsPerPage
	}

	if p.PageIndex < 0 {
		p.PageIndex = defaultPageIndex
	}

	return p
}
