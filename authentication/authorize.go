package auth

import (
	"fmt"
	"net/http"
	"pulsecheck/authentication/core"
	"pulsecheck/common"
)

func Authorize(pass common.AppHandlerFunc) common.AppHandlerFunc {

	return func(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
		_, err := authentication.RequireTokenAuthentication(w, r)
		if err != nil {
			fmt.Printlln("http.StatusUnauthorized")
			return 0, err
		}
		return pass(ctx, w, r)
	}
}
