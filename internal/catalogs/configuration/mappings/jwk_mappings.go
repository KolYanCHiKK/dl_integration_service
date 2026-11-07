package mappings

import (
	"encoding/json"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/configuration/features/jwkcreate/dtos"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"gorm.io/datatypes"
)

type SigningKeyParams struct {
	*ecdsa.GeneratedSigningKey
}

func MapGeneratedSigningKeyToSigningModel(signingKeyParams *SigningKeyParams) models.SigningKey {
	return models.SigningKey{
		KID:           signingKeyParams.KID,
		Algorithm:     signingKeyParams.Algorithm,
		PublicJWK:     datatypes.JSON(signingKeyParams.PublicJWK),
		PrivateKeyEnc: signingKeyParams.PrivateKeyEnc,
		Status:        oauth.SigningKeyStatusActive,
		NotBefore:     time.Now().UTC(),
	}
}

type SigningKeyDTOParams struct {
	PublicJWK json.RawMessage
	Status    string
	NotBefore time.Time
	CreatedAt time.Time
}

func MapSigningKeyDTOParamsToDTO(params *SigningKeyDTOParams) *dtos.JWKCreateResponseDTO {
	return &dtos.JWKCreateResponseDTO{
		Key:       params.PublicJWK,
		Status:    params.Status,
		NotBefore: params.NotBefore,
		CreatedAt: params.CreatedAt,
	}
}
