package util

import (
	"net/http"
)

const (
	SUCCESS          = "SUCCESS"
	DUPLICATE        = "DUPLICATE"
	NOT_FOUND        = "NOT_FOUND"
	SYSTEM_ERROR     = "SYSTEM_ERROR"
	VALIDATION_ERROR = "VALIDATION_ERROR"
)

func GetResCode(state string) int {
	var code int
	switch state {
	case SUCCESS:
		code = http.StatusOK
	case DUPLICATE:
		code = http.StatusBadRequest
	case NOT_FOUND:
		code = http.StatusNotFound
	case VALIDATION_ERROR:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	return code
}
