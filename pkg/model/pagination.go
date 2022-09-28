package model

// Pagination ...
type Pagination struct {
	CurrentPage         int
	PaginationStartPage int
	PaginationEndPage   int
	ItemsPerPage        int
	TotalItems          int
	TotalPages          int
}
