package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Id         uuid.UUID `json:"id"`
	TimeStamp  time.Time `json:"time_stamp"`
	Name       string    `json:"name"`
	Attributes any       `json:"attributes"`
}
