package sorting_domain

import "strings"

type Sorting struct {
	Attribute string
	Direction SortingDirection
}

const (
	SortingDirectionAsc  SortingDirection = "ASC"
	SortingDirectionDesc SortingDirection = "DESC"
)

type SortingDirection string

func (d SortingDirection) Valid() bool {
	direction := SortingDirection(strings.ToUpper(string(d)))
	return direction == SortingDirectionAsc || direction == SortingDirectionDesc
}
