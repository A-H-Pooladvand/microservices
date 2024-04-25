package handlers

import (
	"po/internal/handlers/user"
)

type WebHandlers struct {
	User *user.Handler
}

func NewWebHandlers(user *user.Handler) *WebHandlers {
	return &WebHandlers{User: user}
}
