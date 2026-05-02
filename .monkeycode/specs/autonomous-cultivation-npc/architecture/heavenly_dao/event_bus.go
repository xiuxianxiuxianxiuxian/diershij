package heavenlydao

import "time"

type EventBus interface {
	Publish(event RuleEvent) error
}

type MemoryEventBus struct {
	Events []RuleEvent
}

func (b *MemoryEventBus) Publish(event RuleEvent) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	b.Events = append(b.Events, event)
	return nil
}
