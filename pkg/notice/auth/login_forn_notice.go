package auth

import error2 "github.com/KolYanCHiKK/dl_integration_service/pkg/body/rest/error"

type NoticeKind string

const (
	NoticeError   NoticeKind = "error"
	NoticeSuccess NoticeKind = "success"
)

type Notice struct {
	Kind  NoticeKind `json:"kind"`
	Title string     `json:"title"`
	Text  string     `json:"text"`
}

// По одному уведомлению на каждый элемент errors
func ErrorNotices(e *error2.ErrorBody) []Notice {
	msgs := Messages(e)
	out := make([]Notice, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, Notice{Kind: NoticeError, Title: "Ошибка авторизации", Text: m})
	}
	return out
}

func Messages(e *error2.ErrorBody) []string {
	if len(e.Errors) > 0 {
		return e.Errors
	}
	return []string{"Не удалось выполнить вход"}
}
