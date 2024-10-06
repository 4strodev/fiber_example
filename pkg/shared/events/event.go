package events

import (
	"time"

	"github.com/google/uuid"
)

type EventName string

type Event[T any] struct {
	Id         uuid.UUID `json:"id"`
	TimeStamp  time.Time `json:"time_stamp"`
	Name       EventName `json:"name"`
	Attributes T         `json:"attributes"`
}
