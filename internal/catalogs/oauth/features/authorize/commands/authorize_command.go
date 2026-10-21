package commands

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings/commands"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	oauth2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	logger2 "github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
	"github.com/google/uuid"
	"github.com/ory/fosite"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type AuthorizeCommand struct {
	Postgres *postgres.Postgres
	Cipher   *aead.Cipher
}

func NewAuthorizeCommand(postgres *postgres.Postgres, cipher *aead.Cipher) *AuthorizeCommand {
	return &AuthorizeCommand{
		Postgres: postgres,
		Cipher:   cipher,
	}
}

func (c *AuthorizeCommand) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	clientId, err := uuid.Parse(id)
	if err != nil {
		logger2.CreateErrLog(err.Error(), "[GetClient] Ошибка валидации clientId")
		return nil, fosite.ErrNotFound
	}

	var (
		client        *commands.ClientParametersMapping
		grantTypes    []string
		responseTypes []string
		redirectURIs  []string
		scopes        []string
		audience      []string
		rotated       [][]byte
	)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		cl, err := c.Postgres.GetClient(gctx, clientId)
		if err != nil {
			return err
		}
		client = commands.MapClientToCommand(cl)
		return nil
	})
	g.Go(func() (err error) {
		grantTypes, err = c.Postgres.GetClientGrantType(gctx, clientId)
		return err
	})
	g.Go(func() (err error) {
		responseTypes, err = c.Postgres.GetClientResponseTypes(gctx, clientId)
		return err
	})
	g.Go(func() (err error) {
		redirectURIs, err = c.Postgres.GetClientRedirectURIs(gctx, clientId)
		return err
	})
	g.Go(func() (err error) {
		scopes, err = c.Postgres.GetClientScopes(gctx, clientId)
		return err
	})
	g.Go(func() (err error) {
		audience, err = c.Postgres.GetClientAudience(gctx, clientId)
		return err
	})
	g.Go(func() (err error) {
		rotated, err = c.Postgres.GetClientRotatedSecrets(gctx, clientId)
		return err
	})

	if err = g.Wait(); err != nil {
		return nil, error2.MapClientLoadError(err)
	}

	return &fosite.DefaultClient{
		ID:             clientId.String(),
		Secret:         []byte(client.SecretHash),
		RotatedSecrets: rotated,
		RedirectURIs:   redirectURIs,
		GrantTypes:     grantTypes,
		ResponseTypes:  responseTypes,
		Scopes:         scopes,
		Audience:       audience,
		Public:         client.Public,
	}, nil
}

func (c *AuthorizeCommand) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	return nil
}

func (c *AuthorizeCommand) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return nil
}

func (c *AuthorizeCommand) CreateAuthorizeCodeSession(ctx context.Context, signature string, req fosite.Requester) error {
	if req.GetSession().GetSubject() == "" {
		err := fosite.ErrServerError.WithDebug("authorize code without subject")
		logger2.CreateErrLog(err.Error(), "Отсутствует Subject в запросе /authorize")
		return err
	}
	exp := req.GetSession().GetExpiresAt(fosite.AuthorizeCode)
	if exp.IsZero() {
		err := fosite.ErrServerError.WithDebug("pkce without expiry")
		logger2.CreateErrLog(err.Error(), "Ошибка создания PKCE сессии")
		return err
	}

	mappedRequest, err := commands.MapFositeRequestToDatabaseRequest(req)
	if err != nil {
		return error2.MapSessionError(err)
	}

	err = uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		id, err := c.EnsureRequest(tx, ctx, *mappedRequest)
		if err != nil {
			return err
		}
		return c.Postgres.CreateAuthorizeCode(tx, ctx, *id, signature, req.GetSession().GetExpiresAt(fosite.AuthorizeCode))
	})
	if err != nil {
		return error2.MapSessionError(err)
	}
	return nil
}

func (c *AuthorizeCommand) GetAuthorizeCodeSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	code, err := c.Postgres.GetAuthorizeCode(ctx, signature)
	if err != nil {
		return nil, error2.MapSessionError(err)
	}
	req, err := c.LoadRequester(ctx, code.RequestID, session)
	if err != nil {
		return nil, error2.MapSessionError(err)
	}
	req.GetSession().SetExpiresAt(fosite.AuthorizeCode, code.ExpiresAt) // Установливаем срок жизни токена
	if !code.Active {
		return req, fosite.ErrInvalidatedAuthorizeCode // Отзываем токен если он уже был использован
	}
	return req, nil
}

func (c *AuthorizeCommand) InvalidateAuthorizeCodeSession(ctx context.Context, signature string) error {
	err := c.Postgres.DeactivateAuthorizeCode(ctx, signature)
	if err != nil {
		return error2.MapSessionError(err)
	}

	return nil
}

func (c *AuthorizeCommand) CreatePKCERequestSession(ctx context.Context, signature string, req fosite.Requester) error {
	exp := req.GetSession().GetExpiresAt(fosite.AuthorizeCode)
	if exp.IsZero() {
		err := fosite.ErrServerError.WithDebug("pkce without expiry")
		logger2.CreateErrLog(err.Error(), "Ошибка создания PKCE сессии")
		return err
	}
	form := req.GetRequestForm()
	challenge := form.Get("code_challenge")
	method := form.Get("code_challenge_method")
	if method == "" {
		method = "plain"
	}

	mapped, err := commands.MapFositeRequestToDatabaseRequest(req)
	if err != nil {
		return error2.MapSessionError(err)
	}
	err = uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		id, err := c.EnsureRequest(tx, ctx, *mapped)
		if err != nil {
			return err
		}
		return c.Postgres.CreatePKCERequest(tx, ctx, *id, challenge, method, exp)
	})
	return error2.MapSessionError(err)
}

func (c *AuthorizeCommand) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	pkce, err := c.Postgres.GetPKCERequest(ctx, signature)
	if err != nil {
		return nil, error2.MapSessionError(err)
	}
	req, err := c.LoadRequester(ctx, pkce.RequestID, session)
	if err != nil {
		return nil, error2.MapSessionError(err)
	}

	req.GetSession().SetExpiresAt(fosite.AuthorizeCode, pkce.ExpiresAt) // Установливаем срок жизни токена
	req.Form.Set("code_challenge", pkce.CodeChallenge)
	req.Form.Set("code_challenge_method", pkce.ChallengeMethod)
	return req, nil
}

func (c *AuthorizeCommand) DeletePKCERequestSession(ctx context.Context, signature string) error {
	err := c.Postgres.DeletePKCERequest(ctx, signature)
	if err != nil {
		return error2.MapSessionError(err)
	}

	return nil
}

// Вспомогательные функции для сборки fosite

func (c *AuthorizeCommand) EnsureRequest(tx *gorm.DB, ctx context.Context, request commands.FositeRequestToDatabaseRequest) (*oauth.AuthRequestId, error) {
	if id, found, err := c.Postgres.GetOauthRequestIDByFositeID(tx, ctx, request.FositeRequestID); err != nil {
		return nil, err
	} else if found {
		return &id, nil
	}

	encSession, err := c.Cipher.Encrypt(request.Session, []byte(request.FositeRequestID))
	if err != nil {
		return nil, err
	}

	id, err := c.Postgres.CreateOauthRequest(
		tx, ctx,
		request.FositeRequestID,
		request.ClientID,
		request.Subject,
		request.RequestedAt,
		request.Form,
		encSession,
	)
	if err != nil {
		return nil, err
	}
	if err := c.Postgres.CreateRequestScopes(tx, ctx, *id, oauth2.RequestScopeKindRequested, request.RequestScope); err != nil {
		return nil, err
	}
	if err := c.Postgres.CreateRequestScopes(tx, ctx, *id, oauth2.RequestScopeKindGranted, request.GrantedScope); err != nil {
		return nil, err
	}
	if err := c.Postgres.CreateRequestResources(tx, ctx, *id, oauth2.RequestResourceKindRequested, request.RequestAudience); err != nil {
		return nil, err
	}
	if err := c.Postgres.CreateRequestResources(tx, ctx, *id, oauth2.RequestResourceKindGranted, request.GrantedAudience); err != nil {
		return nil, err
	}

	return id, nil
}

func (c *AuthorizeCommand) LoadRequester(ctx context.Context, requestID oauth.AuthRequestId, session fosite.Session) (*fosite.Request, error) {
	oauthRequest, err := c.Postgres.GetOauthRequest(ctx, requestID)
	if err != nil {
		return nil, err
	}
	client, err := c.GetClient(ctx, oauthRequest.ClientID)
	if err != nil {
		return nil, err
	}

	decryptSession, err := c.Cipher.Decrypt(oauthRequest.Session, []byte(oauthRequest.FositeRequestID))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(decryptSession, &session); err != nil {
		return nil, err
	}

	form := url.Values{}
	if err := json.Unmarshal(oauthRequest.Form, &form); err != nil {
		return nil, err
	}

	reqScope, err := c.Postgres.GetRequestScopes(ctx, requestID, oauth2.RequestScopeKindRequested)
	if err != nil {
		return nil, err
	}
	grScope, err := c.Postgres.GetRequestScopes(ctx, requestID, oauth2.RequestScopeKindGranted)
	if err != nil {
		return nil, err
	}
	reqAud, err := c.Postgres.GetRequestResources(ctx, requestID, oauth2.RequestResourceKindRequested)
	if err != nil {
		return nil, err
	}
	grAud, err := c.Postgres.GetRequestResources(ctx, requestID, oauth2.RequestResourceKindGranted)
	if err != nil {
		return nil, err
	}

	return &fosite.Request{
		ID:                oauthRequest.FositeRequestID,
		RequestedAt:       oauthRequest.RequestedAt,
		Client:            client,
		RequestedScope:    reqScope,
		GrantedScope:      grScope,
		RequestedAudience: reqAud,
		GrantedAudience:   grAud,
		Form:              form,
		Session:           session,
	}, nil
}
