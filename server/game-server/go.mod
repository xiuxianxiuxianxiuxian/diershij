module github.com/cultivation-world/game-server

go 1.21

require (
    github.com/cultivation-world/shared v0.0.0
    github.com/google/uuid v1.6.0
    github.com/jackc/pgx/v5 v5.5.5
    github.com/redis/go-redis/v9 v9.5.1
    google.golang.org/grpc v1.63.2
)

replace github.com/cultivation-world/shared => ../shared
