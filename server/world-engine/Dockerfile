FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY world-engine/go.mod world-engine/go.sum ./world-engine/
RUN cd world-engine && go mod download

COPY shared/ ./shared/
COPY world-engine/ ./world-engine/

RUN cd world-engine && CGO_ENABLED=0 GOOS=linux go build -o world-engine ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/world-engine/world-engine .

EXPOSE 50054

CMD ["./world-engine"]
