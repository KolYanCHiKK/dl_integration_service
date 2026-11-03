package request

import (
	"encoding/json"
	"net/http"
)

func DecodeRestRequest[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var payload T
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}
