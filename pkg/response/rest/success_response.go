package rest

import (
	"encoding/json"
	"net/http"
)

func BuildRestResponse(w http.ResponseWriter, body any, statusCode int, headerMap map[string]string) {
	AddBaseRestSuccessResponseHeaders(w, statusCode)

	if headerMap != nil {
		AddCustomRestResponseHeaders(w, headerMap)
	}

	if body == nil && statusCode == http.StatusNoContent {
		return
	}

	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		GenerateInternalServerErrorResponse(w, nil, err.Error())
	}
}

func AddBaseRestSuccessResponseHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

func AddCustomRestResponseHeaders(w http.ResponseWriter, headerMap map[string]string) {
	for header, value := range headerMap {
		w.Header().Set(header, value)
	}
}
