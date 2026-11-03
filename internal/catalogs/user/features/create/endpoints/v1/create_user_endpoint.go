package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/connection/interfaces"
	error2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/error"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/dtos"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/body/rest/request"
	rest2 "github.com/KolYanCHiKK/dl_integration_service/pkg/response/rest"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CreateUserEndpoint struct {
	*CreateUserHandler
}

func NewCreateUserEndpoint(params *CreateUserHandler) interfaces.Endpoint {
	return &CreateUserEndpoint{params}
}

func (ep *CreateUserEndpoint) MapEndpoint() {
	ep.Mux.Post("/users", ep.handle())
}

func (ep *CreateUserEndpoint) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.DecodeRestRequest[dtos.CreateUserRequest](w, r)
		if err != nil {
			rest2.GenerateBadRequestResponse(w, nil, err.Error())
			return
		}

		validationErrs, err := ep.Dependents.Validator.ValidateBody(body, error2.ValidationCreateUserErrorsMap)
		if err != nil {
			rest2.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}
		if validationErrs != nil {
			rest2.GenerateBadRequestResponse(w, nil, "", validationErrs...)
			return
		}

		// Парсим строку, чтобы получить дату без времени
		var birthDate *datatypes.Date
		if body.BirthDate != "" {
			parsedDate, _ := time.Parse(time.DateOnly, body.BirthDate) // Ошибка отбрасывается валидатором
			date := datatypes.Date(parsedDate)
			birthDate = &date
		}

		// Парсим строку, чтобы получить идентификатор департамента
		var departmentId *user.DepartmentId
		if body.DepartmentId != nil {
			departmentIdParse, _ := uuid.Parse(*body.DepartmentId) // Ошибка отбрасывается валидатором
			departmentId = &departmentIdParse
		}

		// Парсим строку, чтобы получить идентификатор руководителя
		var higherId *user.UserId
		if body.HigherId != nil {
			higherIdParse, _ := uuid.Parse(*body.HigherId) // Ошибка отбасывается валидтором
			higherId = &higherIdParse
		}

		createdUser, err := ep.Dependents.Command.CreateUser(r.Context(), &mappings.CreateUserParams{
			Surname:      body.Surname,
			Name:         body.Name,
			LastName:     body.LastName,
			Phone:        body.Phone,
			BirthDate:    birthDate,
			Login:        body.Login,
			Email:        body.Email,
			DepartmentId: departmentId,
			HigherId:     higherId,
			Password:     body.Password,
		})

		switch {

		case errors.Is(err, error2.LoginAlreadyExists):
			rest2.GenerateConflictResponse(w, nil, "Пользователь с таким логином уже существует", err.Error())
			return

		case errors.Is(err, error2.EmailAlreadyExists):
			rest2.GenerateConflictResponse(w, nil, "Пользователь с таким email уже существует", err.Error())
			return

		case errors.Is(err, error2.PasswordIsTooLong):
			rest2.GenerateBadRequestResponse(w, nil, "Пароль слишком длинный", err.Error())
			return

		case errors.Is(err, error2.DepartmentIsNotExists):
			rest2.GenerateBadRequestResponse(w, nil, "Не найден идентификатор департамента", err.Error())
			return

		case errors.Is(err, error2.HigherIsNotExists):
			rest2.GenerateBadRequestResponse(w, nil, "Не найден идентификатор руководителя", err.Error())
			return

		case err != nil:
			rest2.GenerateInternalServerErrorResponse(w, nil, err.Error())
			return
		}

		rest2.BuildRestResponse(w, mappings.MapCreateUserResponseParamsToDTO(createdUser), 201, nil)

	}
}
