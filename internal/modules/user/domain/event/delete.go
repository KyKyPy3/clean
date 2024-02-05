package event

const UserDeleted = "UserDeleted"

type UserDeletedEvent struct {
	ID string
}

func (e UserDeletedEvent) Kind() string {
	return UserDeleted
}
