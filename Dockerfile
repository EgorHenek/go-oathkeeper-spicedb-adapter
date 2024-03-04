FROM golang:1.22-alpine as builder

WORKDIR /go/src/app
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg/mod CGO_ENABLED=0 go build -v ./cmd/...

FROM cgr.dev/chainguard/static:latest
COPY --from=builder /go/src/app/adapter /usr/local/bin/adapter
ENV PATH="$PATH:/usr/local/bin"
ENTRYPOINT ["adapter"]