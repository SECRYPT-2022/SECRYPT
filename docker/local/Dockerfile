FROM golang:1.17-alpine AS builder

WORKDIR /secrypt

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o secrypt main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates jq

WORKDIR /secrypt

COPY --from=builder /secrypt/secrypt ./
COPY ./docker/local/secrypt.sh ./

# Expose json-rpc, libp2p and grpc ports
EXPOSE 8545 9632 1478

ENTRYPOINT ["./secrypt.sh"]
