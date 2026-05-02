package errors

import "fmt"

type ErrorCode int

const (
    ErrUnknown ErrorCode = iota
    ErrInvalidOperation
    ErrUnauthorized
    ErrEntityNotFound
    ErrRegionNotFound
    ErrInsufficientResources
    ErrInvalidParams
    ErrCooldownActive
    ErrOperationFailed
    ErrInternalError
    ErrServiceUnavailable
)

type GameError struct {
    Code    ErrorCode
    Message string
    Details map[string]interface{}
}

func (e *GameError) Error() string {
    return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func NewGameError(code ErrorCode, message string) *GameError {
    return &GameError{
        Code:    code,
        Message: message,
        Details: make(map[string]interface{}),
    }
}

func (e *GameError) WithDetail(key string, value interface{}) *GameError {
    e.Details[key] = value
    return e
}

var (
    ErrInvalidOperationType = NewGameError(ErrInvalidOperation, "invalid operation type")
    ErrUnauthorizedAccess   = NewGameError(ErrUnauthorized, "unauthorized access")
    ErrEntityNotFound_      = NewGameError(ErrEntityNotFound, "entity not found")
    ErrRegionNotFound_      = NewGameError(ErrRegionNotFound, "region not found")
    ErrInsufficientFunds    = NewGameError(ErrInsufficientResources, "insufficient spirit stones")
    ErrInvalidParams_       = NewGameError(ErrInvalidParams, "invalid parameters")
    ErrCooldownActive_      = NewGameError(ErrCooldownActive, "cooldown is still active")
    ErrBreakthroughFailed   = NewGameError(ErrOperationFailed, "breakthrough failed")
    ErrInternalError_       = NewGameError(ErrInternalError, "internal server error")
)
