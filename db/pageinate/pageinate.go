package pageinate

import (
	"github.com/volatiletech/sqlboiler/queries/qm"
)

//PaginationRequestStruct represents a pagination request
type PaginationRequestStruct struct {
	PageNumber   int `form:"page_number"`
	PerPageCount int `form:"per_page_count"`
}

//PaginationResponseStruct represents the pagination results
type PaginationResponseStruct struct {
	TotalCount int64 `json:"total_count"`
	PageNumber int   `json:"page_number"`
}

//CountStruct represents the count returned by the query
type CountStruct struct {
	TotalCount int64 `boil:"total_count"`
}

//Paginate returns the queryMods for the paginated data
func Paginate(q []qm.QueryMod, r *PaginationRequestStruct, resp *PaginationResponseStruct) []qm.QueryMod {
	pageNumer := r.PageNumber
	if pageNumer <= 0 {
		pageNumer = 1
	}

	pageSize := r.PerPageCount
	if pageSize <= 0 {
		pageSize = 20 //default to 20 items per page
	}

	q = append(q, qm.Offset(pageSize*(pageNumer-1)))
	q = append(q, qm.Limit(pageSize))
	resp.PageNumber = pageNumer
	return q
}
