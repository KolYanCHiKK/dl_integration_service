package dtos

import (
	json2 "encoding/json"
	"time"
)

type JWKCreateResponseDTO struct {
	Key       json2.RawMessage `json:"key"`
	Status    string           `json:"status"`
	NotBefore time.Time        `json:"notBefore"`
	CreatedAt time.Time        `json:"createdAt"`
}
