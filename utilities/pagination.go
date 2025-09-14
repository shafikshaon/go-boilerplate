package utilities

import "math"

// CalculateOffset returns a safe page and perPage along with the SQL offset.
func CalculateOffset(page, perPage int) (safePage, safePerPage, offset int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	return page, perPage, (page - 1) * perPage
}

// TotalPages calculates total pages given total items and page size.
func TotalPages(total, perPage int) int {
	if perPage <= 0 {
		perPage = 10
	}
	return int(math.Ceil(float64(total) / float64(perPage)))
}
