package security

import (
	"context"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
)

var sessionKeyWrangler = aegisContentWrangler{Content: sessionKey}
var sessionUserKeyKeyWrangler = aegisContentWrangler{Content: sessionUserKeyKey}
var sessionUserCodeKeyWrangler = aegisContentWrangler{Content: sessionUserCodeKey}
var sessionTokenKeyWrangler = aegisContentWrangler{Content: sessionTokenKey}
var sessionExpiryKeyWrangler = aegisContentWrangler{Content: sessionExpiryKey}

func get(ctx context.Context) *securityModel.Session {
	si := securityModel.Session{}
	si.SessionID = Current_SessionID(ctx)
	si.SessionToken = Current_SessionToken(ctx)
	si.UserKey = Current_UserKey(ctx)
	si.UserCode = Current_UserCode(ctx)
	si.Expiry = Current_SessionExpiry(ctx)
	return &si
}

func Current_UserCode(ctx context.Context) string {
	return ctx.Value(sessionUserCodeKeyWrangler).(string)
}

func Current_UserKey(ctx context.Context) string {
	return ctx.Value(sessionUserKeyKeyWrangler).(string)
}

func Current_SessionID(ctx context.Context) string {
	return ctx.Value(sessionKey).(string)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return ctx.Value(sessionTokenKeyWrangler).(sessionStore.Session_Store)
}

func Current_SessionExpiry(ctx context.Context) time.Time {
	return ctx.Value(sessionExpiryKey).(time.Time)
}

func setSessionContextValues(ctx context.Context, user securityModel.UserMessage, sessionID string, token sessionStore.Session_Store) context.Context {
	ctx = context.WithValue(ctx, sessionKeyWrangler, sessionID)
	ctx = context.WithValue(ctx, sessionTokenKeyWrangler, token)
	ctx = context.WithValue(ctx, sessionUserKeyKeyWrangler, user.Key)
	ctx = context.WithValue(ctx, sessionUserCodeKeyWrangler, user.Code)
	ctx = context.WithValue(ctx, sessionExpiryKeyWrangler, token.Expiry)
	return ctx
}
