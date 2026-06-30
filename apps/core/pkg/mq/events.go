package mq

type UserRegisteredEvent struct {
	UserID   string
	Username string
	Email    string
}
