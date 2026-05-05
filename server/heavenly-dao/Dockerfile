FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY heavenly-dao/go.mod heavenly-dao/go.sum ./heavenly-dao/
RUN cd heavenly-dao && go mod download

COPY shared/ ./shared/
COPY heavenly-dao/ ./heavenly-dao/

RUN cd heavenly-dao && CGO_ENABLED=0 GOOS=linux go build -o heavenly-dao ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/heavenly-dao/heavenly-dao .

EXPOSE 50053

CMD ["./heavenly-dao"]
