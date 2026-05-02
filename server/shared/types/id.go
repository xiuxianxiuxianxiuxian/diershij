package types

import (
    "crypto/rand"
    "encoding/hex"
)

func GenerateOperationID() OperationID {
    b := make([]byte, 16)
    rand.Read(b)
    return OperationID(hex.EncodeToString(b))
}

func GenerateEntityID() EntityID {
    b := make([]byte, 16)
    rand.Read(b)
    return EntityID(hex.EncodeToString(b))
}
