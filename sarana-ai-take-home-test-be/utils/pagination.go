package utils

import (
	"database/sql"
	"fmt"
)

type PaginationParams struct {
	Search string
	SortBy string
	Order  string
	Page   int
	Limit  int
}

type PaginationResponse struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
}

func ValidatePaginationParams(params *PaginationParams, validSortFields map[string]bool, defaultSortBy string) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.SortBy == "" {
		params.SortBy = defaultSortBy
	}
	if params.Order == "" {
		params.Order = "DESC"
	}

	// Validate sortBy field
	if !validSortFields[params.SortBy] {
		params.SortBy = defaultSortBy
	}

	// Validate order
	if params.Order != "ASC" && params.Order != "DESC" {
		params.Order = "DESC"
	}
}

// BuildPaginatedQuery builds a paginated SQL query with optional search conditions
// baseQuery: the base SELECT query without WHERE clause
// countQuery: the base COUNT query without WHERE clause
// whereCondition: the WHERE clause (e.g., "user_id = $1")
// searchFields: fields to search in (e.g., []string{"title", "content"})
// params: pagination parameters
// baseArgs: arguments for the base WHERE condition
// Returns: final query, count query, arguments array, and argument index
func BuildPaginatedQuery(
	baseQuery string,
	countQuery string,
	whereCondition string,
	searchFields []string,
	params PaginationParams,
	baseArgs []interface{},
) (string, string, []interface{}, error) {
	query := baseQuery
	count := countQuery
	args := make([]interface{}, len(baseArgs))
	copy(args, baseArgs)
	argIndex := len(baseArgs) + 1

	// Add WHERE condition
	if whereCondition != "" {
		query += " WHERE " + whereCondition
		count += " WHERE " + whereCondition
	}

	// Add search condition if search parameter is provided
	if params.Search != "" && len(searchFields) > 0 {
		searchCondition := " AND ("
		for i, field := range searchFields {
			if i > 0 {
				searchCondition += " OR "
			}
			searchCondition += fmt.Sprintf("%s ILIKE $%d", field, argIndex)
		}
		searchCondition += ")"
		query += searchCondition
		count += searchCondition
		args = append(args, "%"+params.Search+"%")
		argIndex++
	}

	// Add sorting
	query += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, params.Order)

	// Add pagination
	offset := (params.Page - 1) * params.Limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	return query, count, args, nil
}

func CalculatePaginationMetadata(total, page, limit int) PaginationResponse {
	totalPages := (total + limit - 1) / limit
	return PaginationResponse{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}
}

func GetTotalCount(db *sql.DB, countQuery string, args []interface{}) (int, error) {
	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	return total, err
}
