package types

import "time"

type ActionType string

const (
    ActionCultivate     ActionType = "cultivate"
    ActionBreakthrough  ActionType = "breakthrough"
    ActionCombat        ActionType = "combat"
    ActionExplore       ActionType = "explore"
    ActionGather        ActionType = "gather"
    ActionCraft         ActionType = "craft"
    ActionCreateMethod  ActionType = "create_method"
    ActionTrade         ActionType = "trade"
    ActionFormSect      ActionType = "form_sect"
    ActionJoinSect      ActionType = "join_sect"
    ActionSendMessage   ActionType = "send_message"
    ActionCastSpell     ActionType = "cast_spell"
    ActionMeditate      ActionType = "meditate"
    ActionSleep         ActionType = "sleep"
    ActionMove          ActionType = "move"
)

type Operation struct {
    ID         string                 `json:"id"`
    ActorID    EntityID               `json:"actor_id"`
    ActionType ActionType             `json:"action_type"`
    Params     map[string]interface{} `json:"params"`
    Timestamp  int64                  `json:"timestamp"`
    Signature  string                 `json:"signature"`
}

func NewOperation(actorID EntityID, actionType ActionType, params map[string]interface{}) *Operation {
    return &Operation{
        ID:         generateOperationID(),
        ActorID:    actorID,
        ActionType: actionType,
        Params:     params,
        Timestamp:  time.Now().UnixNano(),
    }
}

type OperationResult struct {
    Success   bool                   `json:"success"`
    Message   string                 `json:"message"`
    Effects   map[string]interface{} `json:"effects"`
    Timestamp int64                  `json:"timestamp"`
}

type ValidationResult struct {
    Valid    bool     `json:"valid"`
    Errors   []string `json:"errors"`
    Warnings []string `json:"warnings"`
}
