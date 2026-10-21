package configs

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"
	"github.com/joho/godotenv"
	"github.com/ory/fosite"
)

type ApplicationConfig struct {
	Postgres    *DBPostgresConfig
	JWTAuth     *AuthJwtConfig
	OauthFosite *FositeAuthConfig
}
type DBPostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string

	SSLMode         app.SSLMode   // Мод использования сертификатов
	SSLRootCert     string        // Ссылка на файл SSL сертификата
	MaxOpenConn     int           // Максимальное количество открытых соединений
	MaxIdleConn     int           // Максимальное количество свободных соединений
	ConnMaxLifetime time.Duration // Максимальный срок жизни одного соединения
	ConnMaxIdleTime time.Duration // Время жизни соединения без использования
	ConnectTimeout  time.Duration // Таймаут подключения
}

type AuthJwtConfig struct {
	SecretKey       string        // Секретный ключ
	AccessTokenTTL  time.Duration // Время жизни токена
	RefreshTokenTTL time.Duration // Время жизни рефреш токена
	Issuer          string        // Система, которая выпустила токен
}

type FositeAuthConfig struct {
	SessionEncKey         string // Ключ шифрования для сохранения сессии
	FositeGlobalSecret    string // Ключе шифрования для кодов авторизации
	AccessTokenIssuer     string
	TokenURL              string
	ScopeStrategy         func(haystack []string, needle string) bool
	AuthorizeCodeLifespan time.Duration
	AccessTokenLifespan   time.Duration
	RefreshTokenLifespan  time.Duration
	EnforcePKCE           bool
}

func (c *DBPostgresConfig) DSN() string {
	q := url.Values{}
	q.Set("sslmode", string(c.SSLMode))

	if c.SSLRootCert != "" {
		q.Set("sslrootcert", c.SSLRootCert)
	}
	if c.ConnectTimeout > 0 {
		q.Set("connect_timeout", strconv.Itoa(int(c.ConnectTimeout.Seconds())))
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     c.Name,
		RawQuery: q.Encode(),
	}

	return u.String()
}

func (c *DBPostgresConfig) RunDbConfigWithEnv() error {
	r := &envReader{}

	cfg := DBPostgresConfig{
		Host:     r.str("DB_HOST"),
		Port:     r.int("DB_PORT"),
		User:     r.str("DB_USER"),
		Password: r.str("DB_PASSWORD"),
		Name:     r.str("DB_NAME"),

		SSLMode:     app.SSLMode(r.str("DB_SSL_MODE")),
		SSLRootCert: r.sslRootSert("DB_SSL_ROOT_CERT", app.SSLMode(r.str("DB_SSL_MODE"))),

		MaxOpenConn:     r.int("DB_MAX_OPEN_CONN"),
		MaxIdleConn:     r.int("DB_MAX_IDLE_CONN"),
		ConnMaxLifetime: r.duration("DB_CONN_MAX_LIFETIME"),
		ConnMaxIdleTime: r.duration("DB_CONN_MAX_IDLE_TIME"),
		ConnectTimeout:  r.duration("DB_CONNECT_TIMEOUT"),
	}

	if r.err != nil {
		return r.err
	}

	*c = cfg
	return nil
}

func (a *AuthJwtConfig) RunAuthConfigWithEnv() error {
	r := &envReader{}

	cfg := AuthJwtConfig{
		SecretKey:       r.str("JWT_SECRET_KEY"),
		AccessTokenTTL:  r.duration("JWT_ACCESS_TTL"),
		RefreshTokenTTL: r.duration("JWT_REFRESH_TTL"),
		Issuer:          r.str("JWT_ISSUER"),
	}

	if r.err != nil {
		return r.err
	}

	*a = cfg
	return nil
}

func (f *FositeAuthConfig) RunFositeConfigWithEnv() error {
	r := &envReader{}

	cfg := FositeAuthConfig{
		SessionEncKey:         r.str("SESSION_ENC_KEY"),
		FositeGlobalSecret:    r.str("FOSITE_GLOBAL_SECRET"),
		AccessTokenIssuer:     r.str("ACCESS_TOKEN_ISSUER"),
		TokenURL:              r.str("TOKEN_URL"),
		ScopeStrategy:         fosite.ExactScopeStrategy,
		AuthorizeCodeLifespan: r.duration("AUTHORIZE_CODE_LIFESPAN"),
		AccessTokenLifespan:   r.duration("JWT_ACCESS_TTL"),
		RefreshTokenLifespan:  r.duration("JWT_REFRESH_TTL"),
		EnforcePKCE:           r.bool("ENFORCE_PKCE"),
	}

	if r.err != nil {
		return r.err
	}

	*f = cfg
	return nil
}

func LoadApplicationConfig() (*ApplicationConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	r := &envReader{}
	postgresCfg := &DBPostgresConfig{}
	authCfg := &AuthJwtConfig{}
	fositeCfg := &FositeAuthConfig{}

	err := postgresCfg.RunDbConfigWithEnv()
	if err != nil {
		return nil, err
	}

	err = authCfg.RunAuthConfigWithEnv()
	if err != nil {
		return nil, err
	}

	err = fositeCfg.RunFositeConfigWithEnv()
	if err != nil {
		return nil, err
	}

	cfg := &ApplicationConfig{
		Postgres:    postgresCfg,
		JWTAuth:     authCfg,
		OauthFosite: fositeCfg,
	}

	if r.err != nil {
		return nil, r.err
	}

	return cfg, nil
}
