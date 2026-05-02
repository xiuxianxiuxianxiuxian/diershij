module github.com/cultivation-world/ai-scheduler

go 1.21

require (
    github.com/cultivation-world/shared v0.0.0
    google.golang.org/grpc v1.63.2
    google.golang.org/protobuf v1.33.0
)

require (
    golang.org/x/net v0.23.0 // indirect
    golang.org/x/sys v0.18.0 // indirect
    golang.org/x/text v0.14.0 // indirect
    google.golang.org/genproto/googleapis/rpc v0.0.0-20240318140521-94a12d6c2237 // indirect
)

replace github.com/cultivation-world/shared => ../shared
