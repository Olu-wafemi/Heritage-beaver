package auth

import "context"

type ContextKey string

const UserIDKey ContextKey = "user_id"

func UserIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(UserIDKey).(string)
	return id
}