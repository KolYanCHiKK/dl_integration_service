package commands

import (
	"context"
	"errors"

	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/base64"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/utils"
	"golang.org/x/sync/errgroup"
)

type CreateUserCommand struct {
	Postgres *postgres.Postgres
}

func NewCreateUserCommand(postgres *postgres.Postgres) *CreateUserCommand {
	return &CreateUserCommand{Postgres: postgres}
}

func (c *CreateUserCommand) CreateUser(ctx context.Context, params *mappings.CreateUserParams) (*mappings.CreateUserResponseParams, error) {
	var (
		err         error
		loginExists bool
		emailExists bool

		checkHigherExists     bool
		checkDepartmentExists bool
	)

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		exists, err := c.Postgres.IsLoginExists(gctx, params.Login)
		if err != nil {
			return err
		}
		loginExists = exists
		return nil
	})
	g.Go(func() error {
		exists, err := c.Postgres.IsEmailExists(gctx, params.Email)
		if err != nil {
			return err
		}
		emailExists = exists
		return nil
	})
	g.Go(func() error {
		exists, err := c.Postgres.CheckHigherId(gctx, params.HigherId)
		if err != nil {
			return err
		}

		checkHigherExists = exists

		return nil
	})
	g.Go(func() error {
		exists, err := c.Postgres.CheckDepartmentId(gctx, params.DepartmentId)
		if err != nil {
			return err
		}

		checkDepartmentExists = exists

		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, error2.MapUserCrudErrors(err)
	}

	if loginExists {
		return nil, error2.LoginAlreadyExists
	}
	if emailExists {
		return nil, error2.EmailAlreadyExists
	}
	if !checkHigherExists {
		return nil, error2.HigherIsNotExists
	}
	if !checkDepartmentExists {
		return nil, error2.DepartmentIsNotExists
	}

	hashedPassword, err := base64.CreateHashedPassword(params.Password)
	if errors.Is(err, base64.ErrorPasswordIsTooLong) {
		return nil, error2.PasswordIsTooLong
	}
	if err != nil {
		return nil, err
	}

	normalizedPhone, err := utils.MapPhoneToStandardPattern(params.Phone)
	if err != nil {
		return nil, error2.MapUserCrudErrors(err)
	}

	newUser := &models.User{
		Surname:        params.Surname,
		Name:           params.Name,
		LastName:       params.LastName,
		Phone:          normalizedPhone,
		BirthDate:      params.BirthDate,
		Login:          params.Login,
		Email:          params.Email,
		Higher:         params.HigherId,
		DepartmentID:   params.DepartmentId,
		HashedPassword: hashedPassword,
	}

	newUser, err = c.Postgres.CreateUser(ctx, newUser)
	if err != nil {
		return nil, error2.MapUserCrudErrors(err)
	}

	return mappings.MapResponseParamsToCreateUserResponseParams(newUser), nil
}
