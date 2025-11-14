package services

import (
	"database/sql"
	"errors"

	"github.com/rizkyhaksono/sarana-ai-take-home-test/constants"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/models"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/utils"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

type LogsResponse struct {
	Logs                     []models.Log `json:"logs"`
	utils.PaginationResponse `json:",inline"`
}

func (s *LogService) GetLogsWithParams(params utils.PaginationParams) (*LogsResponse, error) {
	// Validate pagination parameters
	validSortFields := map[string]bool{
		"datetime":    true,
		"created_at":  true,
		"method":      true,
		"endpoint":    true,
		"status_code": true,
	}
	utils.ValidatePaginationParams(&params, validSortFields, "datetime")

	// Build paginated query
	baseQuery := "SELECT id, datetime, method, endpoint, headers, request_body, response_body, status_code, created_at FROM logs"
	countQuery := "SELECT COUNT(*) FROM logs"
	whereCondition := ""
	searchFields := []string{"method", "endpoint", "request_body", "response_body"}
	baseArgs := []interface{}{}

	query, countQueryFinal, args, err := utils.BuildPaginatedQuery(
		baseQuery,
		countQuery,
		whereCondition,
		searchFields,
		params,
		baseArgs,
	)
	if err != nil {
		return nil, errors.New(constants.ErrFetchingNotes)
	}

	// Get total count (remove LIMIT and OFFSET from args)
	countArgs := args[:len(args)-2]
	total, err := utils.GetTotalCount(database.DB, countQueryFinal, countArgs)
	if err != nil {
		return nil, errors.New("failed to fetch logs count")
	}

	// Execute query
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, errors.New("failed to fetch logs")
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		var headers, requestBody, responseBody sql.NullString
		if err := rows.Scan(&log.ID, &log.Datetime, &log.Method, &log.Endpoint, &headers, &requestBody, &responseBody, &log.StatusCode, &log.CreatedAt); err != nil {
			return nil, errors.New("failed to fetch logs")
		}
		if headers.Valid {
			log.Headers = headers.String
		}
		if requestBody.Valid {
			log.RequestBody = requestBody.String
		}
		if responseBody.Valid {
			log.ResponseBody = responseBody.String
		}
		logs = append(logs, log)
	}

	// Calculate pagination metadata
	paginationMeta := utils.CalculatePaginationMetadata(total, params.Page, params.Limit)

	return &LogsResponse{
		Logs:               logs,
		PaginationResponse: paginationMeta,
	}, nil
}

func (s *LogService) GetLogByID(logID int) (*models.Log, error) {
	var log models.Log
	var headers, requestBody, responseBody sql.NullString
	err := database.DB.QueryRow(
		"SELECT id, datetime, method, endpoint, headers, request_body, response_body, status_code, created_at FROM logs WHERE id = $1",
		logID,
	).Scan(&log.ID, &log.Datetime, &log.Method, &log.Endpoint, &headers, &requestBody, &responseBody, &log.StatusCode, &log.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("log not found")
		}
		return nil, errors.New("failed to fetch log")
	}

	if headers.Valid {
		log.Headers = headers.String
	}
	if requestBody.Valid {
		log.RequestBody = requestBody.String
	}
	if responseBody.Valid {
		log.ResponseBody = responseBody.String
	}

	return &log, nil
}
