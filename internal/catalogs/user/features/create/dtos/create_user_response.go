package dtos

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
)

type CreateUserData struct {
	ID           user.UserId        `json:"id"`
	Surname      string             `json:"surname,omitempty"`
	Name         string             `json:"name"`
	LastName     string             `json:"lastName,omitempty"`
	Phone        string             `json:"phone,omitempty"`
	BirthDate    *string            `json:"birthDate,omitempty"`
	Higher       *user.UserId       `json:"higher,omitempty"`
	DepartmentId *user.DepartmentId `json:"departmentId,omitempty"`
	Login        string             `json:"login"`
	Email        string             `json:"email"`
}

type CreateUserResponse struct {
	Status string          `json:"status"`
	Data   *CreateUserData `json:"data"`
}
