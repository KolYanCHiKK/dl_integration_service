package commands

import (
	"context"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/oauth"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/uow"
	"github.com/ory/fosite"
	"gorm.io/gorm"
)

type AuthorizeUtils interface {
	LoadRequester(ctx context.Context, requestID oauth.AuthRequestId, session fosite.Session) (*fosite.Request, error)
}

type TokenCommand struct {
	Postgres       *postgres.Postgres
	AuthorizeUtils AuthorizeUtils
}

func NewTokenCommand(postgres *postgres.Postgres, utilCommand AuthorizeUtils) *TokenCommand {
	return &TokenCommand{
		Postgres:       postgres,
		AuthorizeUtils: utilCommand,
	}
}

func (c *TokenCommand) CreateAccessTokenSession(ctx context.Context, signature string, req fosite.Requester) error {
	expiresAt := req.GetSession().GetExpiresAt(fosite.AccessToken)
	if expiresAt.IsZero() {
		err := fosite.ErrServerError.WithDebug("access token expiry is not set")
		logger.CreateErrLog(err.Error(), "Access Token имеет нулевой срок жизни")
		return err
	}

	err := uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		requestID, isFind, err := c.Postgres.GetOauthRequestIDByFositeID(tx, ctx, req.GetID())
		if err != nil {
			return err
		}
		if !isFind {
			return error2.OauthRequestNotFound
		}

		err = c.Postgres.CreateAccessToken(tx, ctx, requestID, signature, expiresAt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	token, err := c.Postgres.GetAccessToken(ctx, signature)
	if err != nil {
		return nil, error2.MapTokenError(err)
	}

	session.SetExpiresAt(fosite.AccessToken, token.ExpiresAt)

	req, err := c.AuthorizeUtils.LoadRequester(ctx, token.RequestID, session)
	if err != nil {
		return nil, error2.MapTokenError(err)
	}
	if !token.Active {
		return req, fosite.ErrInactiveToken
	}

	return req, nil
}

func (c *TokenCommand) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	err := c.Postgres.DeactivateAccessToken(ctx, signature)
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) CreateRefreshTokenSession(ctx context.Context, signature string, accessSignature string, req fosite.Requester) error {
	expiresAt := req.GetSession().GetExpiresAt(fosite.RefreshToken)
	if expiresAt.IsZero() {
		err := fosite.ErrServerError.WithDebug("refresh token expiry is not set")
		logger.CreateErrLog(err.Error(), "Refresh Token имеет нулевой срок жизни")
		return err
	}

	err := uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		requestID, isFind, err := c.Postgres.GetOauthRequestIDByFositeID(tx, ctx, req.GetID())
		if err != nil {
			return err
		}
		if !isFind {
			return error2.OauthRequestNotFound
		}

		err = c.Postgres.CreateRefreshToken(tx, ctx, requestID, signature, accessSignature, expiresAt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	token, err := c.Postgres.GetRefreshToken(ctx, signature)
	if err != nil {
		return nil, error2.MapTokenError(err)
	}

	session.SetExpiresAt(fosite.RefreshToken, token.ExpiresAt)

	req, err := c.AuthorizeUtils.LoadRequester(ctx, token.RequestID, session)
	if err != nil {
		return nil, error2.MapTokenError(err)
	}
	if !token.Active {
		return req, fosite.ErrInactiveToken
	}

	return req, nil
}

func (c *TokenCommand) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	err := c.Postgres.DeactivateRefreshToken(ctx, signature)
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) error {
	err := uow.ExecuteTransaction(ctx, c.Postgres, func(tx *gorm.DB) error {
		err := c.Postgres.RotateRefreshToken(tx, ctx, requestID, refreshTokenSignature)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) RevokeAccessToken(ctx context.Context, requestID string) error {
	err := c.Postgres.RevokeAccessTokensByRequest(ctx, requestID)
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}

func (c *TokenCommand) RevokeRefreshToken(ctx context.Context, requestID string) error {
	err := c.Postgres.RevokeRefreshTokensByRequest(ctx, requestID)
	if err != nil {
		return error2.MapTokenError(err)
	}

	return nil
}
