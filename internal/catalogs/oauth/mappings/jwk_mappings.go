package mappings

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/ecdsa"
	"gorm.io/datatypes"
)

func MapGeneratedSigningKeyToSigningModel(k *ecdsa.GeneratedSigningKey, now time.Time) models.SigningKey {
	return models.SigningKey{
		KID:           k.KID,
		Algorithm:     k.Algorithm,
		PublicJWK:     datatypes.JSON(k.PublicJWK),
		PrivateKeyEnc: k.PrivateKeyEnc,
		Status:        oauth.SigningKeyStatusActive,
		NotBefore:     now,
	}
}
