package common

import (
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
)

//this interface is used for every class which can be
type Routes []Route

type AppHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc AppHandlerFunc
}

type APIRoutes []Routes

// AppContext contains our local context; our database pool, session store, template
// registry and anything else our handlers need to access. We'll create an instance of it
// in our main() function and then explicitly pass a reference to it for our handlers to access.
type AppContext struct {
	Db *sqlx.DB
	//Store *sessions.CookieStore
	//templates   map[string]*template.Template
	//decoder     *schema.Decoder
	//sRedisPool *redis.Pool
	UseMock bool
}

// AppHandler - wrap http.handler embed field *AppContext
// avoid globals and improve middleware chaining and error handling
type AppHandler struct {
	*AppContext
	H AppHandlerFunc
}

//AH function - add http handler to AppHandler type
func (AH AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.appContext as a parameter to our handler type.

	_, err := AH.H(AH.AppContext, w, r)
	if err != nil {

		//		logger.Log(logger.LogModel{
		//			Level:  logger.INFO,
		//			Caller: "common.AppHandler.ServeHTTP",
		//			Msg:    fmt.Sprintf("HTTP %d: %q", status, err),
		//			Err:    err,
		//		})

	}
}

func CaselessMatcher(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
