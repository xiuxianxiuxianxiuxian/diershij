module github.com/cultivation-world/gateway

go 1.21

require (
    github.com/cultivation-world/shared v0.0.0
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/gorilla/websocket v1.5.1
    github.com/redis/go-redis/v9 v9.5.1
    google.golang.org/grpc v1.63.2
)

replace github.com/cultivation-world/shared => ../shared
