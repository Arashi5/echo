package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const ServiceName = "echo"

var options = []option{
	{"config", "string", "", "config file"},

	{"server.http.port", "int", 8010, "server http port"},
	{"server.http.timeout_sec", "int", 86400, "server http connection timeout"},
	{"server.grpc.port", "int", 9010, "server grpc port"},
	{"server.grpc.timeout_sec", "int", 86400, "server grpc connection timeout"},

	{"logger.level", "string", "emerg", "Level of logging. A string that correspond to the following levels: emerg, alert, crit, err, warning, notice, info, debug"},
	{"logger.time_format", "string", "2006-01-02T15:04:05.999999999", "Date format in logs"},

	{"sentry.enabled", "bool", false, "Enables or disables sentry"},
	{"sentry.dsn", "string", "https://829c0fb5737e4fc19997a076d355ece5@sentry.dev.kubedev.ru/4", "Data source name. Sentry addr"},
	{"sentry.environment", "string", "local", "The environment to be sent with events."},

	{"tracer.enabled", "bool", true, "Enables or disables tracing"},
	{"tracer.host", "string", "localhost", "The tracer host"},
	{"tracer.port", "int", 4317, "The tracer port"},

	{"metrics.enabled", "bool", false, "Enables or disables metrics"},
	{"metrics.port", "int", 9155, "server http port"},

	{"limiter.enabled", "bool", false, "Enables or disables limiter"},
	{"limiter.limit", "float64", 10000.0, "Limit tokens per second"},

	{"db.postgres.host", "string", "localhoxt", "postgres master host"},
	{"db.postgres.port", "int", 5432, "postgres master port"},
	{"db.postgres.user", "string", "postgres", "postgres master user"},
	{"db.postgres.password", "string", "postgres", "postgres master password"},
	{"db.postgres.database_name", "string", "echo", "postgres master database name"},
	{"db.postgres.secure", "string", "disable", "postgres master SSL support"},
	{"db.postgres.schema", "string", "echo", "schema"},
	{"db.postgres.limit", "int", 15, "max number of connections pool postgres"},
}

type Config struct {
	Server struct {
		GRPC struct {
			Port       int
			TimeoutSec int `mapstructure:"timeout_sec"`
		}
		HTTP struct {
			Port       int
			TimeoutSec int `mapstructure:"timeout_sec"`
		}
	}
	Logger struct {
		Level      string
		TimeFormat string
	}
	Sentry struct {
		Enabled     bool
		Dsn         string
		Environment string
	}
	Tracer struct {
		Enabled bool
		Host    string
		Port    int
	}
	Metrics struct {
		Enabled bool
		Port    int
	}
	Limiter struct {
		Enabled bool
		Limit   float64
	}
	DB struct {
		Postgres Database
	}
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

type Database struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string `mapstructure:"database_name"`
	Secure       string
	Schema       string
	Limit        uint32
}

// NewConfig returns and prints struct with config parameters
func NewConfig() *Config {
	return &Config{}
}

// read gets parameters from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		case "float64":
			pflag.Float64(o.name, o.value.(float64), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigName(fileName)
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrap(err, "failed to read from file")
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
