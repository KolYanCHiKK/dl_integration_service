package commands

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/frontend"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/mappings"
	"github.com/KolYanCHiKK/dl_integration_service/internal/catalogs/oauth/repository/postgres"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/aead"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/crypt/base64"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
	"github.com/KolYanCHiKK/dl_integration_service/pkg/notice/auth"
)

var loginTmpl = template.Must(template.ParseFS(frontend.Templates, "auth-page/form.html"))

type LoginCommand struct{}

func NewAuthorizeCommand(postgres *postgres.Postgres, cipher *aead.Cipher) *LoginCommand {
	return &LoginCommand{}
}

func (c *LoginCommand) RenderLogin(w http.ResponseWriter, p mappings.LoginPage, notices []auth.Notice, status int) {
	nonce, err := base64.RandomNonce(16)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка рендера формы авторизации")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if notices == nil {
		notices = []auth.Notice{}
	}

	nj, err := json.Marshal(notices)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка рендера формы авторизации")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	p.Nonce, p.NoticesJSON = nonce, string(nj)

	// Рендерим в буфер: если шаблон упадёт, клиент не получит половину страницы
	var buf bytes.Buffer
	if err := loginTmpl.Execute(&buf, p); err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка рендера формы авторизации")
		slog.Error("render login page", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Cache-Control", "no-store")
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy",
		"default-src 'none'; "+
			"style-src 'nonce-"+nonce+"' https://fonts.googleapis.com; "+
			"font-src https://fonts.gstatic.com; "+
			"script-src 'nonce-"+nonce+"'; "+
			"connect-src 'self'; "+ // fetch на POST /users/login
			"form-action 'self'; frame-ancestors 'none'; base-uri 'none'")

	w.WriteHeader(status)

	_, err = w.Write(buf.Bytes())
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка рендера формы авторизации")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
