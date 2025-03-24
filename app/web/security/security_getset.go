package security

import (
	"context"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

func Current_UserCode(ctx context.Context) string {
	return contextHandler.GetSession_UserCode(ctx)
}

func Current_UserKey(ctx context.Context) string {
	return contextHandler.GetSession_UserKey(ctx)
}

func Current_Locale(ctx context.Context) string {
	return contextHandler.GetSession_Locale(ctx)
}

func Current_SessionID(ctx context.Context) string {
	return contextHandler.GetSession_ID(ctx)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return contextHandler.GetSession_Token(ctx).(sessionStore.Session_Store)
	// return ctx.Value(cfg.GetSecuritySessionKey_Token()).(sessionStore.Session_Store)
}

func Current_SessionExpiry(ctx context.Context) time.Time {
	return contextHandler.GetSession_Expiry(ctx)
}

func setSessionContextValues(ctx context.Context, user messageHelpers.UserMessage, sessionID string, session sessionStore.Session_Store) context.Context {
	ctx = contextHandler.SetSession_ID(ctx, sessionID)
	ctx = contextHandler.SetSession_Token(ctx, session)
	ctx = contextHandler.SetSession_UserKey(ctx, user.Key)
	ctx = contextHandler.SetSession_UserCode(ctx, user.Code)
	ctx = contextHandler.SetSession_Expiry(ctx, session.Expiry)
	ctx = contextHandler.SetSession_Locale(ctx, user.Locale)
	ctx = contextHandler.SetSession_Theme(ctx, "")
	ctx = contextHandler.SetSession_Timezone(ctx, "")
	return ctx
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Session(), sessionID)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Token(), session)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserKey(), user.Key)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserCode(), user.Code)
	// ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_ExpiryPeriod(), session.Expiry)
	// return ctx
}
