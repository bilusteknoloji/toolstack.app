FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG BUILD_SHA
ARG BUILD_DATE

RUN find ./site -type f -name "*.html" -exec sed -i "s/\[\[BUILD_SHA\]\]/${BUILD_SHA}/g; s/\[\[BUILD_DATE\]\]/${BUILD_DATE}/g" {} \;

ARG GOOS
ARG GOARCH
RUN CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w" -o server .

FROM alpine:latest AS certs
RUN apk add --update --no-cache ca-certificates

FROM busybox:latest
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/server /server
COPY --from=builder /build/site /site

EXPOSE 8000
ENTRYPOINT ["/server"]

LABEL org.opencontainers.image.source="https://github.com/bilusteknoloji/toolstack.app"
LABEL org.opencontainers.image.authors="Uğur vigo Özyılmazel <vigo@bilus.org>"
LABEL org.opencontainers.image.licenses="MIT"
