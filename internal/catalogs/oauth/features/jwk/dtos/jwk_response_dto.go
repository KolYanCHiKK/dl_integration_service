package dtos

import "encoding/json"

type GetJWKSigningKeyResponseDTO struct {
	Keys []json.RawMessage `json:"keys"`
}
