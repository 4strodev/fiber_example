package auth

import (
	"time"

	"github.com/4strodev/go_monitoring_example/pkg/shared/events"
	"github.com/google/uuid"
)

const AUTH_USER_LOGGED_IN = "auth.user_logged_in"

type UserLoggedInAttributes struct {
	Id string
}

func NewUserLoggedInEvent(id string) events.Event {
	return events.Event{
		Id:        uuid.Must(uuid.NewV7()),
		Name:      AUTH_USER_LOGGED_IN,
		TimeStamp: time.Now(),
		Attributes: UserLoggedInAttributes{
			id,
		},
	}
}
