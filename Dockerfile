FROM golang:bullseye AS builder
WORKDIR /workspace
COPY . .
ARG version="0.0.0"
RUN go build -ldflags "-X main.Version=$version" ./cmd/ztgrep
FROM gcr.io/distroless/base-debian11
COPY --from=builder /workspace/ztgrep /bin/ztgrep
ENTRYPOINT ["/bin/ztgrep"]
