package rest

import (
	"encoding/json"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
)

func DecodeRestRequest(w http.ResponseWriter, r *http.Request, variable any) {
	err := json.NewDecoder(r.Body).Decode(&variable)
	if err != nil {
		rest.GenerateBadRequestResponse(w, nil, err.Error())
	}

}
