package querying

import (
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/queries/qm"
)

const (
	//SortOrderDesc oder field descending
	SortOrderDesc = "desc"

	//SortOrderASC oder dield ascending
	SortOrderASC = "asc"
)

//IsSort checks if sortorder specified and returns the sql-ordering
func IsSort(order string, field string) (bool, string) {
	if strings.EqualFold(SortOrderDesc, order) {
		return true, fmt.Sprintf("%s %s", field, SortOrderDesc)
	}
	if strings.EqualFold(SortOrderASC, order) {
		return true, fmt.Sprintf("%s %s", field, SortOrderASC)
	}
	return false, ""
}

//AddSorting adds a sorting clause to the querymod, if specified
func AddSorting(order string, field string, q []qm.QueryMod) []qm.QueryMod {
	if isSort, order := IsSort(order, field); isSort {
		q = append(q, qm.OrderBy(order))
	}

	return q
}

//Like returns the filterCriteria prefixed and suffixed with the like operator
func Like(s string) string {
	return fmt.Sprintf("%%%s%%", s)
}

//GetSQLKeyString will replace the given keys in the sql string
func GetSQLKeyString(sql string, keys map[string]string) string {
	s := sql
	for key, value := range keys {
		s = strings.Replace(s, key, value, -1)
	}
	return s
}
