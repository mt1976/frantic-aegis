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

func Current_UserLocale(ctx context.Context) string {
	return contextHandler.GetSession_Locale(ctx)
}

func Current_SessionID(ctx context.Context) string {
	return contextHandler.GetSession_ID(ctx)
}

func Current_SessionToken(ctx context.Context) sessionStore.Session_Store {
	return contextHandler.GetSession_Token(ctx).(sessionStore.Session_Store)
}

func Current_SessionExpiry(ctx context.Context) time.Time {
	return contextHandler.GetSession_Expiry(ctx)
}

func Current_SessionTheme(ctx context.Context) string {
	return contextHandler.GetSession_Theme(ctx)
}

func Current_SessionTimezone(ctx context.Context) string {
	return contextHandler.GetSession_Timezone(ctx)
}

func setSessionContextValues(ctx context.Context, user messageHelpers.UserMessage, sessionID string, session sessionStore.Session_Store) context.Context {
	ctx = contextHandler.SetSession_ID(ctx, sessionID)
	ctx = contextHandler.SetSession_Token(ctx, session)
	ctx = contextHandler.SetSession_UserKey(ctx, user.Key)
	ctx = contextHandler.SetSession_UserCode(ctx, user.Code)
	ctx = contextHandler.SetSession_Expiry(ctx, session.Expiry)

	ctx = contextHandler.SetSession_Locale(ctx, session.Locale)
	if session.Locale == "" {
		ctx = contextHandler.SetSession_Locale(ctx, user.Locale)
		if user.Locale == "" {
			ctx = contextHandler.SetSession_Locale(ctx, cfg.GetApplication_Locale())
		}
	}
	ctx = contextHandler.SetSession_Theme(ctx, session.Theme)
	if session.Theme == "" {
		ctx = contextHandler.SetSession_Theme(ctx, user.Theme)
		if user.Theme == "" {
			ctx = contextHandler.SetSession_Theme(ctx, cfg.GetApplication_Theme())
		}
	}
	ctx = contextHandler.SetSession_Timezone(ctx, session.Timezone)
	if session.Timezone == "" {
		ctx = contextHandler.SetSession_Timezone(ctx, user.Timezone)
		if user.Timezone == "" {
			ctx = contextHandler.SetSession_Timezone(ctx, cfg.GetApplication_Timezone())
		}
	}
	return ctx
}
