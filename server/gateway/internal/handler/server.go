package handler

import (
    "fmt"
    "net/http"
    "time"

    "github.com/cultivation-world/gateway/internal/service"
    "github.com/cultivation-world/shared/config"
    "github.com/cultivation-world/shared/types"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type Server struct {
    cfg        *config.Config
    wsHub      *WebSocketHub
    authSvc    *service.AuthService
    gameClient *service.GameServiceClient
    engine     *gin.Engine
}

func NewServer(cfg *config.Config, wsHub *WebSocketHub, authSvc *service.AuthService, gameClient *service.GameServiceClient) *Server {
    gin.SetMode(gin.ReleaseMode)
    engine := gin.New()
    engine.Use(gin.Recovery())
    engine.Use(corsMiddleware())

    s := &Server{
        cfg:        cfg,
        wsHub:      wsHub,
        authSvc:    authSvc,
        gameClient: gameClient,
        engine:     engine,
    }

    s.setupRoutes()
    return s
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func (s *Server) setupRoutes() {
    s.engine.POST("/auth/register", s.handleRegister)
    s.engine.POST("/auth/login", s.handleLogin)
    s.engine.GET("/ws", s.handleWebSocket)
    s.engine.GET("/health", s.handleHealth)
}

func (s *Server) Start() error {
    addr := fmt.Sprintf(":%d", s.cfg.Server.Port)
    return s.engine.Run(addr)
}

func (s *Server) handleRegister(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    entity, token, err := s.authSvc.Register(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "token":   token,
        "entity":  entity,
    })
}

func (s *Server) handleLogin(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    entity, token, err := s.authSvc.Login(req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "token":   token,
        "entity":  entity,
    })
}

func (s *Server) handleWebSocket(c *gin.Context) {
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
        return
    }

    entityID, err := s.authSvc.ValidateToken(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        return
    }

    client := NewWebSocketClient(entityID, conn, s.wsHub, s.gameClient)
    s.wsHub.Register(client)

    go client.WritePump()
    go client.ReadPump()
}

func (s *Server) handleHealth(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
    })
}

func sendMessage(conn *websocket.Conn, msgType types.MessageType, payload interface{}) error {
    msg := types.Message{
        Type:      msgType,
        Payload:   make(map[string]interface{}),
        Timestamp: time.Now().UnixNano(),
    }

    return conn.WriteJSON(msg)
}
