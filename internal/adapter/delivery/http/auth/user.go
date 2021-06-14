package auth

import (
	"context"

	"github.com/pkg/errors"
)

type User struct {
	ID   string
	Role Role
}

type Role string

const (
	TeacherRole = "teacher"
	StudentRole = "student"
)

type ctxKey int

const userCtxKey ctxKey = iota

var ErrNoUserInContext = errors.New("No user in context")

func UserFromCtx(ctx context.Context) (User, error) {
	switch u := ctx.Value(userCtxKey).(type) {
	case User:
		return u, nil
	default:
		return User{}, ErrNoUserInContext
	}
}
