package commands

import (
	"encoding/json"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/utils"
	"github.com/ory/fosite"
	"gorm.io/datatypes"
)

type ClientParametersMapping struct {
	ID                     oauth.ClientId     `json:"clientId"`
	Name                   string             `json:"name"`
	Status                 oauth.ClientStatus `json:"status"`
	Owner                  string             `json:"owner"`
	SecretHash             string             `json:"secretHash"`
	Public                 bool               `json:"public"`
	AuthMethod             string             `json:"authMethod"`
	AuthSigningAlg         string             `json:"authSigningAlg"`
	JWKS                   datatypes.JSON     `json:"jwks"`
	JWKSURI                string             `json:"jwksuri"`
	AccessTokenTTLSeconds  *int               `json:"accessTokenTTLSeconds"`
	RefreshTokenTTLSeconds *int               `json:"refreshTokenTTLSeconds"`
}

func MapClientToCommand(client *models.Client) *ClientParametersMapping {
	return &ClientParametersMapping{
		ID:                     client.ID,
		Name:                   client.Name,
		Status:                 client.Status,
		Owner:                  client.Owner,
		SecretHash:             client.SecretHash,
		Public:                 client.Public,
		AuthMethod:             client.AuthMethod,
		AuthSigningAlg:         client.AuthSigningAlg,
		JWKS:                   client.JWKS,
		JWKSURI:                client.JWKSURI,
		AccessTokenTTLSeconds:  client.AccessTokenTTLSeconds,
		RefreshTokenTTLSeconds: client.RefreshTokenTTLSeconds,
	}
}

type ClientGrantTypeParametersMapping struct {
	ID        oauth.ClientGrantTypeId `json:"clientGrantTypeId"`
	ClientID  oauth.ClientId          `json:"clientId"`
	GrantType oauth.GrantType         `json:"grantType"`
}

func MapClientGrantTypeToCommand(clientGrantType models.ClientGrantType) *ClientGrantTypeParametersMapping {
	return &ClientGrantTypeParametersMapping{
		ID:        clientGrantType.ID,
		ClientID:  clientGrantType.ClientID,
		GrantType: clientGrantType.GrantType,
	}
}

type FositeRequestToDatabaseRequest struct {
	FositeRequestID string
	ClientID        string
	Subject         string
	RequestedAt     time.Time
	Form            datatypes.JSON
	Session         []byte
	RequestScope    []string
	GrantedScope    []string
	RequestAudience []string
	GrantedAudience []string
}

func MapFositeRequestToDatabaseRequest(fr fosite.Requester) (*FositeRequestToDatabaseRequest, error) {
	jsonForm, err := json.Marshal(fr.GetRequestForm())
	if err != nil {
		return nil, err
	}
	sessionByte, err := json.Marshal(fr.GetSession())
	if err != nil {
		return nil, err
	}

	requestScope := utils.Uniq(fr.GetRequestedScopes())
	grantedScope := utils.Uniq(fr.GetGrantedScopes())
	requestAudience := utils.Uniq(fr.GetRequestedAudience())
	grantedAudience := utils.Uniq(fr.GetGrantedAudience())

	return &FositeRequestToDatabaseRequest{
		FositeRequestID: fr.GetID(),
		ClientID:        fr.GetClient().GetID(),
		Subject:         fr.GetSession().GetSubject(),
		RequestedAt:     fr.GetRequestedAt(),
		Form:            jsonForm,
		Session:         sessionByte,
		RequestScope:    requestScope,
		GrantedScope:    grantedScope,
		RequestAudience: requestAudience,
		GrantedAudience: grantedAudience,
	}, nil
}
