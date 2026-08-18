package shared

import "strconv"

// PaginationParams captures common list-query inputs. Keep it intentionally
// small; modules can embed it if needed.
type PaginationParams struct {
	Page  int
	Limit int
}

// DefaultPagination returns a PaginationParams with the supplied defaults.
func DefaultPagination(page, limit int) PaginationParams {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return PaginationParams{Page: page, Limit: limit}
}

// Offset returns the SQL OFFSET for the current page.
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Page returns the current page, clamped to at least 1.
func (p PaginationParams) PageClamped() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// Limit returns the current limit, clamped to a sane range.
func (p PaginationParams) LimitClamped() int {
	if p.Limit <= 0 {
		return 20
	}
	if p.Limit > 100 {
		return 100
	}
	return p.Limit
}

// ParsePagination parses raw string values into a PaginationParams.
func ParsePagination(pageStr, limitStr string) PaginationParams {
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	return DefaultPagination(page, limit)
}

// PageMeta is the pagination metadata attached to list responses.
type PageMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPageMeta computes TotalPages from a total count.
func NewPageMeta(params PaginationParams, total int) PageMeta {
	totalPages := 0
	if total > 0 && params.Limit > 0 {
		totalPages = (total + params.Limit - 1) / params.Limit
	}
	return PageMeta{
		Page:       params.PageClamped(),
		Limit:      params.LimitClamped(),
		Total:      total,
		TotalPages: totalPages,
	}
}
