FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY shared/go.mod shared/go.sum ./shared/
COPY ai-scheduler/go.mod ai-scheduler/go.sum ./ai-scheduler/
RUN cd ai-scheduler && go mod download

COPY shared/ ./shared/
COPY ai-scheduler/ ./ai-scheduler/

RUN cd ai-scheduler && CGO_ENABLED=0 GOOS=linux go build -o ai-scheduler ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/ai-scheduler/ai-scheduler .

EXPOSE 50052

CMD ["./ai-scheduler"]
