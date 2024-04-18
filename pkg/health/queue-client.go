package health

import (
	"context"
	"encoding/json"

	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

func NewNatsPublisher(ctx context.Context, nc *nats.Conn) Service {
	return endpoints{
		LivenessEndpoint: kitnats.NewPublisher(
			nc,
			"experimental-echo-service-health-liveness",
			kitnats.EncodeJSONRequest,
			decodeHTTPLiveness,
		).Endpoint(),
		ReadinessEndpoint: kitnats.NewPublisher(
			nc,
			"experimental-echo-service-health-readiness",
			kitnats.EncodeJSONRequest,
			decodeHTTPReadiness,
		).Endpoint(),
		VersionEndpoint: kitnats.NewPublisher(
			nc,
			"experimental-echo-service-health-version",
			kitnats.EncodeJSONRequest,
			decodeHTTPVersion,
		).Endpoint(),
	}
}

func decodeHTTPLiveness(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var response LivenessResponse
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeHTTPReadiness(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var response ReadinessResponse
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeHTTPVersion(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var response VersionResponse
	if err := json.Unmarshal(msg.Data, &response); err != nil {
		return nil, err
	}
	return response, nil
}
