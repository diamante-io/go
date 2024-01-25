package keystore

import "context"

type contextKey int

const userKey contextKey = iota

func userID(ctx context.Context) string {
	uid, _ := ctx.Value(userKey).(string)
	return uid
}

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey, userID)
}
