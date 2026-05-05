FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY game-server/go.mod game-server/go.sum ./game-server/
RUN cd game-server && go mod download

COPY shared/ ./shared/
COPY game-server/ ./game-server/

RUN cd game-server && CGO_ENABLED=0 GOOS=linux go build -o game-server ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/game-server/game-server .

EXPOSE 50051

CMD ["./game-server"]
