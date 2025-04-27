package models

// SortParams handle sorting.
type SortParams struct {
	Field string
	ASC   bool
}

// Filter defines filter structure. Operators are not yet supported.
type Filter struct {
	Field string
	Value string
}

// FilterParams contains all filters.
type FilterParams struct {
	Filters []Filter
}

// ParamsBag is used for query parameters handling.
type ParamsBag struct {
	Sort   SortParams
	Filter FilterParams
}
