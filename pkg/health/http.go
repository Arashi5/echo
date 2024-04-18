package health

import (
	"context"
	"maps"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/schema"
	json "github.com/json-iterator/go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/arashi5/echo/internal/server"
	"github.com/arashi5/echo/tools/logging"
)

func MakeHTTPHandler(ctx context.Context, s Service) []server.GroupHandler {
	logger := logging.FromContext(ctx)
	logger = log.With(logger, "http handler", "health")

	handlers := make([]server.GroupHandler, 0)

	handlers = append(handlers, server.GroupHandler{
		Path:   "/liveness",
		Method: "GET",
		Handler: accessControl(func(req *fasthttp.RequestCtx) {
			ctx = authTokenAuthorization(ctx, req)
			ctx = httpToContext(ctx, req)

			request, decodeErr := decodeGETLivenessRequest(ctx, req)
			if decodeErr != nil {
				encodeError(ctx, decodeErr, req)
				return
			}

			response, serviceErr := s.Liveness(ctx, request)
			if serviceErr != nil {
				encodeError(ctx, serviceErr, req)
				return
			}

			encodeErr := encodeResponse(ctx, req, response)
			if encodeErr != nil {
				encodeError(ctx, encodeErr, req)
				return
			}
		}),
	})

	handlers = append(handlers, server.GroupHandler{
		Path:   "/readiness",
		Method: "GET",
		Handler: accessControl(func(req *fasthttp.RequestCtx) {
			ctx = authTokenAuthorization(ctx, req)
			ctx = httpToContext(ctx, req)

			request, decodeErr := decodeGETReadinessRequest(ctx, req)
			if decodeErr != nil {
				encodeError(ctx, decodeErr, req)
				return
			}

			response, serviceErr := s.Readiness(ctx, request)
			if serviceErr != nil {
				encodeError(ctx, serviceErr, req)
				return
			}

			encodeErr := encodeResponse(ctx, req, response)
			if encodeErr != nil {
				encodeError(ctx, encodeErr, req)
				return
			}
		}),
	})

	handlers = append(handlers, server.GroupHandler{
		Path:   "/version",
		Method: "GET",
		Handler: accessControl(func(req *fasthttp.RequestCtx) {
			ctx = authTokenAuthorization(ctx, req)
			ctx = httpToContext(ctx, req)

			request, decodeErr := decodeGETVersionRequest(ctx, req)
			if decodeErr != nil {
				encodeError(ctx, decodeErr, req)
				return
			}

			response, serviceErr := s.Version(ctx, request)
			if serviceErr != nil {
				encodeError(ctx, serviceErr, req)
				return
			}

			encodeErr := encodeResponse(ctx, req, response)
			if encodeErr != nil {
				encodeError(ctx, encodeErr, req)
				return
			}
		}),
	})

	return handlers
}

func httpToContext(ctx context.Context, req *fasthttp.RequestCtx) context.Context {
	return context.WithValue(ctx, ContextHTTPKey{}, HTTPInfo{
		Method:   string(req.Method()),
		URL:      string(req.RequestURI()),
		From:     req.RemoteAddr().String(),
		Protocol: string(req.Request.Header.Protocol()),
	})
}

func authTokenAuthorization(ctx context.Context, req *fasthttp.RequestCtx) context.Context {
	const (
		Authorization = "Authorization"
		CompanyUuid   = "Company-UUID"
		ScriptHeader  = "X-Public-Script"
		Page          = "Page"
		XApiKey       = "X-GF-ApiKey"
	)

	ctx = context.WithValue(ctx, Authorization, string(req.Request.Header.Peek(Authorization)))
	ctx = context.WithValue(ctx, ScriptHeader, string(req.Request.Header.Peek(ScriptHeader)))
	ctx = context.WithValue(ctx, CompanyUuid, string(req.Request.Header.Peek(CompanyUuid)))
	ctx = context.WithValue(ctx, XApiKey, string(req.Request.Header.Peek(XApiKey)))

	return ctx
}

func closeHTTPTracer(ctx context.Context) {
	span := stdopentracing.SpanFromContext(ctx)
	span.Finish()
}

func decodeGETLivenessRequest(_ context.Context, r *fasthttp.RequestCtx) (*LivenessRequest, error) {
	var request LivenessRequest

	params := parseUriNamedParams(r)
	maps.Copy(params, parseUriGetParams(r))
	err := schema.NewDecoder().Decode(&request, params)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidArgument, err.Error())
	}

	return &request, nil
}

func decodeGETReadinessRequest(_ context.Context, r *fasthttp.RequestCtx) (*ReadinessRequest, error) {
	var request ReadinessRequest

	params := parseUriNamedParams(r)
	maps.Copy(params, parseUriGetParams(r))
	err := schema.NewDecoder().Decode(&request, params)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidArgument, err.Error())
	}

	return &request, nil
}

func decodeGETVersionRequest(_ context.Context, r *fasthttp.RequestCtx) (*VersionRequest, error) {
	var request VersionRequest

	params := parseUriNamedParams(r)
	maps.Copy(params, parseUriGetParams(r))
	err := schema.NewDecoder().Decode(&request, params)
	if err != nil {
		return nil, errors.Wrap(ErrInvalidArgument, err.Error())
	}

	return &request, nil
}

func encodeResponse(ctx context.Context, req *fasthttp.RequestCtx, response interface{}) error {
	if err, ok := response.(error); ok && err != nil {
		encodeError(ctx, err, req)
		return nil
	}

	req.Response.Header.SetContentType("application/json; charset=utf-8")
	return json.NewEncoder(req).Encode(response)
}

// encodeError handles error from business-layer.
func encodeError(_ context.Context, err error, req *fasthttp.RequestCtx) {
	req.Response.Header.SetStatusCode(getHTTPStatusCode(err))
	req.Response.Header.SetContentType("application/problem+json; charset=utf-8")
	req.Response.Header.Set("X-Api-Error", err.Error())

	json.NewEncoder(req).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"message": func() string {
				if strings.Index(err.Error(), ":") != -1 {
					return err.Error()[:strings.Index(err.Error(), ":")]
				}
				return err.Error()
			}(),
			"type": errors.Cause(err).Error(),
		},
	})
}

// accessControl is CORS middleware.
func accessControl(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Public-Script, X-GF-ApiKey")
		if string(ctx.Method()) == "OPTIONS" {
			return
		}

		h(ctx)
	}
}

func parseUriNamedParams(r *fasthttp.RequestCtx) map[string][]string {
	params := make(map[string][]string)
	r.VisitUserValuesAll(func(key any, value any) {
		params[key.(string)] = []string{value.(string)}
	})

	return params
}

func parseUriGetParams(r *fasthttp.RequestCtx) map[string][]string {
	params := make(map[string][]string)
	r.URI().QueryArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = []string{string(value)}
	})

	return params
}
