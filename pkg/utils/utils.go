package utils

import (
	"net"
	"net/http"

	"github.com/KolYanCHiKK/dl_integration_service/pkg/logger"
)

func Uniq(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := in[:0:0] // пустой срез того же типа
	for _, v := range in {
		if _, ok := seen[v]; !ok { // такого значения ещё не было
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func GetClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		logger.CreateErrLog(err.Error(), "Ошибка получения IP запроса")
		return ""
	}
	return host
}
