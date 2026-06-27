package model

type NewUserCreatedEvent struct {
	UserId   int
	Username string
	Email    string
}

var newUserCreatedHandler func(NewUserCreatedEvent)

func RegisterNewUserCreatedHandler(handler func(NewUserCreatedEvent)) {
	newUserCreatedHandler = handler
}

func triggerNewUserCreated(event NewUserCreatedEvent) {
	if newUserCreatedHandler == nil || event.UserId <= 0 {
		return
	}
	newUserCreatedHandler(event)
}
