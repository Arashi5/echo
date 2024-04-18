package health

import (
	"context"
	"strings"

	"github.com/go-kit/kit/log"
	json "github.com/json-iterator/go"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type healthHttpClient struct {
	client *fasthttp.HostClient
	logger log.Logger
	tracer stdopentracing.Tracer

	maxRedirectsCount int
}

const (
	httpPrefix        = "http://"
	httpsPrefix       = "https://"
	maxRedirectsCount = 5

	healthMethodPath = "/health"
)

// NewHTTPClient returns an Service backed by an HTTP server living at the
// remote instance. We expect instance to come from a service discovery system,
// so likely of the form "host:port". We bake-in certain middlewares,
// implementing the client library pattern.
func NewHTTPClient(instance string, tracer stdopentracing.Tracer, logger log.Logger) (Service, error) {
	// HTTPS settings
	useHttps := true
	if strings.HasPrefix(instance, httpPrefix) {
		useHttps = false
	}

	client := fasthttp.HostClient{Addr: cutProtocolPrefix(instance)}
	client.IsTLS = useHttps

	return healthHttpClient{
		client:            &client,
		maxRedirectsCount: maxRedirectsCount,
		logger:            logger,
		tracer:            tracer,
	}, nil
}

func cutProtocolPrefix(host string) string {
	host, _ = strings.CutPrefix(host, httpPrefix)
	host, _ = strings.CutPrefix(host, httpsPrefix)

	return host
}

func (c healthHttpClient) Liveness(ctx context.Context, reqDto *LivenessRequest) (*LivenessResponse, error) {
	url := fasthttp.AcquireURI()
	errParse := url.Parse([]byte(c.client.Addr), []byte(healthMethodPath))
	if errParse != nil {
		return nil, errParse
	}

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetURI(url)
	httpRequest.Header.SetMethod(fasthttp.MethodPost)
	fasthttp.ReleaseURI(url)

	encodeErr := encodeHTTPLivenessLivenessRequest(ctx, httpRequest, reqDto)
	if encodeErr != nil {
		return nil, encodeErr
	}

	response := fasthttp.AcquireResponse()
	reqError := c.client.DoRedirects(httpRequest, response, c.maxRedirectsCount)
	fasthttp.ReleaseRequest(httpRequest)
	if reqError != nil {
		return nil, reqError
	}

	return decodeHTTPLivenessLivenessRequest(ctx, response)
}

func encodeHTTPLivenessLivenessRequest(_ context.Context, r *fasthttp.Request, request *LivenessRequest) error {
	reqJson, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "encode request body")
	}

	r.SetBody(reqJson)

	return nil
}

func decodeHTTPLivenessLivenessRequest(_ context.Context, r *fasthttp.Response) (*LivenessResponse, error) {
	if r.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(string(r.Header.StatusMessage()))
	}

	var response LivenessResponse
	if err := json.Unmarshal(r.Body(), &response); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}

	return &response, nil
}
func (c healthHttpClient) Readiness(ctx context.Context, reqDto *ReadinessRequest) (*ReadinessResponse, error) {
	url := fasthttp.AcquireURI()
	errParse := url.Parse([]byte(c.client.Addr), []byte(healthMethodPath))
	if errParse != nil {
		return nil, errParse
	}

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetURI(url)
	httpRequest.Header.SetMethod(fasthttp.MethodPost)
	fasthttp.ReleaseURI(url)

	encodeErr := encodeHTTPReadinessReadinessRequest(ctx, httpRequest, reqDto)
	if encodeErr != nil {
		return nil, encodeErr
	}

	response := fasthttp.AcquireResponse()
	reqError := c.client.DoRedirects(httpRequest, response, c.maxRedirectsCount)
	fasthttp.ReleaseRequest(httpRequest)
	if reqError != nil {
		return nil, reqError
	}

	return decodeHTTPReadinessReadinessRequest(ctx, response)
}

func encodeHTTPReadinessReadinessRequest(_ context.Context, r *fasthttp.Request, request *ReadinessRequest) error {
	reqJson, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "encode request body")
	}

	r.SetBody(reqJson)

	return nil
}

func decodeHTTPReadinessReadinessRequest(_ context.Context, r *fasthttp.Response) (*ReadinessResponse, error) {
	if r.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(string(r.Header.StatusMessage()))
	}

	var response ReadinessResponse
	if err := json.Unmarshal(r.Body(), &response); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}

	return &response, nil
}
func (c healthHttpClient) Version(ctx context.Context, reqDto *VersionRequest) (*VersionResponse, error) {
	url := fasthttp.AcquireURI()
	errParse := url.Parse([]byte(c.client.Addr), []byte(healthMethodPath))
	if errParse != nil {
		return nil, errParse
	}

	httpRequest := fasthttp.AcquireRequest()
	httpRequest.SetURI(url)
	httpRequest.Header.SetMethod(fasthttp.MethodPost)
	fasthttp.ReleaseURI(url)

	encodeErr := encodeHTTPVersionVersionRequest(ctx, httpRequest, reqDto)
	if encodeErr != nil {
		return nil, encodeErr
	}

	response := fasthttp.AcquireResponse()
	reqError := c.client.DoRedirects(httpRequest, response, c.maxRedirectsCount)
	fasthttp.ReleaseRequest(httpRequest)
	if reqError != nil {
		return nil, reqError
	}

	return decodeHTTPVersionVersionRequest(ctx, response)
}

func encodeHTTPVersionVersionRequest(_ context.Context, r *fasthttp.Request, request *VersionRequest) error {
	reqJson, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "encode request body")
	}

	r.SetBody(reqJson)

	return nil
}

func decodeHTTPVersionVersionRequest(_ context.Context, r *fasthttp.Response) (*VersionResponse, error) {
	if r.StatusCode() != fasthttp.StatusOK {
		return nil, errors.New(string(r.Header.StatusMessage()))
	}

	var response VersionResponse
	if err := json.Unmarshal(r.Body(), &response); err != nil {
		return nil, errors.Wrap(err, "decode request body")
	}

	return &response, nil
}
