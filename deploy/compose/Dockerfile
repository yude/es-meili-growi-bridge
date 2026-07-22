FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download 2>/dev/null || true
COPY . .
RUN CGO_ENABLED=0 go build -o /bridge ./cmd/es-meili-growi-bridge/

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bridge /usr/local/bin/bridge
EXPOSE 9200
ENTRYPOINT ["bridge"]
