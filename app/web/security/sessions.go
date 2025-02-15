package security

import (
	"context"
	"net/http"
	"strings"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logger"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userID int, userIDValidator func(int) (securityModel.UserMessage, error)) *securityModel.Session {
	clock := timing.Start(domain, "New", "")
	SI := securityModel.Session{}

	SS, err := sessionStore.New(userID)
	if err != nil {
		panic(err)
	}

	// U, err := userStore.Get(userID)
	// if err != nil {
	// 	panic(err)
	// }
	UserMessage, err := userIDValidator(userID)
	if err != nil {
		logger.ErrorLogger.Printf("Error=[%v]", err.Error())
		panic(err)
	}

	SI.Token = SS
	SI.UserID = userID
	SI.SessionID = SS.ID
	SI.Life = 0
	SI.UserCode = UserMessage.Code

	//ctx = setSessionContextValues(ctx, UserMessage, SI.SessionID, SS)

	if appModeDev {
		logger.InfoLogger.Printf("SessionID=[%v]", SI.SessionID)
		logger.InfoLogger.Printf("UserID=[%v]", SI.UserID)
		logger.InfoLogger.Printf("UserCode=[%v]", SI.UserCode)
		logger.InfoLogger.Printf("Token=[%+v]", SI.Token)
		logger.InfoLogger.Printf("Life=[%v]", SI.Life)
		logger.InfoLogger.Printf("SS=[%+v]", SS)
	}
	clock.Stop(1)
	return &SI
}

func GetSessionContext(w http.ResponseWriter, r *http.Request, sessionID string, userValidator func(int) (securityModel.UserMessage, error)) context.Context {

	ctx := r.Context()
	// Get the UserCode from the User Table, via the SessionID

	logger.SecurityLogger.Printf("[%v] EstablishSessionContext: Session=[%v]", strings.ToUpper(domain), sessionID)

	token, err := sessionStore.GetById(sessionID)
	if err != nil {
		logger.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("Session Not Found")
		Violation(w, r, msg.String())
		return ctx
	}

	logger.SecurityLogger.Printf("[%v] EstablishSessionContext: UserID=[%v]", strings.ToUpper(domain), token.UserID)
	clock := timing.Start(domain, "userValidator", "")
	UserMessage, err := userValidator(token.UserID)
	clock.Stop(1)
	if err == commonErrors.UserNotFound {
		logger.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Found")
		Violation(w, r, msg.String())
		return ctx
	}
	if err == commonErrors.UserNotActive {
		logger.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Active")
		Violation(w, r, msg.String())
		return ctx
	}
	if err != nil {
		logger.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Invalid")
		Violation(w, r, msg.String())
		return ctx
	}

	ctx = setSessionContextValues(ctx, UserMessage, sessionID, token)

	if appModeDev {
		logger.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionUserCodeKey, UserMessage.Code)
		logger.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionKey, sessionID)
		logger.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionUserIDKey, UserMessage.ID)
		logger.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionExpiryKey, token.Expiry)
	}

	return ctx
}

func setSessionContextValues(ctx context.Context, user securityModel.UserMessage, sessionID string, token sessionStore.Aegis_SessionStore) context.Context {
	ctx = context.WithValue(ctx, sessionUserCodeKey, user.Code)
	ctx = context.WithValue(ctx, sessionKey, sessionID)
	ctx = context.WithValue(ctx, sessionUserIDKey, user.ID)
	ctx = context.WithValue(ctx, sessionExpiryKey, token.Expiry)
	return ctx
}
