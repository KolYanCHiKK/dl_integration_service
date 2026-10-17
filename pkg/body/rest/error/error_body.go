package error

type ErrorBody struct {
	Bcode       int      `json:"bcode"`
	Description string   `json:"description"`
	Errors      []string `json:"errors"`
}

var (
	BadRequestBody = ErrorBody{
		Bcode:       800,
		Description: "Не пройден форматно-логический контроль (ФЛК)",
		Errors:      nil,
	}
)
