package security

import (
	"context"
	"net/http"
	"strings"

	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userID int, userIDValidator func(int) (securityModel.UserMessage, error)) *securityModel.Session {
	clock := timing.Start(domain, "New", "")
	SI := securityModel.Session{}

	// U, err := userStore.Get(userID)
	// if err != nil {
	// 	panic(err)
	// }
	UserMessage, err := userIDValidator(userID)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		panic(err)
	}

	SS, err := sessionStore.New(UserMessage.ID, UserMessage.Code)
	if err != nil {
		panic(err)
	}

	SI.Token = SS
	SI.UserID = UserMessage.ID
	SI.SessionID = SS.SessionID
	SI.Life = 0
	SI.UserCode = UserMessage.Code

	//ctx = setSessionContextValues(ctx, UserMessage, SI.SessionID, SS)

	if appModeDev {
		logHandler.InfoLogger.Printf("SessionID=[%v]", SI.SessionID)
		logHandler.InfoLogger.Printf("UserID=[%v]", SI.UserID)
		logHandler.InfoLogger.Printf("UserCode=[%v]", SI.UserCode)
		logHandler.InfoLogger.Printf("Token=[%+v]", SI.Token)
		logHandler.InfoLogger.Printf("Life=[%v]", SI.Life)
		logHandler.InfoLogger.Printf("SS=[%+v]", SS)
	}
	clock.Stop(1)
	return &SI
}

func GetSessionContext(w http.ResponseWriter, r *http.Request, sessionID string, userValidator func(int) (securityModel.UserMessage, error)) context.Context {

	ctx := r.Context()
	// Get the UserCode from the User Table, via the SessionID

	logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: Session=[%v]", strings.ToUpper(domain), sessionID)

	token, err := sessionStore.GetBy(sessionStore.FIELD_SessionID, sessionID)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("Session Not Found")
		Violation(w, r, msg.String())
		return ctx
	}

	logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: UserID=[%v] (%v)", strings.ToUpper(domain), token.UserID, token.UserCode)
	clock := timing.Start(domain, "userValidator", "")
	UserMessage, err := userValidator(token.UserID)
	clock.Stop(1)
	if err == commonErrors.ErrorUserNotFound {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Found")
		Violation(w, r, msg.String())
		return ctx
	}
	if err == commonErrors.ErrorUserNotActive {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Active")
		Violation(w, r, msg.String())
		return ctx
	}
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Invalid")
		Violation(w, r, msg.String())
		return ctx
	}

	ctx = setSessionContextValues(ctx, UserMessage, sessionID, token)

	if appModeDev {
		logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionUserCodeKey, UserMessage.Code)
		logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionKey, sessionID)
		logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionUserIDKey, UserMessage.ID)
		logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionUserCodeKey, UserMessage.Code)
		logHandler.SecurityLogger.Printf("[%v] EstablishSessionContext: [%v]=[%v]", strings.ToUpper(domain), sessionExpiryKey, token.Expiry)
	}

	return ctx
}

func setSessionContextValues(ctx context.Context, user securityModel.UserMessage, sessionID string, token sessionStore.Session_Store) context.Context {
	ctx = context.WithValue(ctx, sessionUserCodeKey, user.Code)
	ctx = context.WithValue(ctx, sessionKey, sessionID)
	ctx = context.WithValue(ctx, sessionUserIDKey, user.ID)
	ctx = context.WithValue(ctx, sessionUserCodeKey, user.Code)
	ctx = context.WithValue(ctx, sessionExpiryKey, token.Expiry)
	return ctx
}
