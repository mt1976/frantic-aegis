package security

import (
	"context"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
)

func Get(ctx context.Context) *securityModel.Session {
	si := securityModel.Session{}
	si.SessionID = ctx.Value(sessionKey).(string)
	si.UserID = ctx.Value(sessionUserIDKey).(int)
	si.Token = ctx.Value(sessionTokenKey).(sessionStore.Session_Store)
	si.UserCode = ctx.Value(sessionUserCodeKey).(string)
	return &si
}

func Current_UserCode(ctx context.Context) string {
	return ctx.Value(sessionUserCodeKey).(string)
}

func Current_UserID(ctx context.Context) int {
	return ctx.Value(sessionUserIDKey).(int)
}

func Current_SessionID(ctx context.Context) string {
	return ctx.Value(sessionKey).(string)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return ctx.Value(sessionTokenKey).(sessionStore.Session_Store)
}
