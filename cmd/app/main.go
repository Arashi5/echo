package main

import (
	"context"
	"fmt"
	"github.com/arashi5/echo/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"

	"github.com/go-kit/kit/log/level"

	"github.com/arashi5/echo/configs"
	"github.com/arashi5/echo/internal/server"
	"github.com/arashi5/echo/tools/logging"
	"github.com/arashi5/echo/tools/metrics"
	"github.com/arashi5/echo/tools/sentry"

	"github.com/arashi5/echo/pkg/echo"
	"github.com/arashi5/echo/pkg/health"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}
	// Print config
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %s", err)
		os.Exit(1)
	}
	ctx = logging.WithContext(ctx, logger)

	if cfg.Sentry.Enabled {
		if err := sentry.NewSentry(cfg); err != nil {
			level.Error(logger).Log("err", err, "msg", "failed to init sentry")
		}
	}

	if cfg.Metrics.Enabled {
		ctx = metrics.WithContext(ctx)
	}
	conn, err := pgxConnectionPool(ctx, cfg.DB.Postgres)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	echoService := initEchoService(ctx, cfg, repository.NewRepositoryUnitOfWork(conn))
	healthService := initHealthService(ctx, cfg)

	s, err := server.NewServer(
		server.SetConfig(cfg),
		server.SetLogger(logger),
		server.SetHandler(
			map[string][]server.GroupHandler{
				"echo":   echo.MakeHTTPHandler(ctx, echoService),
				"health": health.MakeHTTPHandler(ctx, healthService),
			}),
		server.SetGRPC(
			echo.JoinGRPC(ctx, echoService),
			health.JoinGRPC(ctx, healthService),
		),
	)
	if err != nil {
		level.Error(logger).Log("init", "server", "err", err)
		os.Exit(1)
	}
	defer s.Close()

	if err := s.AddHTTP(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	if err = s.AddGRPC(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	if err = s.AddMetrics(); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	s.AddSignalHandler()
	s.Run()
}

func initEchoService(ctx context.Context, cfg *configs.Config, repo *repository.Instance) echo.Service {
	echoService := echo.NewEchoService(repo)
	if cfg.Metrics.Enabled {
		echoService = echo.NewMetricsService(ctx, echoService)
	}
	echoService = echo.NewLoggingService(ctx, echoService)
	if cfg.Sentry.Enabled {
		echoService = echo.NewSentryService(echoService)
	}
	return echoService
}

func initHealthService(ctx context.Context, cfg *configs.Config) health.Service {
	healthService := health.NewHealthService()
	if cfg.Metrics.Enabled {
		healthService = health.NewMetricsService(ctx, healthService)
	}
	healthService = health.NewLoggingService(ctx, healthService)
	if cfg.Sentry.Enabled {
		healthService = health.NewSentryService(healthService)
	}
	return healthService
}

func pgxConnectionPool(ctx context.Context, cfgDB configs.Database) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, fmt.Sprintf(`
 		user=%s
		password=%s
        host=%s
        port=%d
        dbname=%s

        sslmode=%s
        search_path=%s
		pool_max_conns=%d`,
		cfgDB.User,
		cfgDB.Password,
		cfgDB.Host,
		cfgDB.Port,
		cfgDB.DatabaseName,
		cfgDB.Secure,
		cfgDB.Schema,
		cfgDB.Limit,
	))
}
