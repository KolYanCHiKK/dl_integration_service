package rest

import (
	"encoding/json"
	"net/http"

	errorpkg "github.com/KolYanCHiKK/dl_integration_service/pkg/body/rest/error"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

func BadRequestResponse(w http.ResponseWriter, errs ...string) {
	AddBaseRestErrorResponseHeaders(w, 400)

	body := errorpkg.BadRequestBody
	body.Errors = append(body.Errors, errs...)

	err := json.NewEncoder(w).Encode(&body)
	if err != nil {
		logger2.CreateErrLog(err.Error(), "Ошибка создания тела ответа JSON")
		http.Error(w, "Error body was not init", 500)
	}

}

func AddBaseRestErrorResponseHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
