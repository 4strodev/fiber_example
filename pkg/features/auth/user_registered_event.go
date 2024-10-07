package auth

import (
	"time"

	"github.com/4strodev/go_monitoring_example/pkg/shared/events"
	"github.com/google/uuid"
)

const AUTH_USER_REGISTERED = "auth.user_registered"

type UserRegisteredAttributes struct {
	Id string
}

func NewUserRegisteredEvent(id string) events.Event {
	return events.Event{
		Id:        uuid.Must(uuid.NewV7()),
		TimeStamp: time.Now(),
		Name:      AUTH_USER_REGISTERED,
		Attributes: UserRegisteredAttributes{
			id,
		},
	}
}
