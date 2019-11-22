FROM dockerfactory.rsglab.com/rsg/golang/golang:1.13.4 AS builder

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
ARG GO111MODULE=on

COPY . /src
RUN cd /src && go build -mod=vendor -o /tmp/go-links ./cmd/go-links

FROM gcr.io/distroless/base

COPY --from=builder /tmp/go-links /usr/bin/

ENTRYPOINT ["/usr/bin/go-links" ]