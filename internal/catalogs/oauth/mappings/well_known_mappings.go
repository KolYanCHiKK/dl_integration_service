package mappings

import (
	"encoding/json"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/jwk/dtos"
	dtos2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/features/openid/dtos"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
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

type OpenIDResponseParams struct {
	Issuer                string
	AuthorizationEndpoint string
	TokenEndpoint         string
	JWKSURI               string
	EndSessionEndpoint    string

	ResponseTypesSupported            []oauth.ResponseType
	GrantTypesSupported               []oauth.GrantType
	SubjectTypesSupported             []app.SubjectType
	IDTokenSigningAlgValuesSupported  []app.SigningAlg
	TokenEndpointAuthMethodsSupported []oauth.AuthMethod
	CodeChallengeMethodsSupported     []app.CodeChallengeAlg
	ScopesSupported                   []oauth.ScopeName
}

func MapOpenIDResponseParamsToDTO(params *OpenIDResponseParams) *dtos2.OpenIDResponseDTO {
	return &dtos2.OpenIDResponseDTO{
		Issuer:                            params.Issuer,
		AuthorizationEndpoint:             params.AuthorizationEndpoint,
		TokenEndpoint:                     params.TokenEndpoint,
		JWKSURI:                           params.JWKSURI,
		EndSessionEndpoint:                params.EndSessionEndpoint,
		ResponseTypesSupported:            params.ResponseTypesSupported,
		GrantTypesSupported:               params.GrantTypesSupported,
		SubjectTypesSupported:             params.SubjectTypesSupported,
		IDTokenSigningAlgValuesSupported:  params.IDTokenSigningAlgValuesSupported,
		TokenEndpointAuthMethodsSupported: params.TokenEndpointAuthMethodsSupported,
		CodeChallengeMethodsSupported:     params.CodeChallengeMethodsSupported,
		ScopesSupported:                   params.ScopesSupported,
	}
}
