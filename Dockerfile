FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/homestead ./cmd/homestead/
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl
COPY --from=builder /bin/homestead /usr/local/bin/homestead
ENV PORT="8990" DATA_DIR="/data"
EXPOSE 8990
HEALTHCHECK --interval=30s --timeout=5s CMD curl -sf http://localhost:8990/health || exit 1
ENTRYPOINT ["homestead"]
