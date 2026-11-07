package mappings

import (
	"encoding/json"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/jwk/dtos"
	"gorm.io/datatypes"
)

type GetJWKParameters struct {
	JsonJWK []json.RawMessage
}

func MapJSONToSliceRawMessages(datatypesJWK []datatypes.JSON) []json.RawMessage {
	var raw []json.RawMessage
	for _, value := range datatypesJWK {
		raw = append(raw, json.RawMessage(value))
	}

	return raw
}

func MapGetJWKParametersToDTO(params *GetJWKParameters) *dtos.GetJWKSigningKeyResponseDTO {
	return &dtos.GetJWKSigningKeyResponseDTO{Keys: params.JsonJWK}
}
