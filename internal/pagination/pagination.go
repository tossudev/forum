package pagination

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Pagination struct {
	Page     int
	PageSize int
}

// PaginationNav contains page data to be used in pagination navigation with go templates
type PaginationNav struct {
	Path       string // The URL path for navigation
	Page       int    // Current (wanted) page
	TotalPages int    // Last possible page
	PageList   []int  // The visible pages shown in the navigation. Get with GetDisplayPageNums.
	PrevPage   int    // The page before the current page
	NextPage   int    // The page after the current page
}

const (
	minPage     = 1
	maxPage     = 100000000 // capped to 100 million prevent int64 overflow on Offset calculation
	minPageSize = 1
	maxPageSize = 100
)

var ErrPaginationParams = errors.New("invalid pagination parameter")

// Limit returns the number of results per page (i.e. page size).
func (p *Pagination) Limit() int {
	return p.PageSize
}

// Offset calculates how many entries to skip to reach the specified page.
// The page numbering is 1-based.
func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Validate enforces pagination parameter limits.
func (p *Pagination) Validate() error {
	if p.Page < minPage {
		return fmt.Errorf("%w: page must not be negative", ErrPaginationParams)
	}

	if p.Page > maxPage {
		return fmt.Errorf("%w: page too large", ErrPaginationParams)
	}

	if p.PageSize < minPageSize || p.PageSize > maxPageSize {
		return fmt.Errorf("%w: size must be in range %d - %d", ErrPaginationParams, minPageSize, maxPageSize)
	}

	return nil
}

// Parse parses the pagination parameters from URL query values into a Pagination struct.
func Parse(query url.Values) (Pagination, error) {
	// Set default vaules
	pagination := Pagination{
		Page:     minPage,
		PageSize: 5,
	}

	// Parse "page" from query
	if queryPage := query.Get("page"); queryPage != "" {
		page, err := strconv.Atoi(queryPage)
		if err != nil {
			var numErr *strconv.NumError
			if errors.As(err, &numErr) && errors.Is(numErr.Err, strconv.ErrRange) {
				return Pagination{}, fmt.Errorf("%w: page too large", ErrPaginationParams)
			}
			return Pagination{}, fmt.Errorf("%w: page must be an integer", ErrPaginationParams)
		}
		pagination.Page = page
	}

	// Parse "size" from query
	if querySize := query.Get("size"); querySize != "" {
		size, err := strconv.Atoi(querySize)
		if err != nil {
			var numErr *strconv.NumError
			if errors.As(err, &numErr) && errors.Is(numErr.Err, strconv.ErrRange) {
				return Pagination{}, fmt.Errorf("%w: size too large", ErrPaginationParams)
			}
			return Pagination{}, fmt.Errorf("%w: size must be an integer", ErrPaginationParams)
		}
		pagination.PageSize = size
	}

	return pagination, nil
}

// CalcTotalPages calculates the total pages of paginated data available (equivalent to the last page number)
func CalcTotalPages(totalCount, pageSize int) int {
	if totalCount == 0 {
		return minPage
	}
	return (totalCount + pageSize - 1) / pageSize // round up trick for totalCount/pageSize
}

// GetDisplayPageNums takes the wanted page, total number of pages, and the maximum number of pages
// wanted to display in navigation and returns a slice of page numbers (up to displayLimit long).
// If the wanted page is greater than the display limit and there are enough pages below and above it, the wanted page will be in the center.
// e.g. want=1, total=10, displayLimit=5 -> []int{1,2,3,4,5}
func GetDisplayPageNums(want, total, displayLimit int) []int {
	var pages []int

	// e.g. want=4, total=10, displayLimit=5 -> []int{1,2,3,4,5}
	if total <= displayLimit || want <= displayLimit {
		max := min(total, displayLimit)
		for i := 1; i <= max; i++ {
			pages = append(pages, i)
		}
		return pages
	}

	// For centering: if there can be enough pages below and above the wanted page, the wanted page will be in the middle.
	// e.g. want=7, total=10, displayLimit=5 -> []int{5,6,7,8,9}
	leftPadding := (displayLimit - 1) / 2 // -1 for center number
	rightPadding := total - leftPadding - 1
	if want+rightPadding <= total {
		for i := want - leftPadding; i <= want+rightPadding; i++ {
			pages = append(pages, i)
		}
		return pages
	}

	// Not enough pages above the wanted page for centering, so put the last page available at the end.
	// e.g. want=9, total=10, displayLimit=5 -> []int{6,7,8,9,10}
	if want+rightPadding > total {
		for i := total - displayLimit + 1; i <= total; i++ {
			pages = append(pages, i)
		}
	}
	return pages
}

func CheckPageTooLarge(paginationData Pagination, count int) error {
	if paginationData.Page > CalcTotalPages(count, paginationData.PageSize) {
		return fmt.Errorf("%w: page number too large", ErrPaginationParams)
	}
	return nil
}
