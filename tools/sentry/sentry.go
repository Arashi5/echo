package sentry

import (
	"github.com/getsentry/sentry-go"
	sentryfasthttp "github.com/getsentry/sentry-go/fasthttp"
	json "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"

	"github.com/arashi5/echo/configs"
)

// NewSentry initializes new global sentry Client.
func NewSentry(cfg *configs.Config) error {
	var debug bool
	if cfg.Logger.Level == "debug" {
		debug = true
	}
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.Sentry.Dsn,
		Debug:            debug,
		AttachStacktrace: true,
		Environment:      cfg.Sentry.Environment,
	}); err != nil {
		return err
	}
	var data map[string]interface{}
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "config",
		Data:     data,
		Message:  "init config",
	})
	return nil
}

func Middleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return sentryfasthttp.New(sentryfasthttp.Options{Repanic: true}).Handle(h)
}
