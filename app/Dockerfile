# Multi-Architecture Dockerfile
# Supports AMD64 and ARM64 (Graviton)

FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

ARG TARGETOS TARGETARCH

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o server .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/server /server

EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/server"]
