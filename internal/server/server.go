// Code generated by gitlab.lenvendo.ru/product/grade-factor/services/service-generator
package server

import (
	"time"

	"github.com/arashi5/echo/configs"
	"github.com/fasthttp/router"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

// NewServer инициализирует сервер.
func NewServer(ops ...Option) (svc *Server, err error) {
	svc = new(Server)

	for _, o := range ops {
		o(svc)
	}

	return svc, nil
}

func SetLogger(logger log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func SetConfig(cfg *configs.Config) Option {
	return func(s *Server) {
		s.cfg = cfg
	}
}

func SetHandler(handlers map[string][]GroupHandler) Option {
	sts := make([]string, 0)
	for s := range handlers {
		sts = append(sts, s)
	}

	return func(s *Server) {
		r := router.New()
		for i := len(sts) - 1; i >= 0; i-- {
			name := sts[i]
			for _, h := range handlers[name] {
				r.Handle(h.Method, h.Path, h.Handler)
			}
		}

		s.handler = r.Handler
	}
}

func SetGRPC(joins ...func(grpc *grpc.Server)) Option {
	return func(s *Server) {
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpctransport.Interceptor),
			grpc.ConnectionTimeout(time.Second*time.Duration(s.cfg.Server.GRPC.TimeoutSec)),
		)
		for _, j := range joins {
			j(grpcServer)
		}
		s.grpc = grpcServer
	}
}
