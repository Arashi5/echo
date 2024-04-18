package health

import (
	"context"
)

var (
	Version   = ""
	Commit    = ""
	BuildTime = ""
)

type healthService struct {
}

func NewHealthService() Service {
	return &healthService{}
}

func (s *healthService) Liveness(ctx context.Context, req *LivenessRequest) (resp *LivenessResponse, err error) {
	return &LivenessResponse{
		Status: "ok",
	}, nil
}

func (s *healthService) Readiness(ctx context.Context, req *ReadinessRequest) (resp *ReadinessResponse, err error) {
	return &ReadinessResponse{
		Status: "ok",
	}, nil
}

func (s *healthService) Version(ctx context.Context, req *VersionRequest) (resp *VersionResponse, err error) {
	return &VersionResponse{
		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}, nil
}
