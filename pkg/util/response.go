package util

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (r *Response) JSON(w http.ResponseWriter, success bool, statusCode int, message string, data interface{}) error {
	res := &Response{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(res)
	if err != nil {
		return err
	}

	return nil
}
