package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func (r *Response) JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	res, err := json.Marshal(data)

	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
}

func (r *Response) ReturnedData(success bool, statusCode int, message string, data interface{}) *Response {
	return &Response{
		Success:    success,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}
