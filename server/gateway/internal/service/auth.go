package service

import (
    "context"
    "errors"
    "time"

    "github.com/cultivation-world/shared/types"
    "github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
    jwtSecret  string
    gameClient *GameServiceClient
}

type Claims struct {
    EntityID types.EntityID `json:"entity_id"`
    Username string         `json:"username"`
    jwt.RegisteredClaims
}

func NewAuthService(jwtSecret string, gameClient *GameServiceClient) *AuthService {
    return &AuthService{
        jwtSecret:  jwtSecret,
        gameClient: gameClient,
    }
}

func (s *AuthService) Register(username, password string) (*types.Entity, string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    entity, err := s.gameClient.CreateEntity(ctx, username, password, types.EntityTypePlayer)
    if err != nil {
        return nil, "", err
    }

    token, err := s.generateToken(entity.ID, username)
    if err != nil {
        return nil, "", err
    }

    return entity, token, nil
}

func (s *AuthService) Login(username, password string) (*types.Entity, string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    entity, err := s.gameClient.AuthenticateEntity(ctx, username, password)
    if err != nil {
        return nil, "", errors.New("invalid credentials")
    }

    token, err := s.generateToken(entity.ID, username)
    if err != nil {
        return nil, "", err
    }

    return entity, token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (types.EntityID, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims.EntityID, nil
    }

    return "", errors.New("invalid token")
}

func (s *AuthService) generateToken(entityID types.EntityID, username string) (string, error) {
    claims := &Claims{
        EntityID: entityID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "cultivation-gateway",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}
