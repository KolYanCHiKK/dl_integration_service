package mappings

type LoginPage struct {
	LoginRequestID string
	Nonce          string
	NoticesJSON    string // Для уведомления, которые нужно показать сразу при ошибке открытия
}
