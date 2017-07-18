package authentication

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"pulsecheck/authentication/parameters"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	tokenDuration = 72
	expireOffset  = 3600
)

//func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
//	authBackend := InitJWTAuthenticationBackend()
//
//	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		} else {
//			return authBackend.PublicKey, nil
//		}
//	})
//
//	if err == nil && token.Valid && !authBackend.IsInBlacklist(req.Header.Get("Authorization")) {
//		next(rw, req)
//	} else {
//		rw.WriteHeader(http.StatusUnauthorized)
//	}
//}

type JWTAuthenticationBackend struct {
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var authBackendInstance *JWTAuthenticationBackend = nil

func InitJWTAuthenticationBackend() *JWTAuthenticationBackend {
	if authBackendInstance == nil {
		authBackendInstance = &JWTAuthenticationBackend{
			privateKey: getPrivateKey(),
			PublicKey:  getPublicKey(),
		}
	}

	return authBackendInstance
}

//var scopes = ["admin","enduser"]string

func (backend *JWTAuthenticationBackend) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(parameters.GetPath().JWTExpirationDelta)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = username
	token.Claims = claims
	//token.Claims["scopes"] = "adminUser","enduser"
	tokenString, err := token.SignedString(backend.privateKey)
	if err != nil {
		//panic(err)
		return "Cannot Generate Token", err
	}
	return tokenString, nil
}

//func (backend *JWTAuthenticationBackend) Authenticate(user *authmodel.User) bool {
//	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testing"), 10)
//
//	testUser := authmodel.User{
//		UUID:     uuid.New(),
//		Username: "gokul",
//		//Password: "test",
//		Password: "testing",
//	}
//return user.Username == testUser.Username && user.Password == testUser.Password
//	//return user.Username == testUser.Username && bcrypt.CompareHashAndPassword([]byte(testUser.Password), []byte(user.Password)) == nil
//}

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

//func (backend *JWTAuthenticationBackend) Logout(tokenString string, token *jwt.Token) error {
//	redisConn := redis.Connect()
//	return redisConn.SetValue(tokenString, tokenString, backend.getTokenRemainingValidity(token.Claims["exp"]))
//}
//
//func (backend *JWTAuthenticationBackend) IsInBlacklist(token string) bool {
//	redisConn := redis.Connect()
//	redisToken, _ := redisConn.GetValue(token)
//
//	if redisToken == nil {
//		return false
//	}
//
//	return true
//}

//x509.parseP

func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(parameters.GetPath().PrivateKeyPath)
	if err != nil {
		fmt.Println("Cannot find pvt key", err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	fmt.Println("pemfileinfo", pemfileinfo, pemfileinfo.Size())
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))
	fmt.Println("Data", data, data.Bytes)
	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	if err != nil {
		fmt.Println("Importing PvtKey err ", err)
	}

	return privateKeyImported
}

func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(parameters.GetPath().PublicKeyPath)
	if err != nil {
		fmt.Println("Cannot find Public  Key Path", err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyFile.Close()

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

//func(token *jwt.Token)(interface{},error){
//	return authBackend.PublicKey,nil
//}

func RequireTokenAuthentication(rw http.ResponseWriter, req *http.Request) (bool, error) {
	authBackend := InitJWTAuthenticationBackend()
	token, err := jwt.ParseFromRequest(
		req, func(token *jwt.Token) (interface{}, error) {
			return authBackend.PublicKey, nil
		})

	if err == nil && token.Valid {
		fmt.Println("valid token", err, token.Valid)
		fmt.Println(token.Claims["scopes"])
		fmt.Println(token.Claims["sub"])
		return true, nil

	} else {
		fmt.Println("invalid token", err)
		return false, err

	}

}
