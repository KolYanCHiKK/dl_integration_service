package dtos

type CreateUserRequest struct {
	Surname      string  `json:"surname" validate:"max=100"`
	Name         string  `json:"name" validate:"required,max=100"`
	LastName     string  `json:"lastName" validate:"max=100"`
	Phone        string  `json:"phone" validate:"omitempty,max=65,phone"`
	BirthDate    string  `json:"birthDate" validate:"omitempty,datetime=2006-01-02"`
	Login        string  `json:"login" validate:"required,max=120"`
	Email        string  `json:"email" validate:"required,max=120,email"`
	DepartmentId *string `json:"departmentId" validate:"omitempty,uuid"`
	HigherId     *string `json:"higherId" validate:"omitempty,uuid"`
	Password     string  `json:"password" validate:"required,max=72"`
}
