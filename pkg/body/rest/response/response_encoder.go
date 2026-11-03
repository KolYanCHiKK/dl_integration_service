package response

import (
	"encoding/json"
	"net/http"
)

func EncodeRestRequest[T any](w http.ResponseWriter, body T) error {
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		return err
	}

	return nil
}
