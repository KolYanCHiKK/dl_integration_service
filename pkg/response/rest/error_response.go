package rest

import (
	"encoding/json"
	"net/http"

	errorpkg "github.com/KolYanCHiKK/dl_integration_service/pkg/body/rest/error"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

func GenerateBadRequestResponse(w http.ResponseWriter, headerMap map[string]string, message string, errs ...string) {
	AddBaseRestErrorResponseHeaders(w, 400)

	if headerMap != nil {
		AddCustomRestResponseHeaders(w, headerMap)
	}

	body := errorpkg.BadRequestBody
	if message != "" {
		body.Errors = append(body.Errors, message)
	} else {
		body.Errors = append(body.Errors, errs...)
	}

	// Логируем ошибки
	logger2.CreateValidationErrLogs("Rest API", errs...)

	err := json.NewEncoder(w).Encode(&body)
	if err != nil {
		ToFormCreateBodyError(err, w)
		return
	}

}

func GenerateInternalServerErrorResponse(w http.ResponseWriter, headerMap map[string]string, errs ...string) {
	AddBaseRestErrorResponseHeaders(w, 500)

	if headerMap != nil {
		AddCustomRestResponseHeaders(w, headerMap)
	}

	body := errorpkg.InternalServerErrorBody
	body.Errors = append(body.Errors, errs...)

	// Логируем ошибки
	logger2.CreateValidationErrLogs("Rest API", errs...)

	err := json.NewEncoder(w).Encode(&body)
	if err != nil {
		ToFormCreateBodyError(err, w)
	}

}

func GenerateNotFoundErrorResponse(w http.ResponseWriter, headerMap map[string]string, message string, errs ...string) {
	AddBaseRestErrorResponseHeaders(w, 404)

	if headerMap != nil {
		AddCustomRestResponseHeaders(w, headerMap)
	}

	body := errorpkg.NotFoundErrorBody
	if message != "" {
		body.Errors = append(body.Errors, message)
	} else {
		body.Errors = append(body.Errors, errs...)
	}

	// Логируем ошибки
	logger2.CreateValidationErrLogs("Rest API", errs...)

	err := json.NewEncoder(w).Encode(&body)
	if err != nil {
		ToFormCreateBodyError(err, w)
	}

}

func GenerateConflictResponse(w http.ResponseWriter, headerMap map[string]string, message string, errs ...string) {
	AddBaseRestErrorResponseHeaders(w, http.StatusConflict)

	if headerMap != nil {
		AddCustomRestResponseHeaders(w, headerMap)
	}

	body := errorpkg.BadRequestBody
	if message != "" {
		body.Errors = append(body.Errors, message)
	} else {
		body.Errors = append(body.Errors, errs...)
	}

	// Логируем ошибки
	logger2.CreateValidationErrLogs("Rest API", errs...)

	err := json.NewEncoder(w).Encode(&body)
	if err != nil {
		ToFormCreateBodyError(err, w)
		return
	}

}

func ToFormCreateBodyError(err error, w http.ResponseWriter) {
	logger2.CreateErrLog(err.Error(), "Ошибка создания тела ответа JSON")
	http.Error(w, "Error body was not init", 500)
}

func AddBaseRestErrorResponseHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
