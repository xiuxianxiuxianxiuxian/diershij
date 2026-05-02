package types

import "github.com/google/uuid"

func generateOperationID() string {
    return uuid.New().String()
}

func GenerateEntityID() EntityID {
    return EntityID(uuid.New().String())
}
