package mappings

import (
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/create/dtos"
	dtos2 "github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/features/profile/dtos"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/user/models"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/app"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/utils"
	"gorm.io/datatypes"
)

type CreateUserParams struct {
	Surname      string
	Name         string
	LastName     string
	Phone        string
	BirthDate    *datatypes.Date
	Login        string
	Email        string
	DepartmentId *user.DepartmentId
	HigherId     *user.UserId
	Password     string
}

type CreateUserResponseParams struct {
	ID           user.UserId
	Surname      string
	Name         string
	LastName     string
	Phone        string
	BirthDate    *datatypes.Date
	DepartmentId *user.DepartmentId
	Higher       *user.UserId
	Login        string
	Email        string
}

func MapResponseParamsToCreateUserResponseParams(user *models.User) *CreateUserResponseParams {
	return &CreateUserResponseParams{
		ID:           user.ID,
		Surname:      user.Surname,
		Name:         user.Name,
		LastName:     user.LastName,
		Phone:        user.Phone,
		BirthDate:    user.BirthDate,
		DepartmentId: user.DepartmentID,
		Higher:       user.Higher,
		Login:        user.Login,
		Email:        user.Email,
	}
}

func MapCreateUserResponseParamsToDTO(createParams *CreateUserResponseParams) *dtos.CreateUserResponse {
	if createParams == nil {
		return nil
	}

	birthDate := time.Time(*createParams.BirthDate).Format(time.DateOnly)

	return &dtos.CreateUserResponse{
		Status: app.SuccessCreatedStatus,
		Data: &dtos.CreateUserData{
			ID:           createParams.ID,
			Surname:      createParams.Surname,
			Name:         createParams.Name,
			LastName:     createParams.LastName,
			Phone:        createParams.Phone,
			BirthDate:    &birthDate,
			Higher:       createParams.Higher,
			DepartmentId: createParams.DepartmentId,
			Login:        createParams.Login,
			Email:        createParams.Email,
		},
	}
}

type GetUserProfileResponseParams struct {
	ID               user.UserId
	Surname          string
	Name             string
	LastName         string
	Phone            string
	BirthDate        *datatypes.Date
	DepartmentParams *DepartmentParams
	HigherParams     *HigherParams
	Login            string
	Email            string
}

func MapResponseParamsToGetUserProfileResponseParams(user *models.User, getProfileParams *GetUserProfileResponseParams) *GetUserProfileResponseParams {
	if getProfileParams == nil {
		return nil
	}

	getProfileParams.ID = user.ID
	getProfileParams.Surname = user.Surname
	getProfileParams.Name = user.Name
	getProfileParams.LastName = user.LastName
	getProfileParams.Phone = user.Phone
	getProfileParams.BirthDate = user.BirthDate
	getProfileParams.Login = user.Login
	getProfileParams.Email = user.Email

	return getProfileParams
}

func MapGetUserProfileResponseParamsToDTO(getProfileParams *GetUserProfileResponseParams) *dtos2.GetUserProfileResponse {
	if getProfileParams == nil {
		return nil
	}

	birthDate := time.Time(*getProfileParams.BirthDate).Format(time.DateOnly)
	userDTO := &dtos2.UserDTO{
		ID:        getProfileParams.ID,
		Surname:   getProfileParams.Surname,
		Name:      getProfileParams.Name,
		LastName:  getProfileParams.LastName,
		Phone:     getProfileParams.Phone,
		BirthDate: &birthDate,
		Login:     getProfileParams.Login,
		Email:     getProfileParams.Email,
	}

	return &dtos2.GetUserProfileResponse{
		User:       userDTO,
		Higher:     MapHigherParamsToHigherDTO(getProfileParams.HigherParams),
		Department: MapDepartmentParamsToDepartmentDTO(getProfileParams.DepartmentParams),
	}
}

type DepartmentParams struct {
	Id   user.DepartmentId
	Name string
}

func MapDepartmentParamsToDepartmentDTO(params *DepartmentParams) *dtos2.DepartmentDTO {
	if params == nil {
		return nil
	}

	return &dtos2.DepartmentDTO{
		Id:   params.Id,
		Name: params.Name,
	}
}

type HigherParams struct {
	Id       user.UserId
	Surname  string
	Name     string
	LastName string
}

func MapHigherParamsToHigherDTO(params *HigherParams) *dtos2.HigherDTO {
	if params == nil {
		return nil
	}

	return &dtos2.HigherDTO{
		Id:             params.Id,
		NormalizedName: utils.ShortFullName(params.Surname, params.Name, params.LastName),
	}
}
