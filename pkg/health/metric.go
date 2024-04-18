package health

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"

	tool "github.com/arashi5/echo/tools/metrics"
)

// NewMetricService returns an instance of an instrumenting Service.
func NewMetricsService(ctx context.Context, s Service) Service {
	counter, latency := tool.FromContext(ctx)
	return &metricService{counter, latency, s}
}

type metricService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *metricService) Liveness(ctx context.Context, req *LivenessRequest) (resp *LivenessResponse, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "health", "handler", "Liveness", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "health", "handler", "Liveness", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.Liveness(ctx, req)
}

func (s *metricService) Readiness(ctx context.Context, req *ReadinessRequest) (resp *ReadinessResponse, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "health", "handler", "Readiness", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "health", "handler", "Readiness", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.Readiness(ctx, req)
}

func (s *metricService) Version(ctx context.Context, req *VersionRequest) (resp *VersionResponse, err error) {
	defer func(begin time.Time) {
		go func() {
			s.requestCount.With("service", "health", "handler", "Version", "code", strconv.Itoa(getHTTPStatusCode(err))).Add(1)
			s.requestLatency.With("service", "health", "handler", "Version", "code", strconv.Itoa(getHTTPStatusCode(err))).Observe(time.Since(begin).Seconds())
		}()
	}(time.Now())
	return s.Service.Version(ctx, req)
}
