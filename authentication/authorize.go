package authentication

import (
	//	"github.com/gorilla/context"
	"context"
	"crypto/rsa"
	"crypto/x509"
	//"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	//"github.com/gkewl/pulsecheck/logger"
	"github.com/gkewl/pulsecheck/model"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}
	return authBackendInstance
}

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
			actorid = int64(claims["actorid"].(float64))
			role = claims["scopes"].(string)
		}
	}
	return
}

// Validate Token
func GetValidateToken(tokenStr string) (int64, error) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})
	if err != nil {
		return 0, eh.NewError(eh.ErrAuthTokenError, err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return int64(claims["actorid"].(float64)), nil
	}
	return 0, nil
}

//GenerateToken is the method that generates the token: for user a claims with expiration is attached, there
// are no exp claims on a machine token making it a permanent token
func (backend *JWTAuthenticationBackend) GenerateToken(username string, actorid int64, actorType string, role string) (model.TokenInfo, error) {
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
	claims["actorid"] = actorid
	claims["scopes"] = role
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
	data, _ := pem.Decode([]byte(config.GetEnv(config.MOS_PRIVATE_KEY_SECRET)))
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		fmt.Println("Importing PvtKey err ", err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	data, _ := pem.Decode([]byte(config.GetEnv(config.MOS_PUBLIC_KEY_SECRET)))
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

func (backend *JWTAuthenticationBackend) Authenticate(user model.User) bool {

	// The username and password we want to check
	//username := "rgunari"
	username := user.Username
	password := user.Password

	// bindusername := config.GetEnv(config.MOS_LDAP_USERNAME_PREFIX) + user.Username
	// bindpassword := user.Password
	// port, err := strconv.Atoi(config.GetEnv(config.MOS_LDAP_PORT))
	// if err != nil {
	// 	logger.LogInfo(fmt.Sprintf("Getting port failed. Error:%s", err.Error()), user.Xid)
	// }

	// l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", config.GetEnv(config.MOS_LDAP_URL), port))

	// if err != nil {
	// 	logger.LogInfo(fmt.Sprintf("Cannot connect to LDAP. Error:%s", err.Error()), user.Xid)
	// }
	// defer l.Close()

	// // First bind with a read only user
	// err = l.Bind(bindusername, bindpassword)
	// if err != nil {
	// 	logger.LogInfo(fmt.Sprintf("Invalid LDAP credentials. Error:%s", err.Error()), user.Xid)
	// }

	// // Search for the given username
	// searchRequest := ldap.NewSearchRequest(
	// 	config.GetEnv(config.MOS_LDAP_OU_DC),
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(&(sAMAccountName=%s))", username),
	// 	[]string{"dn"},
	// 	nil,
	// )

	// sr, err := l.Search(searchRequest)
	// if err != nil {
	// 	logger.LogInfo(fmt.Sprintf("Cannot search user. Error:%s", err.Error()), user.Xid)
	// 	return false
	// }

	// if len(sr.Entries) != 1 {
	// 	logger.LogInfo(fmt.Sprintf("User does not exist or too many entries:%s", sr.Entries[0]), user.Xid)
	// 	return false
	// }

	// userdn := sr.Entries[0].DN
	// // Bind as the user to verify their password
	// err = l.Bind(userdn, password)
	// if err != nil {
	// 	logger.LogInfo(fmt.Sprintf("Invalid password for user:%s   Error:%d", bindusername, err.Error()), user.Xid)
	// 	return false
	// }

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

//Authorize is the method the wraps a handlerfunc to serverhttp
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

func GetCurrentUser(r *http.Request) (string, int64) {
	username := GetCurrentUsername(r)
	userid := GetCurrentUserId(r)

	return username, userid
}

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

func GetPath() model.Settings {
	settings := model.Settings{}
	tokenexpiration, err := strconv.Atoi(config.GetEnv(config.MOS_TOKEN_TIMEOUT_IN_HOURS))
	if err != nil {
		tokenexpiration = 72 //user token expirationm time. Machine, we can check in generate token.
	}
	settings.JWTExpirationDelta = tokenexpiration

	return settings
}
