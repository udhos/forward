# STEP 1 build executable binary

#FROM golang:alpine as builder
FROM golang:1.20.4-alpine3.17 as builder

RUN apk update
RUN apk add git

# Create appuser on builder image
RUN adduser -D -g '' appuser

COPY ./cmd/ /tmp/app/cmd
COPY go.* /tmp/app/
WORKDIR /tmp/app
RUN go mod tidy
RUN go env -w CGO_ENABLED=0
RUN go build -o /tmp/forward github.com/udhos/forward/cmd/forward

#
# STEP 2 build a small image from alpine
#
# curl: https://security.alpinelinux.org/srcpkg/curl
#
FROM alpine:3.17.3
COPY --from=builder /tmp/forward /bin/forward
RUN apk add curl=8.1.0-r1 libcrypto3=3.0.8-r4 libssl3=3.0.8-r4
RUN adduser -D -g '' user
USER user
ENTRYPOINT ["/bin/forward"]