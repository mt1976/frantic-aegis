package security

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mt1976/frantic-core/logHandler"
)

var domain = "security"

func ExtractSessionTokenFromReferer(r *http.Request) string {
	referer := r.Referer()
	refURI, err := url.Parse(referer)
	if err != nil {
		logHandler.ErrorLogger.Printf("[%v] Error=[%v]", strings.ToUpper(domain), err.Error())
		return ""
	}
	//fmt.Printf("key: %v\n", key)
	//fmt.Printf("refURI: %+v\n", refURI)
	//fmt.Printf("referer: %+v\n", referer)
	sessionID := refURI.Query().Get(sessionKey)
	if sessionID == "" {
		//get the last part of the referer delimited by /
		parts := strings.Split(referer, "/")
		sessionID = parts[len(parts)-1]
	}
	//fmt.Printf("ret: %+v\n", ret)
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
