package security

import (
	"context"
	"time"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
)

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
	return ctx.Value(cfg.GetSecuritySessionKey_UserCode()).(string)
}

func Current_UserKey(ctx context.Context) string {
	return ctx.Value(cfg.GetSecuritySessionKey_UserKey()).(string)
}

func Current_SessionID(ctx context.Context) string {
	return ctx.Value(cfg.GetSecuritySessionKey_Session()).(string)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return ctx.Value(cfg.GetSecuritySessionKey_Token()).(sessionStore.Session_Store)
}

func Current_SessionExpiry(ctx context.Context) time.Time {
	return ctx.Value(cfg.GetSecuritySessionKey_ExpiryPeriod()).(time.Time)
}

func setSessionContextValues(ctx context.Context, user securityModel.UserMessage, sessionID string, session sessionStore.Session_Store) context.Context {
	ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Session(), sessionID)
	ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_Token(), session)
	ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserKey(), user.Key)
	ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_UserCode(), user.Code)
	ctx = context.WithValue(ctx, cfg.GetSecuritySessionKey_ExpiryPeriod(), session.Expiry)
	return ctx
}
