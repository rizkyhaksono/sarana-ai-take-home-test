package middleware

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rizkyhaksono/sarana-ai-take-home-test/database"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		method := c.Method()
		endpoint := c.Path()

		var requestBody string
		if c.Body() != nil {
			requestBody = string(c.Body())
		}

		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headerKey := string(key)
			headerValue := string(value)

			if strings.ToLower(headerKey) == "authorization" {
				headerValue = "***MASKED***"
			}
			headers[headerKey] = headerValue
		})
		headersJSON, _ := json.Marshal(headers)

		responseBody := &bytes.Buffer{}

		err := c.Next()

		responseBodyStr := string(c.Response().Body())
		statusCode := c.Response().StatusCode()
		// duration := time.Since(startTime)

		go func() {
			_, dbErr := database.DB.Exec(`
				INSERT INTO logs (datetime, method, endpoint, headers, request_body, response_body, status_code)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
			`, startTime, method, endpoint, string(headersJSON), requestBody, responseBodyStr, statusCode)

			if dbErr != nil {
				println("Failed to save log to database:", dbErr.Error())
			}

			// Send to Loki
			// lokiMessage := fmt.Sprintf("duration=%dms request_body=%s response_body=%s",
			// 	duration.Milliseconds(), requestBody, responseBodyStr)

			// lokiErr := utils.SendLogToLoki(c.Context(), getLogLevel(statusCode), method, endpoint, lokiMessage, statusCode)
			// if lokiErr != nil {
			// 	println("Failed to send log to Loki:", lokiErr.Error())
			// }
		}()

		_ = responseBody

		return err
	}
}

func getLogLevel(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	default:
		return "info"
	}
}
