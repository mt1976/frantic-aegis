package security

import (
	"net/http"

	"github.com/mt1976/frantic-core/logger"
)

// func templatedHTML() string {
// 	where := paths.Res().String() + "/html/templates.html"
// 	//logger.InfoLogger.Printf("[TEMPLATE] Template Loc=[%v]", where)
// 	return where
// }

func trace(r *http.Request) {
	mesg := "[SECURITY] Method=[%s] URI=[%s] Header[%v]"
	logger.TraceLogger.Printf(mesg, r.Method, r.URL, r.Header)
}
