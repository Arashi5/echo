package health

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/grpc"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/arashi5/echo/internal/echoexperimentalechoservicepb"
	"github.com/arashi5/echo/tools/logging"
)

type grpcServer struct {
	liveness  grpctransport.Handler
	readiness grpctransport.Handler
	version   grpctransport.Handler
}

type ContextGRPCKey struct{}

type GRPCInfo struct{}

// NewGRPCServer makes a set of endpoints available as a gRPC health Server.
func NewGRPCServer(ctx context.Context, s Service) pb.HealthServiceServer {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "grpc handler", "health")
	//tracer := tracing.FromContext(ctx)

	options := []grpctransport.ServerOption{
		// grpctransport.ServerErrorLogger(logger),
		grpctransport.ServerBefore(grpcToContext()),
		//grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "grpc server", logger)),
		//grpctransport.ServerFinalizer(closeGRPCTracer()),
	}

	return &grpcServer{
		liveness: grpctransport.NewServer(
			makeLivenessEndpoint(s),
			decodeGRPCLivenessRequest,
			encodeGRPCLivenessResponse,
			options...,
		),
		readiness: grpctransport.NewServer(
			makeReadinessEndpoint(s),
			decodeGRPCReadinessRequest,
			encodeGRPCReadinessResponse,
			options...,
		),
		version: grpctransport.NewServer(
			makeVersionEndpoint(s),
			decodeGRPCVersionRequest,
			encodeGRPCVersionResponse,
			options...,
		),
	}
}

func JoinGRPC(ctx context.Context, s Service) func(*googlegrpc.Server) {
	return func(g *googlegrpc.Server) {
		pb.RegisterHealthServiceServer(g, NewGRPCServer(ctx, s))
	}
}

func grpcToContext() grpc.ServerRequestFunc {
	return func(ctx context.Context, md metadata.MD) context.Context {
		return context.WithValue(ctx, ContextGRPCKey{}, GRPCInfo{})
	}
}
func closeGRPCTracer() grpc.ServerFinalizerFunc {
	return func(ctx context.Context, err error) {
		span := stdopentracing.SpanFromContext(ctx)
		span.Finish()
	}
}

func (s *grpcServer) Liveness(ctx context.Context, req *pb.LivenessRequest) (*pb.LivenessResponse, error) {
	_, rep, err := s.liveness.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LivenessResponse), nil
}

func (s *grpcServer) Readiness(ctx context.Context, req *pb.ReadinessRequest) (*pb.ReadinessResponse, error) {
	_, rep, err := s.readiness.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ReadinessResponse), nil
}

func (s *grpcServer) Version(ctx context.Context, req *pb.VersionRequest) (*pb.VersionResponse, error) {
	_, rep, err := s.version.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.VersionResponse), nil
}

func decodeGRPCLivenessRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.LivenessRequest)
	if !ok {
		return nil, errors.New("decodeGRPCLivenessRequest wrong request")
	}

	req := PBToLivenessRequest(inReq)
	return *req, nil
}

func decodeGRPCReadinessRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.ReadinessRequest)
	if !ok {
		return nil, errors.New("decodeGRPCReadinessRequest wrong request")
	}

	req := PBToReadinessRequest(inReq)
	return *req, nil
}

func decodeGRPCVersionRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*pb.VersionRequest)
	if !ok {
		return nil, errors.New("decodeGRPCVersionRequest wrong request")
	}

	req := PBToVersionRequest(inReq)
	return *req, nil
}

func encodeGRPCLivenessResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*LivenessResponse)
	if !ok {
		return nil, errors.New("encodeGRPCLivenessResponse wrong response")
	}

	return LivenessResponseToPB(inResp), nil
}

func encodeGRPCReadinessResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*ReadinessResponse)
	if !ok {
		return nil, errors.New("encodeGRPCReadinessResponse wrong response")
	}

	return ReadinessResponseToPB(inResp), nil
}

func encodeGRPCVersionResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*VersionResponse)
	if !ok {
		return nil, errors.New("encodeGRPCVersionResponse wrong response")
	}

	return VersionResponseToPB(inResp), nil
}

func LivenessRequestToPB(d *LivenessRequest) *pb.LivenessRequest {
	if d == nil {
		return nil
	}

	resp := pb.LivenessRequest{}

	return &resp
}

func PBToLivenessRequest(d *pb.LivenessRequest) *LivenessRequest {
	if d == nil {
		return nil
	}

	resp := LivenessRequest{}

	return &resp
}

func LivenessResponseToPB(d *LivenessResponse) *pb.LivenessResponse {
	if d == nil {
		return nil
	}

	resp := pb.LivenessResponse{
		Status: d.Status,
	}

	return &resp
}

func PBToLivenessResponse(d *pb.LivenessResponse) *LivenessResponse {
	if d == nil {
		return nil
	}

	resp := LivenessResponse{
		Status: d.Status,
	}

	return &resp
}

func ReadinessRequestToPB(d *ReadinessRequest) *pb.ReadinessRequest {
	if d == nil {
		return nil
	}

	resp := pb.ReadinessRequest{}

	return &resp
}

func PBToReadinessRequest(d *pb.ReadinessRequest) *ReadinessRequest {
	if d == nil {
		return nil
	}

	resp := ReadinessRequest{}

	return &resp
}

func ReadinessResponseToPB(d *ReadinessResponse) *pb.ReadinessResponse {
	if d == nil {
		return nil
	}

	resp := pb.ReadinessResponse{
		Status: d.Status,
	}

	return &resp
}

func PBToReadinessResponse(d *pb.ReadinessResponse) *ReadinessResponse {
	if d == nil {
		return nil
	}

	resp := ReadinessResponse{
		Status: d.Status,
	}

	return &resp
}

func VersionRequestToPB(d *VersionRequest) *pb.VersionRequest {
	if d == nil {
		return nil
	}

	resp := pb.VersionRequest{}

	return &resp
}

func PBToVersionRequest(d *pb.VersionRequest) *VersionRequest {
	if d == nil {
		return nil
	}

	resp := VersionRequest{}

	return &resp
}

func VersionResponseToPB(d *VersionResponse) *pb.VersionResponse {
	if d == nil {
		return nil
	}

	resp := pb.VersionResponse{
		BuildTime: d.BuildTime,
		Commit:    d.Commit,
		Version:   d.Version,
	}

	return &resp
}

func PBToVersionResponse(d *pb.VersionResponse) *VersionResponse {
	if d == nil {
		return nil
	}

	resp := VersionResponse{
		BuildTime: d.BuildTime,
		Commit:    d.Commit,
		Version:   d.Version,
	}

	return &resp
}
