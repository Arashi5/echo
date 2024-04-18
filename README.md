# Сервис сборки

## Кофигурация

Сервис может быть сконфигурирован тремя способами:

- через флаги
- через переменные окружения
- через файл конфига в формате toml

Возможные параметры конфигурации:

| Параметр                                  |   Переменная окружения                    | Значение по умолчанию               | Описание                              |
| -------------                             | :-------------                            |:------------------------------------|:-------------                         |
| server.grpc.port                          | ECHO_SERVICE_SERVER_GRPC_PORT                     | 9090                                | grpc port server                      |
| server.grpc.timeout_sec                   | ECHO_SERVICE_SERVER_GRPC_TIMEOUT_SEC              | 86400                               | server grpc connection timeout        |
| server.http.port                          | ECHO_SERVICE_SERVER_HTTP_PORT                     | 8080                                | http port server                      |
| server.http.timeout_sec                   | ECHO_SERVICE_SERVER_HTTP_TIMEOUT_SEC              | 86400                               | server http connection timeout    |
| limiter.enabled                           | ECHO_SERVICE_LIMITER_ENABLED                      | false                               | Enables or disables limiter |
| limiter.limit                             | ECHO_SERVICE_LIMITER_LIMIT                        | 10000.0                             | Limit tokens per second     |
| sentry.enabled                            | ECHO_SERVICE_SENTRY_ENABLED                       | false                               | Enables or disables sentry                            |
| sentry.dsn                                | ECHO_SERVICE_SENTRY_DSN                           | localhost                           |Sentry addres |
| sentry.environment                        | ECHO_SERVICE_SENTRY_ENVIRONMENT                   | dev                                 | The environment to be sent with events |
| tracer.enabled                            | ECHO_SERVICE_TRACER_ENABLED                       | false                               | флаг, если указан, то в opentracing будут отправляться трассировки путей запросов (если передан через флаги, то любое значение будет соотвествоать true)     |
| tracer.host                               | ECHO_SERVICE_TRACER_HOST                          | 127.0.0.1                           | хост трасировщика                                     |
| tracer.port                               | ECHO_SERVICE_TRACER_PORT                          | 5775                                | порт трасировщика                                     |
| tracer.name                               | ECHO_SERVICE_TRACER_NAME                          | STMS                                | название трасировщика                                     |
| metrics.enabled                           | ECHO_SERVICE_METRICS_ENABLED                      | false                               | Enables or disables metric                            |
| metrics.port                              | ECHO_SERVICE_METRICS_PORT                         | 9153                                | metrics server http port                              |
| logger.level                              | ECHO_SERVICE_LOGGER_LEVEL                         | emerg                               | log level ([syslog](https://en.wikipedia.org/wiki/Syslog#Severity_level))              |
| logger.time.format                        | ECHO_SERVICE_LOGGER_TIME_FORMST                   | 2006-01-02T15:04:05.999999999Z07:00 |[time format for logger](https://golang.org/src/time/format.go)                |

Флаги имеют наивысший приортитет, файл конфига - наинизший

### Конфигурация через флаги

```bash
$ ./echo-service --server.grpc.port=9090 --server.http.port=8080
```

### Конфигурация через переменные окружения

Имена переменных окружения должны начинаться с префикса ECHO_SERVICE (echo service), точки заменяются знаком
подчеркивания

```
STMS_SERVER_GRPC_PORT=9090
```

### Конфигурация через файл конфига

При запуске сервиса можно указать путь до файла с конфигурацией

```bash
$ ./echo-servic --config=./path/to/config.toml
```

Пример файла конфигурации см. в папке `/configs`: config.toml.dist

## База данных

В качестве базы данных используется Postgres.

Для миграций используется пакет [migrate](https://github.com/golang-migrate/migrate)

`/migrations` - миграции

Создать миграцию

```
$ migrate create -ext sql -dir ./migrations create_echo_table
```

Накатить/откатить миграции up/down

```
$ migrate -source file://migrations -database "postgres://postgres:postgres@localhost:5432/echo?sslmode=disable" up/down
```