package health

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	pb "github.com/arashi5/echo/internal/echoexperimentalechoservicepb"
)

// NewGRPCClient returns an Service backed by a gRPC server at the other end
// of the conn. The caller is responsible for constructing the conn, and
// eventually closing the underlying transport. We bake-in certain middlewares,
// implementing the client library pattern.
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) Service {
	// global client middlewares
	options := []grpctransport.ClientOption{
		grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
	}

	return endpoints{
		// Each individual endpoint is an grpc/transport.Client (which implements
		// endpoint.Endpoint) that gets wrapped with various middlewares. If you
		// made your own client library, you'd do this work there, so your server
		// could rely on a consistent set of client behavior.
		LivenessEndpoint: grpctransport.NewClient(
			conn,
			"echoexperimentalechoservicepb.HealthService",
			"Liveness",
			encodeGRPCLivenessRequest,
			decodeGRPCLivenessResponse,
			pb.LivenessResponse{},
			options...,
		).Endpoint(),
		ReadinessEndpoint: grpctransport.NewClient(
			conn,
			"echoexperimentalechoservicepb.HealthService",
			"Readiness",
			encodeGRPCReadinessRequest,
			decodeGRPCReadinessResponse,
			pb.ReadinessResponse{},
			options...,
		).Endpoint(),
		VersionEndpoint: grpctransport.NewClient(
			conn,
			"echoexperimentalechoservicepb.HealthService",
			"Version",
			encodeGRPCVersionRequest,
			decodeGRPCVersionResponse,
			pb.VersionResponse{},
			options...,
		).Endpoint(),
	}
}

func encodeGRPCLivenessRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*LivenessRequest)
	if !ok {
		return nil, errors.New("encodeGRPCLivenessRequest wrong request")
	}

	return LivenessRequestToPB(inReq), nil
}

func encodeGRPCReadinessRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*ReadinessRequest)
	if !ok {
		return nil, errors.New("encodeGRPCReadinessRequest wrong request")
	}

	return ReadinessRequestToPB(inReq), nil
}

func encodeGRPCVersionRequest(_ context.Context, request interface{}) (interface{}, error) {
	inReq, ok := request.(*VersionRequest)
	if !ok {
		return nil, errors.New("encodeGRPCVersionRequest wrong request")
	}

	return VersionRequestToPB(inReq), nil
}

func decodeGRPCLivenessResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.LivenessResponse)
	if !ok {
		return nil, errors.New("decodeGRPCLivenessResponse wrong response")
	}

	resp := PBToLivenessResponse(inResp)

	return *resp, nil
}

func decodeGRPCReadinessResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.ReadinessResponse)
	if !ok {
		return nil, errors.New("decodeGRPCReadinessResponse wrong response")
	}

	resp := PBToReadinessResponse(inResp)

	return *resp, nil
}

func decodeGRPCVersionResponse(_ context.Context, response interface{}) (interface{}, error) {
	inResp, ok := response.(*pb.VersionResponse)
	if !ok {
		return nil, errors.New("decodeGRPCVersionResponse wrong response")
	}

	resp := PBToVersionResponse(inResp)

	return *resp, nil
}
