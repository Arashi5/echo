package limiting

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/valyala/fasthttp"
	"golang.org/x/time/rate"

	"github.com/arashi5/echo/tools/logging"
)

type Limiter interface {
	Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler
}

func NewLimiter(ctx context.Context, limit float64) Limiter {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "component", "limiter")
	return &limiter{
		limiter: rate.NewLimiter(rate.Limit(limit), 1),
		logger:  logger,
	}
}

type limiter struct {
	limiter *rate.Limiter
	logger  log.Logger
}

func (l *limiter) Middleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if l.limiter.Allow() == false {
			code := fasthttp.StatusTooManyRequests
			msg := fasthttp.StatusMessage(code)

			level.Debug(l.logger).Log("code", code, "msg", msg)
			ctx.Error(msg, code)

			return
		}
		next(ctx)
	}
}
