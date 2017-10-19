package authentication

import (
	//	"github.com/gorilla/context"

	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	//"encoding/json"

	"fmt"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	//"github.com/gkewl/pulsecheck/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gkewl/pulsecheck/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// define global map; initialize as empty with the trailing {}
var roleMapping = map[string]int{
	"Guest":     constant.Guest,
	"Superuser": constant.Superuser,
	"Admin":     constant.Admin,
}

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var authBackendInstance *JWTAuthenticationBackend

// InitJWTAuthenticationBackend
func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}
	return authBackendInstance
}

// RequireTokenAuthentication
func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
	authBackend := InitJWTAuthenticationBackend()

	token, err := request.ParseFromRequest(req, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	x := map[string]interface{}{}

	if err == nil && token.Valid {
		x = token.Claims.(jwt.MapClaims)
	}
	if err != nil {
		if err == request.ErrNoTokenInRequest {
			err = eh.NewError(eh.ErrAuthTokenNotAvailable, "")
		} else {
			err = eh.NewError(eh.ErrAuthTokenError, err.Error())
		}

	}
	return x, err
}

// GetUserInfoFromToken parses the raw token and returns user info
func GetUserInfoFromToken(tokenStr string) (username string, actorid int64, role string) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			username = claims["sub"].(string)
			actorid = int64(claims["userid"].(float64))
			//role = claims["scopes"].(string)
		}
	}
	return
}

// GetValidateToken Validate Token
func GetValidateToken(tokenStr string) (int64, error) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return 0, eh.NewError(eh.ErrAuthTokenError, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return int64(claims["userid"].(float64)), nil
	}
	return 0, nil
}

//GenerateToken is the method that generates the token: for user a claims with expiration is attached, there
// are no exp claims on a machine token making it a permanent token
func (backend *JWTAuthenticationBackend) GenerateToken(username string, userInfo model.UserCompany, actorType string) (model.TokenInfo, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := token.Claims.(jwt.MapClaims)
	expireTime := time.Now().Add(time.Hour * time.Duration(GetPath().JWTExpirationDelta)).Unix()
	if strings.ToLower(actorType) == "user" {
		claims["exp"] = expireTime
	} else {
		expireTime = 0
	}

	claims["iat"] = time.Now().Unix()
	//claims["iss"] = config.AuthIssuer()

	claims["sub"] = username
	claims["userid"] = userInfo.UserID
	claims["companyid"] = userInfo.CompanyID
	claims["scope"] = userInfo.Role

	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		return model.TokenInfo{}, err
	}

	tokenInfo := model.TokenInfo{Token: tokenString, Exp: expireTime}
	return tokenInfo, nil
}

func (backend *JWTAuthenticationBackend) getTokenRemainingValidity(timestamp interface{}) int {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remainer := tm.Sub(time.Now())
		if remainer > 0 {
			return int(remainer.Seconds() + expireOffset)
		}
	}
	return expireOffset
}

func getPrivateKey() *rsa.PrivateKey {
	// privateKeyFile, err := os.Open(config.GetEnv(config.PULSE_PRIVATE_KEY_FILEPATH))
	// if err != nil {
	// 	fmt.Println("Cannot find pvt key", err)
	// }

	// pemfileinfo, _ := privateKeyFile.Stat()
	// var size = pemfileinfo.Size()
	// pembytes := make([]byte, size)

	// buffer := bufio.NewReader(privateKeyFile)
	// _, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(constant.PULSE_PRIVATE_KEY))
	// privateKeyFile.Close()
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println("Importing PvtKey err ", err)
	}
	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	// privateKeyFile, err := os.Open(config.GetEnv(config.PULSE_PUBLIC_KEY_FILEPATH))
	// if err != nil {
	// 	fmt.Println("Cannot find pvt key", err)
	// }

	// pemfileinfo, _ := privateKeyFile.Stat()
	// var size = pemfileinfo.Size()
	// pembytes := make([]byte, size)

	// buffer := bufio.NewReader(privateKeyFile)
	// _, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(constant.PULSE_PUBLIC_KEY))
	// privateKeyFile.Close()
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		fmt.Println("Public Key Import Err")
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		fmt.Println("Public Key not ok")
	}

	return rsaPub
}

func (backend *JWTAuthenticationBackend) Authenticate(user model.AuthenticateUser) bool {
	return true
}

//Authorize is the method the wraps a handlerfunc to serverhttp
func Authorize(pass common.AppHandlerFunc) common.AppHandlerFunc {

	return func(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
		claims, err := RequireTokenAuthentication(w, r)
		if err != nil {
			return 0, eh.RespWithError(w, r, nil, err)
		}
		//context.Set(r, constant.Claims, claims)
		goCtx := context.WithValue(r.Context(), constant.Claims, claims)
		r = r.WithContext(goCtx)

		return pass(ctx, w, r)
	}
}

//AuthorizeWithRole is the method the wraps a handlerfunc to serverhttp
func AuthorizeWithRole(pass common.AppHandlerFunc, apiRole int) common.AppHandlerFunc {

	return func(ctx *common.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
		claims, err := RequireTokenAuthentication(w, r)
		if err != nil {
			newErr := eh.NewError(eh.ErrAuthUserNameNotAvailable, "")
			return 0, eh.RespWithError(w, r, nil, newErr)
		}

		//validate role
		if role, ok := claims["scopes"].(string); ok {
			if roleMapping[role] >= apiRole {
				//context.Set(r, constant.Claims, claims)
				goCtx := context.WithValue(r.Context(), constant.Claims, claims)
				r = r.WithContext(goCtx)
				return pass(ctx, w, r)
			}
		}

		newErr := eh.NewError(eh.ErrUnAuthorizedUserForAPI, "")
		return 0, eh.RespWithError(w, r, nil, newErr)
	}
}

// AnnotateRequestWithAuthorizedUser checks the inbound request for a token
// with a scopes claim that is greater or equal to the desired api role. If
// the requirement is met, an http.Request is returned with the claims attached.
// If there was no token or insufficient privilege, the error has been written
// to the response writer and nil is returned.
func AnnotateRequestWithAuthorizedUser(w http.ResponseWriter, r *http.Request, apiRole int) *http.Request {
	var newErr error
	claims, err := RequireTokenAuthentication(w, r)
	if err == nil {
		//validate role
		if role, ok := claims["scopes"].(string); ok {
			if roleMapping[role] >= apiRole {
				goCtx := context.WithValue(r.Context(), constant.Claims, claims)
				return r.WithContext(goCtx)
			}
		}
		newErr = eh.NewError(eh.ErrUnAuthorizedUserForAPI, "")
	} else {
		newErr = err
	}
	eh.RespWithError(w, r, nil, newErr)
	return nil
}

// GetCurrentUsername
func GetCurrentUsername(r *http.Request) string {
	claims := r.Context().Value(constant.Claims)
	if claims == nil {
		return ""
	}

	//get map from interface
	claimMap := claims.(map[string]interface{})
	if username, ok := claimMap["sub"].(string); ok {
		return username
	}

	return ""
}

// GetCurrentUserId
func GetCurrentUserId(r *http.Request) int64 {
	//	claims, ok := context.GetOk(r, constant.Claims)
	//	if !ok {
	//		return 0
	//	}
	claims := r.Context().Value(constant.Claims)
	if claims == nil {
		return 0
	}
	//get map from interface
	claimMap := claims.(map[string]interface{})
	if id, ok := claimMap["actorid"].(float64); ok {
		return int64(id)
	}
	return 0
}

// GetCurrentUserRole
func GetCurrentUserRole(r *http.Request) string {
	claims := r.Context().Value(constant.Claims)
	if claims == nil {
		return ""
	}

	//get map from interface
	claimMap := claims.(map[string]interface{})
	if role, ok := claimMap["scopes"].(string); ok {
		return role
	}
	return ""
}

// HasPermission
func HasPermission(r *http.Request, apiRole int) bool {

	roleMapping := map[string]int{
		"Guest": constant.Guest, "Superuser": constant.Superuser, "Admin": constant.Admin,
	}

	role := GetCurrentUserRole(r)

	if roleMapping[role] >= apiRole {
		return true
	}
	return false
}

// GetCurrentUser
func GetCurrentUser(r *http.Request) (string, int64) {
	username := GetCurrentUsername(r)
	userid := GetCurrentUserId(r)

	return username, userid
}

// GetUserInfoFromContext
func GetUserInfoFromContext(ctx context.Context) (name string, ID int64, role string) {
	if ctx == nil {
		return "", 0, ""
	}
	claims := ctx.Value(constant.Claims)
	if claims == nil {
		return "", 0, ""
	}
	// get map from interface
	var ok bool
	claimMap := claims.(map[string]interface{})
	if fid, ok := claimMap["actorid"].(float64); !ok {
		return "", 0, ""
	} else {
		ID = int64(fid)
	}
	if name, ok = claimMap["sub"].(string); !ok {
		return "", 0, ""
	}
	if role, ok = claimMap["scopes"].(string); !ok {
		return "", 0, ""
	}
	return
}

// RoleIdForName returns the integer role id for the given role
// name or -1 if the name is not valid
func RoleIdForName(role string) int {
	if val, present := roleMapping[role]; present {
		return val
	}
	return -1
}

// RoleName returns the string name for a role id or "" if the id
// is not valid
func RoleName(role int) string {
	for name, val := range roleMapping {
		if val == role {
			return name
		}
	}
	return ""
}

// GetPath
func GetPath() model.Settings {
	settings := model.Settings{}
	tokenexpiration, err := strconv.Atoi(config.GetEnv(config.PULSE_TOKEN_TIMEOUT_IN_HOURS))
	if err != nil {
		tokenexpiration = 72 //user token expirationm time. Machine, we can check in generate token.
	}
	settings.JWTExpirationDelta = tokenexpiration

	return settings
}
