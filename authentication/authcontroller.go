package auth

import (
	"encoding/json"
	"net/http"
	"pulsecheck/authentication/model"
	"pulsecheck/authentication/services"
	"pulsecheck/common"
)

func GetRoutes() common.Routes {
	return common.Routes{
		common.Route{
			"Login",
			"POST",
			"/token-auth",
			LoginUser,
		},
		common.Route{
			"Machinelogin",
			"POST",
			"/machine_token-auth",
			MachineLoginUser,
		},
	}
}

func LoginUser(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	requestUser := new(authmodel.User)

	decoder := json.NewDecoder(r.Body)
	//; err != nil {
	//	newErr := errorhandler.ErrStepInvalidEntity
	//		newErr = errors.Annotate(newErr, "Error for json conversion: "+err.Error())
	//		return 0, GetErrorInfo(w, newErr, "Login User", "*")
	// 	}

	decoder.Decode(&requestUser)
	responseStatus, token := services.Login(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
	return 0, nil
}

func MachineLoginUser(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	requestUser := new(authmodel.User)

	decoder := json.NewDecoder(r.Body)
	// 	 ; err != nil {
	//	newErr := errorhandler.ErrStepInvalidEntity
	//		newErr = errors.Annotate(newErr, "Error for json conversion: "+err.Error())
	//		return 0, GetErrorInfo(w, newErr, "Login User", "*")
	// 	}

	decoder.Decode(&requestUser)
	responseStatus, token := services.MachineLogin(requestUser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseStatus)
	w.Write(token)
	return 0, nil
}

//func RefreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//	requestUser := new(models.User)
//	decoder := json.NewDecoder(r.Body)
//	decoder.Decode(&requestUser)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(services.RefreshToken(requestUser))
//}
