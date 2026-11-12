package utils

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"
	"unicode"
	"unicode/utf8"
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

func NonNil(v []string) []string {
	if v == nil {
		return []string{}
	}
	return v
}

func GetClientIP(r *http.Request, trustedProxies []netip.Prefix) string {
	remote, err := netip.ParseAddrPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	addr := remote.Addr().Unmap()

	if isTrustedProxies(addr, trustedProxies) {
		if h := r.Header.Get("X-Real-IP"); h != "" {
			if ip, err := netip.ParseAddr(h); err == nil {
				return ip.String()
			}
		}
	}

	return addr.String()
}

func isTrustedProxies(addr netip.Addr, trusted []netip.Prefix) bool {
	for _, p := range trusted {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}

func ShortFullName(surname, name, patronymic string) string {
	surname = strings.TrimSpace(surname)
	name = strings.TrimSpace(name)
	patronymic = strings.TrimSpace(patronymic)

	parts := make([]string, 0, 3)

	if surname != "" {
		parts = append(parts, surname, initial(name), initial(patronymic))
	} else {
		parts = append(parts, name, patronymic)
	}

	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}

	return strings.Join(result, " ")
}

func initial(s string) string {
	if s == "" {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + "."
}

func MapPhoneToStandardPattern(phone string) (string, error) {
	runePhone := []rune(phone)
	if string(runePhone[0]) == "+" {
		return string(runePhone), nil
	} else if string(runePhone[0]) != "8" {
		return "", errors.New("error phone number validation")
	}

	var rightPhone []rune
	rightPhone = append(rightPhone, []rune("+7")...)
	for i := 1; i < len(runePhone); i++ {
		rightPhone = append(rightPhone, runePhone[i])
	}

	return string(rightPhone), nil
}
