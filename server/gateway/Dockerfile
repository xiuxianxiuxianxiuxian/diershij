FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY gateway/go.mod gateway/go.sum ./gateway/
RUN cd gateway && go mod download

COPY shared/ ./shared/
COPY gateway/ ./gateway/

RUN cd gateway && CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/gateway/gateway .

EXPOSE 8080

CMD ["./gateway"]
