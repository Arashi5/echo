package health

import (
	"context"
	"strconv"

	"github.com/getsentry/sentry-go"
)

func NewSentryService(s Service) Service {
	return &sentryService{s}
}

type sentryService struct {
	Service
}

type sentryLog interface {
	SentryLog() []interface{}
}

func (s *sentryService) getSentryLog(req interface{}, resp interface{}) (out map[string][]interface{}) {
	out = make(map[string][]interface{})
	if sentry, ok := interface{}(req).(sentryLog); ok {
		out["request"] = append(out["request"], sentry.SentryLog()...)
	}

	if sentry, ok := interface{}(resp).(sentryLog); ok {
		out["response"] = append(out["response"], sentry.SentryLog()...)
	}
	return
}

func (s *sentryService) Liveness(ctx context.Context, req *LivenessRequest) (resp *LivenessResponse, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "Liveness")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.Liveness(ctx, req)
}

func (s *sentryService) Readiness(ctx context.Context, req *ReadinessRequest) (resp *ReadinessResponse, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "Readiness")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.Readiness(ctx, req)
}

func (s *sentryService) Version(ctx context.Context, req *VersionRequest) (resp *VersionResponse, err error) {
	defer func() {
		if err != nil {
			log := s.getSentryLog(req, resp)
			sentry.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", strconv.Itoa(getHTTPStatusCode(err)))
				scope.SetTag("method", "Version")
				scope.SetExtra("request", log["request"])
				scope.SetExtra("response", log["response"])
			})
			sentry.CaptureException(err)
		}
	}()
	return s.Service.Version(ctx, req)
}
