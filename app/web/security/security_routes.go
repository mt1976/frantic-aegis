package security

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mt1976/frantic-core/logHandler"
	"github.com/mt1976/frantic-core/messageHelpers"
)

func AnnounceInsecureRoute(route string) string {

	prefix := fmt.Sprintf("[ROUTE] Path=[%v]", route)
	padTo := 50
	if len(prefix) < padTo {
		prefix = prefix + strings.Repeat(" ", padTo-len(prefix))
	}
	logHandler.ApiLogger.Printf("%s %v://%v:%v%v", prefix, serverProtocol, serverHost, serverPort, route)
	return route
}

func AnnounceSecureRoute(route string) string {
	route = route + "/:" + sessionKey
	//logger.SecurityLogger.Printf("[%v] Secure Route: %v", strings.ToUpper(domain), route)
	return AnnounceInsecureRoute(route)
}

func EntryPoint(h httprouter.Handle, userKeyValidator func(string) (messageHelpers.UserMessage, error), userNameValidator func(string) (messageHelpers.UserMessage, error), authValidator func(string, string) error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// r.Header.Del("Authorization")
		// w.Header().Del("Authorization")

		// c := config.LoadConfigFile()
		// sessionKeyName := c.Security.SessionKeyName

		securityCheckOK := false

		userName := r.FormValue("username")
		password := r.FormValue("password")

		if userName == "" || password == "" {
			logHandler.SecurityLogger.Printf("[%v] Unauthorized Access Attempt - %v", strings.ToUpper(domain), http.StatusText(http.StatusUnauthorized))
			msg, _ := trnsl8.Get("Username and/or Password not provided")
			Violation(w, r, msg.String())
		}

		// Check this is a valid userMessage
		//userMessage, err := userStore.GetByUserName(userName)
		userMessage, err := userNameValidator(userName)
		if err != nil {
			logHandler.SecurityLogger.Printf("[%v] User not found - %v", strings.ToUpper(domain), userName)
			msg, _ := trnsl8.Get("User not found")
			Violation(w, r, msg.String())
		}
		logHandler.InfoLogger.Printf("User: %+v", userMessage)
		logHandler.InfoLogger.Printf("User Key: %+v", userMessage.Key)
		logHandler.InfoLogger.Printf("User Code: %+v", userMessage.Code)
		logHandler.SecurityLogger.Printf("[%v] User found - %v [%+v]", strings.ToUpper(domain), userName, userMessage)

		//err = auth.ValidateUserIDAndPassword(user.ID, password)
		err = authValidator(userMessage.Key, password)
		if err != nil {
			logHandler.SecurityLogger.Printf("[%v] Password mismatch - %v", strings.ToUpper(domain), userName)
			msg, _ := trnsl8.Get("Password mismatch")
			Violation(w, r, msg.String())
		} else {
			securityCheckOK = true
		}

		// Check this is a valid password
		if securityCheckOK {
			// Delegate request to the given handle
			si := New(r.Context(), userMessage.Key, userKeyValidator)
			r.Header.Add(SESSION_KEY, si.SessionID)

			newURL, _ := r.URL.Parse("/home")
			newPath := newURL.Query()
			newPath.Add(SESSION_KEY, si.SessionID)
			newURL.RawQuery = newPath.Encode()

			http.Redirect(w, r, newURL.String(), http.StatusFound)
		} else {
			// Request Basic Authentication otherwise
			msg := http.StatusText(http.StatusUnauthorized)

			if err != nil {
				msg = err.Error()
			}

			msgx, _ := trnsl8.Get(msg)
			msg = msgx.String()
			Violation(w, r, msg)
		}
	}
}

func Violation(w http.ResponseWriter, r *http.Request, msg string) {

	oldTableName := domain
	domain = "!!!!!EXCEPTION!!!!!"
	logHandler.SecurityLogger.Printf("[%v] Unauthorized Access Attempt - %v", strings.ToUpper(domain), http.StatusText(http.StatusUnauthorized))
	logHandler.SecurityLogger.Printf("[%v] Reason : [%v]", strings.ToUpper(domain), msg)
	logHandler.SecurityLogger.Printf("[%v] From : [%v]", strings.ToUpper(domain), r.RemoteAddr)
	logHandler.SecurityLogger.Printf("[%v] Request : [%v]", strings.ToUpper(domain), r.RequestURI)
	logHandler.SecurityLogger.Printf("[%v] Method : [%v]", strings.ToUpper(domain), r.Method)
	logHandler.SecurityLogger.Printf("[%v] User Agent : [%v]", strings.ToUpper(domain), r.UserAgent())
	logHandler.SecurityLogger.Printf("[%v] Referer : [%v]", strings.ToUpper(domain), r.Referer())
	logHandler.SecurityLogger.Printf("[%v] Host : [%v]", strings.ToUpper(domain), r.Host)
	logHandler.SecurityLogger.Printf("[%v] Remote Address : [%v]", strings.ToUpper(domain), r.RemoteAddr)
	domain = oldTableName

	if msg == "" {
		msg = http.StatusText(http.StatusUnauthorized)
	}

	uri := url.URL{}
	uri.Path = "/fail"
	v := url.Values{}
	//uriQuery := uri.Query()
	//fmt.Printf("v: %v\n", v)
	v.Add(msgTypeKey, "error")
	//fmt.Printf("v: %v\n", v)
	msg2, _ := trnsl8.Get("Unauthorized Access")
	v.Add(msgTitleKey, msg2.String())
	//fmt.Printf("v: %v\n", v)
	msg3, _ := trnsl8.Get("Access is prohibited - %v")
	v.Add(msgContentKey, fmt.Sprintf(msg3.String(), msg))
	//fmt.Printf("v: %v\n", v)
	msg4, _ := trnsl8.Get("Security")
	v.Add(msgActionKey, msg4.String())
	//fmt.Printf("v: %v\n", v)
	uri.RawQuery = v.Encode()

	//fmt.Printf("Query: %+v\n", v.Encode())
	//fmt.Printf("URI: %+v\n", uri.String())

	////spew.Dump(uri)

	logHandler.SecurityLogger.Printf("[%v] Redirecting to [%v]", domain, uri.String())

	http.Redirect(w, r, uri.String(), http.StatusFound)

	// if err != nil {
	// 	logger.ErrorLogger.Fatalf("Error=[%v]", err.Error())
	// }

}

func Validate(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logHandler.SecurityLogger.Printf("[%v] Session Key Name: [%v]", strings.ToUpper(domain), sessionKey)

		sessionID := extractSessionID(ps, sessionKey, r)

		logHandler.SecurityLogger.Printf("[%v] Session Key FOUND!: [%v]", strings.ToUpper(domain), sessionID)

		if sessionID != "" {
			// Delegate request to the given handle
			logHandler.SecurityLogger.Printf("[%v] Session Validated - [%v] - [%v]", strings.ToUpper(domain), r.RequestURI, http.StatusText(http.StatusAccepted))
			query := r.URL.Query()
			query.Add(sessionKey, sessionID)
			r.URL.RawQuery = query.Encode()
			//	logHandler.SecurityLogger.Printf("url adding [%v=%v]\n", sessionKey, sessionID)
			r.URL.RawQuery = r.URL.Query().Encode()
			//	logHandler.SecurityLogger.Printf("ps=%+v", ps)
			//	logHandler.SecurityLogger.Printf("r=%+v", r.URL.Query().Encode())
			h(w, r, ps)
		} else {
			// Error Response
			logHandler.SecurityLogger.Printf("[%v] Session Violation - [%v] - [%v]", strings.ToUpper(domain), r.RequestURI, http.StatusText(http.StatusUnauthorized))
			msg2, _ := trnsl8.Get("Session Key not found")
			msg := http.StatusText(http.StatusUnauthorized) + " - " + msg2.String()
			Violation(w, r, msg)
		}
	}
}
