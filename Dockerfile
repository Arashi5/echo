FROM golang:1.22-alpine AS build

ENV APP=./cmd/app
ENV BIN=/bin/experimental-echo-service
ENV PATH_ROJECT=${GOPATH}/src/github.com/arashi5/echo
ENV GO111MODULE=on
ENV GOSUMDB=off
ENV GOFLAGS=-mod=vendor
ARG GITLAB_DEPLOYMENT_PRIVATE_KEY
ENV GITLAB_DEPLOYMENT_PRIVATE_KEY ${GITLAB_DEPLOYMENT_PRIVATE_KEY:-unknown}
ARG VERSION
ENV VERSION ${VERSION}
ARG BUILD_TIME
ENV BUILD_TIME ${BUILD_TIME:-unknown}
ARG COMMIT
ENV COMMIT ${COMMIT:-unknown}

#ADD https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz /tmp/
#RUN cd  /tmp/ && tar -xzf migrate.linux-amd64.tar.gz &&  rm migrate.linux-amd64.tar.gz  && chmod +x /tmp/migrate

WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}

WORKDIR ${PATH_ROJECT}
COPY . ${PATH_ROJECT}
CMD go test -v -race -timeout=5s ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
        -X github.com/arashi5/echo/pkg/health.Version=${VERSION} \
        -X github.com/arashi5/echo/pkg/health.Commit=${COMMIT} \
        -X github.com/arashi5/echo/pkg/health.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo -o ${BIN} ${APP}

FROM alpine:3.11 as production
COPY --from=build ${BIN} ${BIN}

#COPY --from=build /tmp/migrate /bin/migrate

#WORKDIR /migrations

#COPY ./migrations /migrations

ENTRYPOINT ["/bin/experimental-echo-service"]
