package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	app2 "github.com/KolYanCHiKK/dl_integration_service/internal/shared/constants/app"
	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"
)

type envReader struct {
	err error
}

func (r *envReader) fail(err error) {
	if r.err == nil {
		r.err = err
	}
}

func (r *envReader) str(key string) string {
	v := os.Getenv(key)
	if v == "" {
		r.fail(fmt.Errorf("variable %s is missing", key))
		return v
	}
	return v
}

func (r *envReader) sslRootSert(key string, sslmode app.SSLMode) string {
	v := os.Getenv(key)
	if v == "" && sslmode != app2.SSLModeDisable {
		r.fail(fmt.Errorf("variable %s is missing", key))
		return v
	}
	return v
}

func (r *envReader) int(key string) int {
	v := os.Getenv(key)
	if v == "" {
		r.fail(fmt.Errorf("variable %s is missing", key))
		return 0
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		r.fail(fmt.Errorf("variable %s: invalid int %q", key, v))
		return 0
	}
	return n
}

func (r *envReader) duration(key string) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		r.fail(fmt.Errorf("variable %s is missing", key))
		return 0
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		r.fail(fmt.Errorf("variable %s: invalid duration %q", key, v))
		return 0
	}
	return d
}
