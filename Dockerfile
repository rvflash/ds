# Download the dependencies (mod)
FROM golang:1.15 as modules

## Linux dependencies.
RUN apt-get update && apt-get install -y git

ADD go.mod go.sum /mod/
RUN cd /mod && go mod download


# Build the application (app)
FROM golang:1.15 as builder
COPY --from=modules /go/pkg /go/pkg

## Add a non-privileged user
RUN useradd -u 10001 ds

## Copy sources
RUN mkdir -p /app
ADD . /app
WORKDIR /app

## Build the binary
ARG APP_VERSION=v0.0.1-noop
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 go build \
    -ldflags "-X main.buildVersion=${APP_VERSION}" \
    -o bin/ds .


# Run the binary
FROM scratch

## Certificates and privileges
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER ds

## Finally exposes the command line tool
COPY --from=builder /app/bin/ds /bin/ds
WORKDIR /pkg
CMD ["/bin/ds"]