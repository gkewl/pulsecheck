package services

import (
	"encoding/json"
	"pulsecheck/authentication/core"
	"pulsecheck/authentication/model"
	"pulsecheck/authentication/parameters"
	//	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

func Login(requestUser *authmodel.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	if authBackend.Authenticate(requestUser) {
		token, err := authBackend.GenerateToken(requestUser.Actorid)
		if err != nil {
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(parameters.TokenAuthentication{token})
			return http.StatusOK, response
		}
	}

	return http.StatusUnauthorized, []byte("")
}

func MachineLogin(requestUser *authmodel.User) (int, []byte) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	token, err := authBackend.GenerateToken(requestUser.Actorid)
	if err != nil {
		return http.StatusInternalServerError, []byte("")
	} else {
		response, _ := json.Marshal(parameters.TokenAuthentication{token})
		return http.StatusOK, response
	}

	return http.StatusUnauthorized, []byte("")
}

//func RefreshToken(requestUser *auth.User) []byte {
//	authBackend := authentication.InitJWTAuthenticationBackend()
//	token, err := authBackend.GenerateToken(requestUser.UUID)
//	if err != nil {
//		fmt.Println("Cannot Refresh Token:" ,err)
//	}
//	response, err := json.Marshal(parameters.TokenAuthentication{token})
//	if err != nil {
//		fmt.Println("Cannot Marshal Token:",err)
//	}
//	return response
//}

//func Logout(req *http.Request) error {
//	authBackend := authentication.InitJWTAuthenticationBackend()
//	tokenRequest, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
//		return authBackend.PublicKey, nil
//	})
//	if err != nil {
//		return err
//	}
//	tokenString := req.Header.Get("Authorization")
//	return authBackend.Logout(tokenString, tokenRequest)
//}
