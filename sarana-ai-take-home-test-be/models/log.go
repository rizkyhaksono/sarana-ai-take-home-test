package models

import "time"

type Log struct {
	ID           int       `json:"id"`
	Datetime     time.Time `json:"datetime"`
	Method       string    `json:"method"`
	Endpoint     string    `json:"endpoint"`
	Headers      string    `json:"headers"`
	RequestBody  string    `json:"request_body"`
	ResponseBody string    `json:"response_body"`
	StatusCode   int       `json:"status_code"`
	CreatedAt    time.Time `json:"created_at"`
}
