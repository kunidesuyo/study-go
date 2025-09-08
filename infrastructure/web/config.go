package web

import (
	"strings"

	"go-api-arch-clean-template/pkg"
)

type Config struct {
	Host             string
	Port             string
	CorsAllowOrigins []string
}

func NewConfigWeb() *Config {
	return &Config{
		Host: pkg.GetEnvDefault("WEB_HOST", "0.0.0.0"),
		Port: pkg.GetEnvDefault("WEB_PORT", "8080"),
		CorsAllowOrigins: strings.Split(pkg.GetEnvDefault(
			"WEB_CORS_ALLOW_ORIGINS",
			"http://0.0.0.0:8001"), ","),
	}
}
