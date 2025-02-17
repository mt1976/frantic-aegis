package security

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mt1976/frantic-aegis/app/web/security/securityModel"
	logger "github.com/mt1976/frantic-core/logHandler"
)

func AnnounceInsecureRoute(route string) string {

	prefix := fmt.Sprintf("[ROUTE] Path=[%v]", route)
	padTo := 50
	if len(prefix) < padTo {
		prefix = prefix + strings.Repeat(" ", padTo-len(prefix))
	}
	logger.ApiLogger.Printf("%s %v://%v:%v%v", prefix, serverProtocol, serverHost, serverPort, route)
	return route
}

func AnnounceSecureRoute(route string) string {
	route = route + "/:" + sessionKey
	//logger.SecurityLogger.Printf("[%v] Secure Route: %v", strings.ToUpper(domain), route)
	return AnnounceInsecureRoute(route)
}

func EntryPoint(h httprouter.Handle, userIDValidator func(int) (securityModel.UserMessage, error), userNameValidator func(string) (securityModel.UserMessage, error), authValidator func(int, string) error) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// r.Header.Del("Authorization")
		// w.Header().Del("Authorization")

		// c := config.LoadConfigFile()
		// sessionKeyName := c.Security.SessionKeyName

		securityCheckOK := false

		userName := r.FormValue("username")
		password := r.FormValue("password")

		if userName == "" || password == "" {
			logger.SecurityLogger.Printf("[%v] Unauthorized Access Attempt - %v", strings.ToUpper(domain), http.StatusText(http.StatusUnauthorized))
			msg, _ := trnsl8.Get("Username and/or Password not provided")
			Violation(w, r, msg.String())
		}

		// Check this is a valid user
		//user, err := userStore.GetByUserName(userName)
		user, err := userNameValidator(userName)
		if err != nil {
			logger.SecurityLogger.Printf("[%v] User not found - %v", strings.ToUpper(domain), userName)
			msg, _ := trnsl8.Get("User not found")
			Violation(w, r, msg.String())
		}

		//err = auth.ValidateUserIDAndPassword(user.ID, password)
		err = authValidator(user.ID, password)
		if err != nil {
			logger.SecurityLogger.Printf("[%v] Password mismatch - %v", strings.ToUpper(domain), userName)
			msg, _ := trnsl8.Get("Password mismatch")
			Violation(w, r, msg.String())
		} else {
			securityCheckOK = true
		}

		// Check this is a valid password
		if securityCheckOK {
			trace(r)
			// Delegate request to the given handle
			si := New(r.Context(), user.ID, userIDValidator)
			//r.URL.Query().Add(sessionKey, si.SessionID)
			r.Header.Add(sessionKey, si.SessionID)
			newURI := "/home?" + sessionKey + "=" + si.SessionID
			http.Redirect(w, r, newURI, http.StatusFound)
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
	logger.SecurityLogger.Printf("[%v] Unauthorized Access Attempt - %v", strings.ToUpper(domain), http.StatusText(http.StatusUnauthorized))
	logger.SecurityLogger.Printf("[%v] Reason : [%v]", strings.ToUpper(domain), msg)
	logger.SecurityLogger.Printf("[%v] From : [%v]", strings.ToUpper(domain), r.RemoteAddr)
	logger.SecurityLogger.Printf("[%v] Request : [%v]", strings.ToUpper(domain), r.RequestURI)
	logger.SecurityLogger.Printf("[%v] Method : [%v]", strings.ToUpper(domain), r.Method)
	logger.SecurityLogger.Printf("[%v] User Agent : [%v]", strings.ToUpper(domain), r.UserAgent())
	logger.SecurityLogger.Printf("[%v] Referer : [%v]", strings.ToUpper(domain), r.Referer())
	logger.SecurityLogger.Printf("[%v] Host : [%v]", strings.ToUpper(domain), r.Host)
	logger.SecurityLogger.Printf("[%v] Remote Address : [%v]", strings.ToUpper(domain), r.RemoteAddr)
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

	logger.SecurityLogger.Printf("[%v] Redirecting to [%v]", domain, uri.String())

	http.Redirect(w, r, uri.String(), http.StatusFound)

	// if err != nil {
	// 	logger.ErrorLogger.Fatalf("Error=[%v]", err.Error())
	// }

}

func Validate(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		logger.SecurityLogger.Printf("[%v] Session Key Name: [%v]", strings.ToUpper(domain), sessionKey)

		sessionID := extractSessionID(ps, sessionKey, r)

		logger.SecurityLogger.Printf("[%v] Session Key FOUND!: [%v]", strings.ToUpper(domain), sessionID)

		if sessionID != "" {
			// Delegate request to the given handle
			logger.SecurityLogger.Printf("[%v] Session Validated - [%v] - [%v]", strings.ToUpper(domain), r.RequestURI, http.StatusText(http.StatusAccepted))
			r.URL.Query().Add(sessionKey, sessionID)
			logger.SecurityLogger.Printf("url adding sessionKey: %v sessionID:%v\n", sessionKey, sessionID)
			r.URL.RawQuery = r.URL.Query().Encode()
			logger.SecurityLogger.Printf("ps=%+v", ps)
			logger.SecurityLogger.Printf("r=%+v", r.URL.Query().Encode())
			h(w, r, ps)
		} else {
			// Error Response
			logger.SecurityLogger.Printf("[%v] Session Violation - [%v] - [%v]", strings.ToUpper(domain), r.RequestURI, http.StatusText(http.StatusUnauthorized))
			msg2, _ := trnsl8.Get("Session Key not found")
			msg := http.StatusText(http.StatusUnauthorized) + " - " + msg2.String()
			Violation(w, r, msg)
		}
	}
}
