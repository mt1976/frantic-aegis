package security

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mt1976/frantic-aegis/app/dao/sessionStore"
	"github.com/mt1976/frantic-core/commonErrors"
	"github.com/mt1976/frantic-core/contextHandler"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
	"github.com/mt1976/frantic-core/timing"
)

func New(ctx context.Context, userKey string, userIDValidator func(string) (messageHelpers.UserMessage, error)) *messageHelpers.SessionMessage {
	clock := timing.Start(domain, "New", "")
	SI := messageHelpers.SessionMessage{}

	// U, err := userStore.Get(userID)
	// if err != nil {
	// 	panic(err)
	// }
	UserMessage, err := userIDValidator(userKey)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		panic(err)
	}

	SS, err := sessionStore.New(ctx, UserMessage.Key, UserMessage.Code)
	if err != nil {
		panic(err)
	}

	//SI.SessionToken = SS
	SI.UserKey = UserMessage.Key
	SI.UserCode = UserMessage.Code
	SI.SessionID = SS.SessionID
	SI.Expiry = SS.Expiry
	SI.User.Locale = UserMessage.Locale

	//ctx = setSessionContextValues(ctx, UserMessage, SI.SessionID, SS)

	if appModeDev {
		logHandler.InfoLogger.Printf("SessionID=[%v]", SI.SessionID)
		//logHandler.InfoLogger.Printf("SessionToken=[%+v]", SI.SessionToken)
		logHandler.InfoLogger.Printf("UserKey=[%v]", SI.UserKey)
		logHandler.InfoLogger.Printf("UserCode=[%v]", SI.UserCode)
		logHandler.InfoLogger.Printf("Life=[%v]", SI.Expiry)
		logHandler.InfoLogger.Printf("SS=[%+v]", SS)
		logHandler.InfoLogger.Printf("UserMessage=[%+v]", UserMessage)
		logHandler.InfoLogger.Printf("SessionMessage=[%+v]", SI)
		logHandler.InfoLogger.Printf("Locale=[%v]", SI.Locale)
		logHandler.InfoLogger.Printf("Theme=[%v]", SI.Spare1)
		logHandler.InfoLogger.Printf("Timezone=[%v]", SI.Spare2)
	}
	clock.Stop(1)
	return &SI
}

func GetSessionContext(w http.ResponseWriter, r *http.Request, ps httprouter.Params, userValidator func(string) (messageHelpers.UserMessage, error)) (context.Context, *http.Request) {

	//ps httprouter.Params
	sessionID := ps.ByName(cfg.GetSecuritySessionKey_Session())
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking headers [%v]", strings.ToUpper(domain), r.Header)
	}
	sessionID = extractSessionID(ps, cfg.GetSecuritySessionKey_Session(), r)
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] Unable to find session id to set context", strings.ToUpper(domain))
		msg, _ := trnsl8.Get("Session Error")
		Violation(w, r, msg.String())
		return r.Context(), r
	}

	ctx := r.Context()
	// Get the UserCode from the User Table, via the SessionID

	logHandler.SecurityLogger.Printf("[%v] GetSessionContext: Session=[%v]", strings.ToUpper(domain), sessionID)

	userSessionTokenRecord, err := sessionStore.GetBy(sessionStore.FIELD_SessionID, sessionID)
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("Session Not Found")
		Violation(w, r, msg.String())
		return r.Context(), r
	}

	logHandler.SecurityLogger.Printf("[%v] GetSessionContext: UserKey=[%v] (%v)", strings.ToUpper(domain), userSessionTokenRecord.UserKey, userSessionTokenRecord.UserCode)
	clock := timing.Start(domain, "userValidator", "")
	UserMessage, err := userValidator(userSessionTokenRecord.UserKey)
	clock.Stop(1)
	if err == commonErrors.ErrorUserNotFound {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Found")
		Violation(w, r, msg.String())
		return r.Context(), r
	}
	if err == commonErrors.ErrorUserNotActive {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Not Active")
		Violation(w, r, msg.String())
		return r.Context(), r
	}
	if err != nil {
		logHandler.ErrorLogger.Printf("Error=[%v]", err.Error())
		msg, _ := trnsl8.Get("User Invalid")
		Violation(w, r, msg.String())
		return r.Context(), r
	}

	ctx = setSessionContextValues(ctx, UserMessage, sessionID, userSessionTokenRecord)

	//	if appModeDev {
	contextHandler.Debug(ctx, "Request-Session-OK")

	//	}

	return ctx, r.WithContext(ctx)
}

func ExtractSessionTokenFromReferer(r *http.Request) string {
	referer := r.Referer()
	refURI, err := url.Parse(referer)
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(domain), err.Error())
		return ""
	}

	sessionID := refURI.Query().Get(sessionKey)
	if sessionID == "" {
		//get the last part of the referer delimited by /
		parts := strings.Split(referer, "/")
		sessionID = parts[len(parts)-1]
	}
	return sessionID
}

func extractSessionID(ps httprouter.Params, sessionKeyName string, r *http.Request) string {
	sessionID := ps.ByName(sessionKeyName)
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking headers [%v]", strings.ToUpper(domain), r.Header)
		sessionID = r.Header.Get(sessionKeyName)
	}
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking params [%v]", strings.ToUpper(domain), r.URL.Query())
		sessionID = r.URL.Query().Get(sessionKeyName)
	}
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking form [%v]", strings.ToUpper(domain), r.Form)
		sessionID = r.FormValue(sessionKeyName)
	}
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking context [%v]", strings.ToUpper(domain), r.Context())
		ID := r.Context().Value(sessionKeyName)
		if ID != nil {
			sessionID = ID.(string)
		}
	}
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking cookies [%v]", strings.ToUpper(domain), r.Cookies())
		cookie, err := r.Cookie(sessionKeyName)
		if err == nil {
			sessionID = cookie.Value
		}
	}
	if sessionID == "" {
		logHandler.SecurityLogger.Printf("[%v] No Session Key Found, checking referer [%v]", strings.ToUpper(domain), r.Referer())
		sessionID = ExtractSessionTokenFromReferer(r)
	}
	// remove the last char if it is a ?
	sessionID = strings.TrimSuffix(sessionID, "?")
	return sessionID
}

func ExtractSessionID(ps httprouter.Params, r *http.Request) string {
	return extractSessionID(ps, sessionKey, r)
}
