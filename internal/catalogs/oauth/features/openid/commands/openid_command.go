package commands

import (
	"context"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	app2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/app"
	oauth2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
)

type OpenIDCommand struct {
	Postgres *postgres.Postgres
}

func NewOpenIDCommand(postgres *postgres.Postgres) *OpenIDCommand {
	return &OpenIDCommand{
		Postgres: postgres,
	}
}

var (
	responseTypesSupported = []oauth.ResponseType{
		oauth2.ResponseTypeCode,
	}

	grantTypesSupported = []oauth.GrantType{
		oauth2.GrantTypeAuthorizationCode, oauth2.GrantTypeRefreshToken, oauth2.GrantTypeClientCredentials,
	}

	subjectTypesSupported = []app.SubjectType{
		app2.SubjectTypePublic,
	}

	signingAlgSupported = []app.SigningAlg{
		app2.SigningAlgES256,
	}

	authMethodSupported = []oauth.AuthMethod{
		oauth2.AuthMethodSecretBasic, oauth2.AuthMethodSecretPost,
		oauth2.AuthMethodPrivateKeyJWT, oauth2.AuthMethodNone,
	}

	codeChallengeMethodSupported = []app.CodeChallengeAlg{
		app2.CodeChallengeAlgS256,
	}
)

func (c *OpenIDCommand) GetOpenIDConfigurator(ctx context.Context, issuer string) (*mappings.OpenIDResponseParams, error) {
	scopes, err := c.Postgres.GetAllScopes(ctx)
	if err != nil {
		return nil, error2.MapWellKnowError(err)
	}

	return &mappings.OpenIDResponseParams{
		Issuer:                            issuer,
		AuthorizationEndpoint:             issuer + "/oauth2/authorize",
		TokenEndpoint:                     issuer + "/oauth2/token",
		JWKSURI:                           issuer + "/oauth2/well-known/jwk.json",
		EndSessionEndpoint:                issuer + "/users/logout",
		ResponseTypesSupported:            responseTypesSupported,
		GrantTypesSupported:               grantTypesSupported,
		SubjectTypesSupported:             subjectTypesSupported,
		IDTokenSigningAlgValuesSupported:  signingAlgSupported,
		TokenEndpointAuthMethodsSupported: authMethodSupported,
		CodeChallengeMethodsSupported:     codeChallengeMethodSupported,
		ScopesSupported:                   scopes,
	}, nil
}
