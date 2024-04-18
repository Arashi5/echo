package health

import (
	"context"
	"encoding/json"

	"github.com/arashi5/echo/tools/logging"
	"github.com/go-kit/kit/log"
	kitnats "github.com/go-kit/kit/transport/nats"
	"github.com/nats-io/nats.go"
)

func NewNatsSubscriber(ctx context.Context, nc *nats.Conn, svc Service) error {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "nats handler")
	{
		ra := kitnats.NewSubscriber(
			makeLivenessEndpoint(svc),
			decodeNATSLivenessRequestRequest,
			kitnats.EncodeJSONResponse,
			kitnats.SubscriberErrorLogger(logger),
		)
		if _, err := nc.QueueSubscribe("experimental-echo-service-health-liveness", "health", ra.ServeMsg(nc)); err != nil {
			return err
		}
	}
	{
		ra := kitnats.NewSubscriber(
			makeReadinessEndpoint(svc),
			decodeNATSReadinessRequestRequest,
			kitnats.EncodeJSONResponse,
			kitnats.SubscriberErrorLogger(logger),
		)
		if _, err := nc.QueueSubscribe("experimental-echo-service-health-readiness", "health", ra.ServeMsg(nc)); err != nil {
			return err
		}
	}
	{
		ra := kitnats.NewSubscriber(
			makeVersionEndpoint(svc),
			decodeNATSVersionRequestRequest,
			kitnats.EncodeJSONResponse,
			kitnats.SubscriberErrorLogger(logger),
		)
		if _, err := nc.QueueSubscribe("experimental-echo-service-health-version", "health", ra.ServeMsg(nc)); err != nil {
			return err
		}
	}

	return nil
}

func decodeNATSLivenessRequestRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request LivenessRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeNATSReadinessRequestRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request ReadinessRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeNATSVersionRequestRequest(_ context.Context, msg *nats.Msg) (interface{}, error) {
	var request VersionRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		return nil, err
	}
	return request, nil
}
