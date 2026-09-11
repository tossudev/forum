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

const (
	minPage     = 0
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
// The page numbering is 0-based.
func (p *Pagination) Offset() int {
	return p.Page * p.PageSize
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
		PageSize: 10,
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
