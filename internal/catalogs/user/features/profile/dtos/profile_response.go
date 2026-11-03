package dtos

import (
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/user"
)

type DepartmentDTO struct {
	Id   user.DepartmentId `json:"id"`
	Name string            `json:"name"`
}

type HigherDTO struct {
	Id             user.UserId `json:"id"`
	NormalizedName string      `json:"normalizedName"`
}

type UserDTO struct {
	ID        user.UserId `json:"id"`
	Surname   string      `json:"surname,omitempty"`
	Name      string      `json:"name"`
	LastName  string      `json:"lastName,omitempty"`
	Phone     string      `json:"phone,omitempty"`
	BirthDate *string     `json:"birthDate,omitempty"`
	Login     string      `json:"login"`
	Email     string      `json:"email"`
}

type GetUserProfileResponse struct {
	User       *UserDTO       `json:"user"`
	Department *DepartmentDTO `json:"department,omitempty"`
	Higher     *HigherDTO     `json:"higher,omitempty"`
}
