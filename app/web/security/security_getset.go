package security

import (
	"context"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
	"github.com/mt1976/frantic-core/contextHandler"
)

func Current_UserCode(ctx context.Context) string {
	return contextHandler.GetUserCode(ctx)
}

func Current_UserKey(ctx context.Context) string {
	return contextHandler.GetUserKey(ctx)
}

func Current_SessionID(ctx context.Context) string {
	return contextHandler.GetSessionID(ctx)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return contextHandler.GetSessionToken(ctx).(sessionStore.Session_Store)
	// return ctx.Value(cfg.GetSecuritySessionKey_Token()).(sessionStore.Session_Store)
}

func Current_SessionExpiry(ctx context.Context) time.Time {
	return contextHandler.GetSessionExpiry(ctx)
}

func setSessionContextValues(ctx context.Context, user securityModel.UserMessage, sessionID string, session sessionStore.Session_Store) context.Context {
	ctx = contextHandler.SetSessionID(ctx, sessionID)
	ctx = contextHandler.SetSessionToken(ctx, session)
	ctx = contextHandler.SetUserKey(ctx, user.Key)
	ctx = contextHandler.SetUserCode(ctx, user.Code)
	ctx = contextHandler.SetSessionExpiry(ctx, session.Expiry)
	return ctx
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Session(), sessionID)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Token(), session)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserKey(), user.Key)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserCode(), user.Code)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_ExpiryPeriod(), session.Expiry)
	// return ctx
}
